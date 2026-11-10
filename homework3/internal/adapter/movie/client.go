package movie

import (
	"context"
	"encoding/json"
	"fmt"
	"homework3/internal/domain"
	"net/http"
	"time"
)

type MovieClient struct {
	httpClient      *http.Client
	BaseURL         string
	UpstreamTimeout time.Duration
}

type movieResponse struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Director string `json:"director"`
	Year     int    `json:"year"`
}

func NewMovieClient(baseURL string, upstreamTimeout time.Duration) *MovieClient {
	return &MovieClient{&http.Client{}, baseURL, upstreamTimeout}
}

func (c *MovieClient) GetMovie(ctx context.Context, id int) (domain.Movie, error) {
	url := fmt.Sprintf("%s/%d/info.0.json", c.BaseURL, id)
	ctx, cancel := context.WithTimeout(ctx, c.UpstreamTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return domain.Movie{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return domain.Movie{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return domain.Movie{}, domain.ErrMovieNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return domain.Movie{}, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var response movieResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return domain.Movie{}, err
	}
	return domain.Movie{
		ID:       response.ID,
		Title:    response.Title,
		Director: response.Director,
		Year:     response.Year,
	}, nil
}
