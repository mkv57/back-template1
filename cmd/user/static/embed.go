// Package static contains embedded files.
package static

import (
	"embed"
	"io/fs"
)

var (
	//go:embed swagger-ui/index.html
	swaggerUI embed.FS
	// SwaggerUI contains overlay files for third_party.SwaggerUI.
	//nolint:gochecknoglobals // used for go embed.
	SwaggerUI, _ = fs.Sub(swaggerUI, "swagger-ui")
)
