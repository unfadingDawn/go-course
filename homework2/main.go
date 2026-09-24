package main

import (
	"context"
	"flag"
	"fmt"
	"homework2/film_requests"
	"log"
	"os/signal"
	"syscall"
	"time"
)

type Request struct {
	ID, Timeout int
}

type Result struct {
	Err  error
	Film film_requests.Film
}

func worker(ctx context.Context, _ int, jobs <-chan Request, results chan<- Result) {
	for args := range jobs {
		ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(args.Timeout)*time.Second)
		film, err := film_requests.GetFilm(ctxTimeout, args.ID)
		results <- Result{err, film}
		cancel()
	}
}

func main() {
	var from, to int
	flag.IntVar(&from, "from", -1, "id of first film")
	flag.IntVar(&to, "to", -1, "id of last film")
	workers := flag.Int("workers", 10, "quantity of workers in worker pool")
	timeout := flag.Int("timeout", 5, "timeout in seconds")
	flag.Parse()
	if from == -1 || to == -1 {
		flag.Usage()
		return
	}
	if from > to {
		log.Fatalf("'to' must be greater than 'from'")
	}
	numFilms := to - from + 1
	jobs := make(chan Request, numFilms)
	results := make(chan Result, numFilms)

	ctx, stop := signal.NotifyContext(
		context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	for w := 1; w <= *workers; w++ {
		go worker(ctx, w, jobs, results)
	}
	go func() {
		for i := from; i <= to; i++ {
			jobs <- Request{i, *timeout}
		}
		close(jobs)
	}()
	for i := 0; i < numFilms; i++ {
		result := <-results
		if result.Err == nil {
			fmt.Println(result.Film)
		} else {
			log.Println(result.Err)
		}
	}
}
