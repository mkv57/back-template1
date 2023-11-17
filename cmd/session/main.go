package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gofrs/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sipki-tech/database/connectors"
	"google.golang.org/grpc/grpclog"
	"gopkg.in/yaml.v3"

	"github.com/ZergsLaw/back-template/cmd/session/internal/adapters/queue"
	"github.com/ZergsLaw/back-template/cmd/session/internal/adapters/repo"
	"github.com/ZergsLaw/back-template/cmd/session/internal/api"
	"github.com/ZergsLaw/back-template/cmd/session/internal/app"
	"github.com/ZergsLaw/back-template/cmd/session/internal/auth"
	"github.com/ZergsLaw/back-template/internal/flags"
	"github.com/ZergsLaw/back-template/internal/grpchelper"
	"github.com/ZergsLaw/back-template/internal/logger"
	"github.com/ZergsLaw/back-template/internal/metrics"
	"github.com/ZergsLaw/back-template/internal/serve"
)

type (
	config struct {
		AuthKey string      `yaml:"auth_key"`
		Server  server      `yaml:"server"`
		DB      dbConfig    `yaml:"db"`
		Queue   queueConfig `yaml:"queue"`
	}
	server struct {
		Host string `yaml:"host"`
		Port ports  `yaml:"port"`
	}
	ports struct {
		GRPC   uint16 `yaml:"grpc"`
		Metric uint16 `yaml:"metric"`
	}
	dbConfig struct {
		MigrateDir string                 `yaml:"migrate_dir"`
		Driver     string                 `yaml:"driver"`
		Cockroach  connectors.CockroachDB `yaml:"cockroach"`
	}
	queueConfig struct {
		URLs     []string `yaml:"urls"`
		Username string   `yaml:"username"`
		Password string   `yaml:"password"`
	}
)

var (
	cfgFile  = &flags.File{DefaultPath: "config.yml", MaxSize: 1024 * 1024}
	logLevel = &flags.Level{Level: slog.LevelDebug}
)

const version = "v0.1.0"

func main() {
	flag.Var(cfgFile, "cfg", "path to config file")
	flag.Var(logLevel, "log_level", "log level")
	flag.Parse()

	log := buildLogger(logLevel.Level)
	grpclog.SetLoggerV2(grpchelper.NewLogger(log))

	appName := filepath.Base(os.Args[0])
	ctxParent := logger.NewContext(context.Background(), log.With(slog.String(logger.Version.String(), version)))
	ctx, cancel := signal.NotifyContext(ctxParent, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM)
	defer cancel()
	go forceShutdown(ctx)

	err := start(ctx, cfgFile, appName)
	if err != nil {
		log.Error("shutdown",
			slog.String(logger.Error.String(), err.Error()),
		)
		os.Exit(2)
	}
}

func start(ctx context.Context, cfgFile io.Reader, appName string) error {
	cfg := config{}
	err := yaml.NewDecoder(cfgFile).Decode(&cfg)
	if err != nil {
		return fmt.Errorf("yaml.NewDecoder.Decode: %w", err)
	}

	reg := prometheus.NewPedanticRegistry()

	return run(ctx, cfg, reg, appName)
}

func run(ctx context.Context, cfg config, reg *prometheus.Registry, namespace string) error {
	log := logger.FromContext(ctx)
	m := metrics.New(reg, namespace)

	r, err := repo.New(ctx, reg, namespace, repo.Config{
		Cockroach:  cfg.DB.Cockroach,
		MigrateDir: cfg.DB.MigrateDir,
		Driver:     cfg.DB.Driver,
	})
	if err != nil {
		return fmt.Errorf("repo.New: %w", err)
	}
	defer func() {
		err := r.Close()
		if err != nil {
			log.Error("close database connection", slog.String(logger.Error.String(), err.Error()))
		}
	}()

	q, err := queue.New(ctx, reg, namespace, queue.Config{
		URLs:     cfg.Queue.URLs,
		Username: cfg.Queue.Username,
		Password: cfg.Queue.Password,
	})
	if err != nil {
		return fmt.Errorf("queue.New: %w", err)
	}
	defer func() {
		err := q.Close()
		if err != nil {
			log.Error("close queue connection", slog.String(logger.Error.String(), err.Error()))
		}
	}()

	authModule := auth.New(cfg.AuthKey)
	module := app.New(r, authModule, idGenerator{}, q)
	grpcAPI := api.New(ctx, m, module, reg, namespace)

	err = serve.Start(
		ctx,
		serve.Metrics(log.With(slog.String(logger.Module.String(), "metric")), cfg.Server.Host, cfg.Server.Port.Metric, reg),
		serve.GRPC(log.With(slog.String(logger.Module.String(), "gRPC")), cfg.Server.Host, cfg.Server.Port.GRPC, grpcAPI),
		module.Process,
		q.Monitor,
		q.Process,
	)
	if err != nil {
		return fmt.Errorf("serve.Start: %w", err)
	}

	return nil
}

func buildLogger(level slog.Level) *slog.Logger {
	return slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{ //nolint:exhaustruct
				AddSource: true,
				Level:     level,
			},
		),
	)
}

var _ app.ID = &idGenerator{}

type idGenerator struct{}

// New implements app.ID.
func (idGenerator) New() uuid.UUID {
	return uuid.Must(uuid.NewV4())
}

func forceShutdown(ctx context.Context) {
	log := logger.FromContext(ctx)
	const shutdownDelay = 15 * time.Second

	<-ctx.Done()
	time.Sleep(shutdownDelay)

	log.Error("failed to graceful shutdown")
	os.Exit(2)
}
