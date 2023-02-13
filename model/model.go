package model

type Movie struct {
	Id     int    `json:"Id"`
	Title  string `json:"Title"`
	Rating int    `json:"Rating"`
	Notes  string `json:"Notes"`
	Year   string `json:"Year"`
	Poster string `json:"Poster"`
}

type NewMovie struct {
	Id     int    `json:"Id"`
	Title  string `json:"Title"`
	Rating string    `json:"Rating"`
	Notes  string `json:"Notes"`
}


type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []Movie
}

type Metadata struct {
	id           int
	movie_id     int
	poster       string
	release_date string
}
