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

type (
	SwaggerConfig struct {
		DocPath  string
		BasePath string
		Flavor   string
		SubPath  string
		JSONName string
		NoUI     bool
	}
)

func Start(e *echo.Echo, address string) error {
	return StartWithConfig(e, address, SwaggerConfig{})
}

func StartWithConfig(e *echo.Echo, address string, config SwaggerConfig) error {
	if config.DocPath == "" {
		config.DocPath = "./swagger.yml"
	}

	specDoc, err := loads.Spec(config.DocPath)
	if err != nil {
		panic(err)
	}
	b, err := json.MarshalIndent(specDoc.Spec(), "", "	")
	if err != nil {
		panic(err)
	}

	basePath := config.BasePath
	if basePath == "" || !path.IsAbs(basePath) {
		basePath = "/" + basePath
	}

	jsonName := config.JSONName
	if jsonName == "" {
		jsonName = "swagger.json"
	} else if path.Ext(jsonName) != ".json" {
		panic(errors.New("JsonName must have .json extension."))
	}

	flavor := config.Flavor
	if flavor == "" {
		flavor = "redoc"
	} else if flavor != "redoc" && flavor != "swagger" {
		panic(errors.New("Flavor must be redoc or swagger."))
	}

	subPath := config.SubPath
	if subPath == "" {
		subPath = "docs"
	}

	e.GET(path.Join(config.BasePath, jsonName), func(c echo.Context) error {
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		return c.String(http.StatusOK, string(b))
	})

	if !config.NoUI {
		if flavor == "redoc" {
			redocOpts := redocmiddleware.RedocOpts{
				BasePath: basePath,
				SpecURL:  path.Join(basePath, jsonName),
				Path:     subPath,
			}

			redocOpts.EnsureDefaults()

			tmpl := template.Must(template.New("redoc").Parse(redocTemplate))

			buf := bytes.NewBuffer(nil)
			_ = tmpl.Execute(buf, redocOpts)
			html := buf.String()

			e.GET(path.Join(basePath, subPath), func(c echo.Context) error {
				return c.HTML(http.StatusOK, html)
			})
		} else {
			tmpl := template.Must(template.New("redoc").Parse(swaggerTemplate))

			swaggerOpts := map[string]interface{}{
				"SpecURL":  path.Join(basePath, jsonName),
				"CSS":      swaggerCSS,
				"BundleJS": swaggerBundleJS,
				"PresetJS": swaggerPreset,
			}

			buf := bytes.NewBuffer(nil)
			_ = tmpl.Execute(buf, swaggerOpts)
			html := buf.String()

			e.GET(path.Join(basePath, subPath), func(c echo.Context) error {
				return c.HTML(http.StatusOK, html)
			})
		}
	}

	return e.Start(address)
}
