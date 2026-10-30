package main

import (
	"homework3/internal/client"
	"homework3/internal/config"
	moviehttp "homework3/internal/transport/http"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	c := client.NewMovieClient(cfg.MovieBaseURL, cfg.UpstreamTimeout, cfg.OperationTimeout)

	r := moviehttp.NewRouter(c)

	if err := http.ListenAndServe(cfg.HTTPAddr, r); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
