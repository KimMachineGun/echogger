package echogger

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"path"

	"github.com/go-openapi/loads"
	redocmiddleware "github.com/go-openapi/runtime/middleware"

	"github.com/labstack/echo"
)

// Config configures the echogger middlewares
type Config struct {
	DocPath  string
	BasePath string
	Flavor   string
	SubPath  string
	JSONName string
	NoUI     bool
}

// Middleware creates a middleware to serve a documentation site
func Middleware() echo.MiddlewareFunc {
	return MiddlewareWithConfig(Config{})
}

// MiddlewareWithConfig creates a middleware with options to serve a documentation site
func MiddlewareWithConfig(config Config) echo.MiddlewareFunc {
	config.EnsureDefaults()

	specDoc, err := loads.Spec(config.DocPath)
	if err != nil {
		panic(err)
	}

	b, err := json.MarshalIndent(specDoc.Spec(), "", "	")
	if err != nil {
		panic(err)
	}

	html := ""

	if !config.NoUI {
		if config.Flavor == "redoc" {
			redocOpts := redocmiddleware.RedocOpts{
				BasePath: config.BasePath,
				SpecURL:  path.Join(config.BasePath, config.JSONName),
				Path:     config.SubPath,
			}

			redocOpts.EnsureDefaults()

			tmpl := template.Must(template.New("redoc").Parse(redocTemplate))

			buf := bytes.NewBuffer(nil)
			_ = tmpl.Execute(buf, redocOpts)
			html = buf.String()
		} else {
			swaggerOpts := map[string]interface{}{
				"SpecURL":  path.Join(config.BasePath, config.JSONName),
				"CSS":      swaggerCSS,
				"BundleJS": swaggerBundleJS,
				"PresetJS": swaggerPreset,
			}

			tmpl := template.Must(template.New("swagger").Parse(swaggerTemplate))

			buf := bytes.NewBuffer(nil)
			_ = tmpl.Execute(buf, swaggerOpts)
			html = buf.String()
		}
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			switch c.Path() {
			case path.Join(config.BasePath, config.JSONName):
				c.Response().Header().Set("Access-Control-Allow-Origin", "*")
				return c.String(http.StatusOK, string(b))
			case path.Join(config.BasePath, config.SubPath):
				return c.HTML(http.StatusOK, html)
			default:
				return next(c)
			}
		}
	}
}

// EnsureDefaults in case some options are missing
func (c *Config) EnsureDefaults() {
	if c.DocPath == "" {
		c.DocPath = "./swagger.yaml"
	}

	if c.BasePath == "" || !path.IsAbs(c.BasePath) {
		c.BasePath = "/" + c.BasePath
	}

	if c.JSONName == "" {
		c.JSONName = "swagger.json"
	} else if path.Ext(c.JSONName) != ".json" {
		panic(errors.New("JSONName must have .json extension"))
	}

	if c.Flavor == "" {
		c.Flavor = "swagger"
	} else if c.Flavor != "redoc" && c.Flavor != "swagger" {
		panic(errors.New("Flavor must be redoc or swagger"))
	}

	if c.SubPath == "" {
		c.SubPath = "docs"
	}
}
