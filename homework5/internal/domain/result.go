package domain

type MovieLoadError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	ID      int64  `json:"id"`
}

type LoadMoviesResult struct {
	Movies []Movie          `json:"movies"`
	Errors []MovieLoadError `json:"errors"`
}
