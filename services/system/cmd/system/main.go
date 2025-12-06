package main

//go:generate oapi-codegen --package=api --generate types,skip-prune,spec -o ./../../api/dummy/dummy-openapi.gen.go ./../../api/dummy/dummy-openapi.yaml
//go:generate oapi-codegen --package=api --config=./../../api/config.yaml -o ./../../api/system-openapi.gen.go ./../../api/system-openapi.yaml

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"

	openapi "github.com/RaviChandarkolla/go-echo-containerisation/services/system/api"

	api2 "github.com/RaviChandarkolla/go-echo-containerisation/services/system/internal/api"
)

var (
	port = flag.Int("port", GetEnvIntParam("SYSTEM_PORT", 8001), "System port")
)

// GetEnvIntParam : return integer environmental param if exists, otherwise return default
func GetEnvIntParam(param string, dflt int) int {
	if v, exists := os.LookupEnv(param); exists {
		i, err := strconv.Atoi(v)
		if err != nil {
			return dflt
		}
		return i
	}
	return dflt
}

func main() {
	openapi3.DefineStringFormat("uuid", openapi3.FormatOfStringForUUIDOfRFC4122)
	swagger, err := openapi.GetSwagger()
	if err != nil {
		fmt.Printf("error loading openapi spec: %s\n", err)
		os.Exit(1)
	}

	e := echo.New()
	// validate requests against schema
	e.Use(middleware.OapiRequestValidator(swagger))
	dummyAPI := api2.NewServer()
	openapi.RegisterHandlers(e, dummyAPI)

	serverAddr := fmt.Sprintf(":%d", *port)
	wg := sync.WaitGroup{}

	// start API server
	go func() {
		defer wg.Done()
		e.Logger.Fatal(e.Start(serverAddr))
	}()

}
