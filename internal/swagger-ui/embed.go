// Package swaggerui contains embedded files.
package swaggerui

import (
	"embed"
	"io/fs"
)

var (
	//go:embed dist/favicon-16x16.png
	//go:embed dist/favicon-32x32.png
	//go:embed dist/oauth2-redirect.html
	//go:embed dist/swagger-ui-bundle.js
	//go:embed dist/swagger-ui.css
	//go:embed dist/swagger-ui-standalone-preset.js
	swaggerUI embed.FS
	// SwaggerUI contains static files for Swagger UI required by
	// our ../swagger-ui/index.html.
	//nolint:gochecknoglobals // used for go embed.
	SwaggerUI, _ = fs.Sub(swaggerUI, "dist")
)
