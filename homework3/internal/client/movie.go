package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"homework3/internal/model"
	"net/http"
	"time"
)

var ErrNoMovies = errors.New("films not found")

type MovieClient struct {
	httpClient       *http.Client
	BaseURL          string
	UpstreamTimeout  time.Duration
	OperationTimeout time.Duration
}

func NewMovieClient(baseURL string, upstreamTimeout, operationTimeout time.Duration) *MovieClient {
	return &MovieClient{&http.Client{}, baseURL, upstreamTimeout, operationTimeout}
}

func (c *MovieClient) GetMovie(ctx context.Context, id int) (model.Movie, error) {
	url := fmt.Sprintf(
		"%s/%d/info.0.json",
		c.BaseURL,
		id,
	)
	ctx, cancel := context.WithTimeout(ctx, c.UpstreamTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return model.Movie{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return model.Movie{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.Movie{}, fmt.Errorf(
			"unexpected status: %d",
			resp.StatusCode,
		)
	}

	var movie model.Movie

	if err := json.NewDecoder(resp.Body).Decode(&movie); err != nil {
		return model.Movie{}, err
	}

	return movie, nil
}

func (c *MovieClient) ExistsMovie(ctx context.Context, id int) (bool, error) {
	url := fmt.Sprintf(
		"%s/%d/info.0.json",
		c.BaseURL,
		id,
	)
	ctx, cancel := context.WithTimeout(ctx, c.UpstreamTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return false, err
	}

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}
	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}
	return false, fmt.Errorf("error code: %d", resp.StatusCode)
}

func (c *MovieClient) GetLastMovie(ctx context.Context) (model.Movie, error) {
	low, high := 0, 1
	for {
		exists, err := c.ExistsMovie(ctx, high)

		if err != nil {
			return model.Movie{}, err
		}

		if !exists {
			break
		}

		low = high
		high *= 2
	}
	for high-low > 1 {
		mid := low + (high-low)/2

		exists, err := c.ExistsMovie(ctx, mid)
		if err != nil {
			return model.Movie{}, err
		}

		if exists {
			low = mid
		} else {
			high = mid
		}
	}
	if high == 1 {
		return model.Movie{}, ErrNoMovies
	}
	movie, err := c.GetMovie(ctx, low)
	if err != nil {
		return model.Movie{}, err
	}
	return movie, err
}
