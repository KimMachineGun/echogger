# echogger
Easy Swagger UI for your Echo API

## Require
- [golang/dep](https://github.com/golang/dep)

## Installation
If you want to install echogger. 
```
go get -u github.com/DEATH-TROOPER/echogger
```
Then you can see the echogger directory like this.
```
src
├── github.com
|    └── DEATH-TROOPER
|         └── echogger
|              ├── echogger.go
|              └── templates.go
└── echo-server
     ├── main.go
     └── swagger.yml
     
```
Echogger need other libraries. So you need to install them with `dep ensure`  

## Example
```
package main

import (
	"net/http"

	"github.com/DEATH-TROOPER/echogger"
	"github.com/labstack/echo"
)

const PORT = "8080"

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello world")
	})

	config := echogger.SwaggerConfig{
		Flavor:   "swagger",
		BasePath: "v1",
		SubPath:  "docs",
		DocPath:  "./swagger.yml",
	}

	echogger.StartWithConfig(e, ":"+PORT, config)
}
```
If you use `echogger.Start()`, config is default value.  

### Default
>```
>Flavor:   "redoc"
>BasePath: "/"
>SubPath:  "docs"
>DocPath:  "./swagger.yml"
>JSONName: "swagger.json"
>```


