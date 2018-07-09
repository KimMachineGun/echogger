# echogger [![Build Status](https://travis-ci.org/KimMachineGun/echogger.svg?branch=master)](https://travis-ci.org/KimMachineGun/echogger) [![Go Report Card](https://goreportcard.com/badge/github.com/KimMachineGun/echogger)](https://goreportcard.com/report/github.com/KimMachineGun/echogger)
Easy Swagger UI for your [Echo](https://echo.labstack.com/) API

## Installation
If you want to install echogger. 
```
go get -u github.com/KimMachineGun/echogger
```

## Example
```
package main

import (
	"net/http"

	"github.com/KimMachineGun/echogger"

	"github.com/labstack/echo"
)

const PORT = "8080"

func main() {
	e := echo.New()

	config := echogger.Config{
		Flavor:   "swagger",
		BasePath: "v1",
		SubPath:  "document",
		DocPath:  "./swagger.yaml",
		JSONName: "spec.json",
	}

	e.Use(echogger.MiddlewareWithConfig(config))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello world")
	})

	e.Start(":" + PORT)
}

```
If you use `echogger.Middleware()`, config is default value.  

> ### Default Value
>```
>Flavor:   "swagger"
>BasePath: "/"
>SubPath:  "docs"
>DocPath:  "./swagger.yaml"
>JSONName: "swagger.json"
>```


