package api

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"storage-gateway/application/api/handlers/get_object"
	"storage-gateway/application/api/handlers/put_object"
	"storage-gateway/application/api/middlewares"
	"storage-gateway/config"
	"storage-gateway/domain/services"
	"storage-gateway/internal/context-wrapper"
	"storage-gateway/internal/log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type API struct {
	server *echo.Echo
	config config.Config
	Addr   string
}

func NewApi(nps *services.NodePoolService, config config.Config) *API {
	return &API{
		server: echoServer(nps, config),
		config: config,
		Addr:   apiAddr(config.Api),
	}
}

func (a *API) Start() error {
	log.Infof("HTTP listener running on %s", a.Addr)
	return a.server.Start(a.Addr)
}

// Shutdown gracefully shuts down the API server
func (a *API) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(a.config.App.ShutdownTimeoutInSeconds)*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		log.Fatalf("could not shutdown API server gracefully with error %s", err)
	}
}

// echoServer sets up an Echo server with various middlewares for handling HTTP requests
func echoServer(nps *services.NodePoolService, config config.Config) *echo.Echo {
	e := echo.New()

	e.Logger.SetLevel(log.Lvl(config.App.LogLevel))

	e.Server.ReadHeaderTimeout = time.Duration(config.Api.ReadHeaderTimeoutInSeconds) * time.Second

	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		Skipper: middleware.DefaultSkipper,
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			log.Errort(context_wrapper.GetCorrelationID(c.Request().Context()), fmt.Sprintf("[PANIC RECOVER] Error: %s, Stack trace: %s", err, string(stack)))
			return err
		},
	}))

	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: time.Duration(config.Api.TimeoutInSeconds) * time.Second,
	}))

	e.Use(middleware.LoggerWithConfig(
		middleware.LoggerConfig{
			Output: e.Logger.Output(),
			Format: `{"time":"${time_rfc3339_nano}","trackID":"${id}","remote_ip":"${remote_ip}",` +
				`"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
				`"status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}"` +
				`,"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n",
		}))

	e.Use(middlewares.CorrelationID())

	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	e.GET("healthcheck", func(c echo.Context) error {
		return c.JSON(http.StatusOK, nil)
	})

	getObjectHandler := get_object.NewGetObjectHandler(services.NewGetObjectService(nps))
	e.GET("/object/:objectID", func(c echo.Context) error {
		return getObjectHandler.GetObject(c)
	})

	putObjectHandler := put_object.NewPutObjectHandler(services.NewPutObjectService(nps))
	e.PUT("/object/:objectID", func(c echo.Context) error {
		return putObjectHandler.PutObject(c)
	})

	return e
}

func apiAddr(cfg config.Api) string {
	return net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))
}
