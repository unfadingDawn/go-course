package main

import (
	movieadapter "homework3/internal/adapter/movie"
	"homework3/internal/config"
	moviehttp "homework3/internal/handler/http"
	"homework3/internal/usecase"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	c := movieadapter.NewMovieClient(cfg.MovieBaseURL, cfg.UpstreamTimeout)
	uc := usecase.NewMovieService(c)

	r := moviehttp.NewRouter(uc, cfg.OperationTimeout)

	if err := http.ListenAndServe(cfg.HTTPAddr, r); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
