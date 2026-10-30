package model

type Movie struct {
	Title    string `json:"title"`
	Director string `json:"director"`
	ID       int    `json:"id"`
	Year     int    `json:"year"`
}
