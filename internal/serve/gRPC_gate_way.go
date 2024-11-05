package serve

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/mvrilo/go-redoc"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/cors"
	"google.golang.org/grpc"

	"github.com/ZergsLaw/back-template1/internal/grpchelper"
)

// GateWayConfig is config for building gRPC-Gateway proxy for WEB clients.
type GateWayConfig struct {
	FS             embed.FS
	Spec           string
	GRPCServerPort uint16
	Reg            *prometheus.Registry
	Namespace      string
	GRPCGWPattern  string                                                           // Pattern for http.ServeMux to serve grpc-gateway.
	DocsUIPattern  string                                                           // Pattern for http.ServeMux to serve Swagger UI.
	Register       func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error // Register gRPC server.
	DevMode        bool
}

// GRPCGateWay starts HTTP-proxy server for gRPC serer, for using gRPC endpoints from WEB.
func GRPCGateWay(log *slog.Logger, host string, port uint16, cfg GateWayConfig) func(context.Context) error {
	return func(ctx context.Context) error {
		const subsystem = "grpc_gateway_client"

		clientMetric := grpchelper.NewClientMetrics(cfg.Reg, cfg.Namespace, subsystem)

		conn, err := grpchelper.Dial(ctx,
			net.JoinHostPort(host, fmt.Sprintf("%d", cfg.GRPCServerPort)),
			log,
			clientMetric,
			[]grpc.UnaryClientInterceptor{},
			[]grpc.StreamClientInterceptor{},
			[]grpc.DialOption{},
		)
		if err != nil {
			return fmt.Errorf("grpc_helper.Dial: %w", err)
		}

		gw := runtime.NewServeMux()
		err = cfg.Register(ctx, gw, conn)
		if err != nil {
			return fmt.Errorf("cfg.Register: %w", err)
		}

		mux := http.NewServeMux()
		mux.Handle(cfg.GRPCGWPattern, noCache(corsAllowAll(gw)))
		mux.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK) // TODO: Add checking all external dependency.
		}))

		if cfg.DevMode {
			doc := redoc.Redoc{
				SpecPath: cfg.Spec,
				SpecFile: cfg.Spec,
				SpecFS:   &cfg.FS,
			}
			mux.Handle(cfg.DocsUIPattern, doc.Handler())
		}

		return HTTP(log, host, port, mux)(ctx)
	}
}

func noCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Expires", "0")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		next.ServeHTTP(w, r)
	})
}

func corsAllowAll(next http.Handler) http.Handler {
	handler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"*"},
		AllowCredentials: false,
	})

	return handler.Handler(next)
}
