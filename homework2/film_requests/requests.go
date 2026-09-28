package film_requests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Film struct {
	Title    string `json:"title"`
	Director string `json:"director"`
	ID       int    `json:"id"`
	Year     int    `json:"year"`
}

func (f Film) String() string {
	var empty Film
	if empty != f {
		return fmt.Sprintf("%d - %s - %d - %s", f.ID, f.Title, f.Year, f.Director)
	}
	return ""
}

func GetFilm(ctx context.Context, id int) (Film, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://homeworksite.site/"+strconv.Itoa(id)+"/info.0.json", nil)
	if err != nil {
		return Film{}, fmt.Errorf("create request for film with id %d failed: %v", id, err)
	}
	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return Film{}, fmt.Errorf("create request for film with id %d failed: %v", id, err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return Film{}, fmt.Errorf("server returned status %d for film with id %d", res.StatusCode, id)
	}
	var film Film
	if err := json.NewDecoder(res.Body).Decode(&film); err != nil {
		return Film{}, fmt.Errorf("getting film with id %d failed: %v", id, err)
	}
	return film, nil
}
