package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os/signal"
	"syscall"

	"storage-gateway/application/api"
	"storage-gateway/config"
	"storage-gateway/domain/services"
	"storage-gateway/infrastructure/discovery-service"
	"storage-gateway/internal/log"
)

func main() {
	configFile := flag.String("conf", "config/config.local.json", "Config file path")
	flag.Parse()

	shutdownCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	appConfig, err := config.Read(*configFile)
	if err != nil {
		log.Fatalf("could not read config file with error %s", err)
	}

	log.SetupLogging(appConfig.App.LogLevel)

	dds, err := discovery_service.NewDockerDiscoveryService()
	if err != nil {
		log.Fatal(err.Error())
	}

	nps := services.NewNodePoolService(dds)

	go func() {
		if err = nps.StartRefreshingNodes(); err != nil {
			log.Fatalf("could not start refresh nodes scheduler with error %s", err)
		}
	}()

	gateway := api.NewApi(nps, *appConfig)

	go runApiHandler(gateway)

	<-shutdownCtx.Done()
	gateway.Shutdown()
	nps.StopRefreshingNodes()
}

func runApiHandler(gateway *api.API) {
	if err := gateway.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("could not start API server with error %s", err)
	}
}
