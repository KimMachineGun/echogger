package echogger

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/labstack/echo"
)

func TestMiddlewareWithConfig(t *testing.T) {
	assert := assert.New(t)
	e := echo.New()

	testConfigs := []Config{
		{
			DocPath:  "testSpec.yaml",
			BasePath: "base",
			Flavor:   "",
			SubPath:  "docs",
			JSONName: "docs.json",
			NoUI:     false,
		},
		{
			DocPath: "testSpec.yaml",
		},
	}

	for _, config := range testConfigs {

		config.EnsureDefaults()
		middleware := MiddlewareWithConfig(config)

		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()

		// Test API Docs
		c := e.NewContext(req, rec)
		c.SetPath(path.Join(config.BasePath, config.SubPath))

		h := middleware(echo.NotFoundHandler)

		if assert.NoError(h(c)) {
			assert.Equal(http.StatusOK, rec.Code)
			assert.NotEmpty(rec.Body.String())
		}

		// Test Json API Specs
		req = httptest.NewRequest(echo.GET, "/", nil)
		rec.Flush()

		c = e.NewContext(req, rec)
		c.SetPath(path.Join(config.BasePath, config.JSONName))

		if assert.NoError(h(c)) {
			assert.Equal(http.StatusOK, rec.Code)
			assert.NotEmpty(rec.Body.String())
		}
	}
}

func TestMiddlewareWithoutFile(t *testing.T) {
	assert := assert.New(t)

	defer func() {
		_, ok := recover().(error).(*os.PathError)
		assert.True(ok)
	}()

	Middleware()
}
