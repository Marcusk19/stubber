package controller

import (
	"database/sql"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/Marcusk19/stubber/data"
	"github.com/Marcusk19/stubber/model"
	"github.com/Marcusk19/stubber/util"
	"github.com/gorilla/mux"
)

const (
	by_id     = 1
	by_rating = 2
	by_title  = 3
)

// Function to return all movies from the database ordered by ID
func MoviesById(w http.ResponseWriter, r *http.Request) {
	var response model.Response

	arrMovie := AllMovies(by_id)

	response.Status = 200
	response.Message = "Success"
	response.Data = arrMovie
	log.Print("[INFO]Fetched list of movies")

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

// Same as MoviesById but ordered by rating
func MoviesByRating(w http.ResponseWriter, r *http.Request) {
	var response model.Response

	arrMovie := AllMovies(by_rating)

	response.Status = 200
	response.Message = "Success"
	response.Data = arrMovie
	log.Print("[INFO]Fetched list of movies")

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

// returns all movies in specified order
func AllMovies(order int) []model.Movie {

	var movie model.Movie
	var arrMovie []model.Movie
	var rows *sql.Rows
	var err error

	db := data.Connect()
	defer db.Close()

	switch order {
	case by_id:
		rows, err = db.Query("SELECT * FROM movies ORDER BY id ASC")
	case by_rating:
		rows, err = db.Query("SELECT * FROM movies ORDER BY rating DESC")
	case by_title:
		rows, err = db.Query("SELECT * FROM movies ORDER BY title ASC")
	default:
		rows, err = db.Query("SELECT * FROM movies ORDER BY rating DESC")
	}

	util.HandleError(err)

	for rows.Next() {
		err = rows.Scan(&movie.Id, &movie.Title, &movie.Rating, &movie.Notes, &movie.Year)
		movie.Poster = util.GetPosterPath(movie.Id)
		if err != nil {
			util.HandleError(err)
		} else {
			arrMovie = append(arrMovie, movie)
		}

	}

	return arrMovie
}

// Function to insert/create a new movie into the db
func InsertMovie(w http.ResponseWriter, r *http.Request) {
	// preflight cors
	if r.Method == "OPTIONS" {
		handleCors(w)
		return
	}
	requestBody, error := ioutil.ReadAll(r.Body)
	util.HandleError(error)
	var movie model.NewMovie
	json.Unmarshal(requestBody, &movie)

	db := data.Connect()
	defer db.Close()

	row := db.QueryRow("SELECT id FROM movies ORDER BY id DESC limit 1") // sql query to get highest value in id col.
	var id int
	countError := row.Scan(&id)
	if countError == sql.ErrNoRows {
		id = 0
	}
	id++
	title := movie.Title
	// string to int conversion for rating
	rating_str := movie.Rating
	rating, _ := strconv.Atoi(rating_str)

	notes := movie.Notes

	log.Print("Received movie with content: ")
	log.Print(movie)

	if title == "" {
		log.Print("Request failed to specify title")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if notes == "" {
		notes = "none"
	}

	_, err := db.Exec("INSERT INTO movies(id, title, rating, notes) VALUES($1, $2, $3, $4)", id, title, rating, notes)
	util.HandleError(err)

	metadata_id := util.GenerateMetadata(title, id)
	log.Print(metadata_id)

	if metadata_id == 0 {
		log.Print("unable to find metadata for " + movie.Title)
		return
	}

	year := util.GetReleaseDate(metadata_id)
	_, err = db.Exec("UPDATE movies SET year=$1 WHERE id=$2", year, id)
	util.HandleError(err)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(movie)
}

// function to update the data associated with a movie (given movie ID)
func UpdateMovie(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)

	var updateMovie model.Movie
	json.Unmarshal(reqBody, &updateMovie)

	db := data.Connect()
	defer db.Close()

	id := updateMovie.Id
	title := updateMovie.Title
	rating := updateMovie.Rating

	_, err := db.Exec("UPDATE movies SET title=$1, rating=$2 WHERE id=$3", title, rating, id)

	checkErr(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updateMovie)

}

// function to remove a movie from database
func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	// cors handler
	var response model.Response
	if r.Method == "OPTIONS" {
		handleCors(w)
		return
	}
	db := data.Connect()
	defer db.Close()
	params := mux.Vars(r)
	id := params["id"]
	log.Print("Attempting to delete movie " + id)

	_, err := db.Exec("DELETE FROM movies WHERE id=$1", id)
	if err != nil {
		log.Print(err)
		return
	}

	response.Status = http.StatusOK
	response.Message = "sucessfully deleted"

	log.Print("[INFO] DELETED movie: " + id)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

// function to get movie by id
func GetMovie(w http.ResponseWriter, r *http.Request) {
	var movie model.Movie
	var response model.Response
	var arrMovie []model.Movie
	db := data.Connect()
	defer db.Close()

	// obtain url parameter (aka movie id)
	params := mux.Vars(r)
	id := params["id"]

	res := db.QueryRow("SELECT * FROM movies WHERE id=$1", id)
	err := res.Scan(&movie.Id, &movie.Title, &movie.Rating)
	util.HandleError(err)
	if err != nil {
		// catch error if movie doesn't exist in table
		log.Print(err)
		movie.Id = 0
		movie.Title = ""
		movie.Rating = 0
		movie.Notes = ""
		movie.Year = ""
	}
	arrMovie = append(arrMovie, movie)

	response.Status = 200
	response.Message = "Success"
	response.Data = arrMovie
	log.Print("[INFO]Fetched movie by id")

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)

}

// Upload file to /tmp/uploaded_rankings directory
func UploadFile(w http.ResponseWriter, r *http.Request) {
	var response model.Response

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	// helper function for uploading a file
	log.Print("Uploading file...")

	// Parse multipart form
	// 10 << 20 specifies max upload of 10 MB files
	r.ParseMultipartForm(900 << 20)
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		log.Println("Error receiving the file")
		log.Println(err)
		return
	}
	defer file.Close()
	log.Printf("File to upload: %+v\n", handler.Filename)
	log.Printf("File Size: %+v\n", handler.Size)
	log.Printf("MIME Header: %+v\n", handler.Header)

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile("/tmp/uploaded_rankings", "data-*.csv")
	if err != nil {
		log.Print(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer tempFile.Close()

	log.Print("Created tmp file for writing")
	// write this file to temp file
	_, err = io.Copy(tempFile, file)
	if err != nil {
		log.Print("[ERROR] problem copying")
	}

	log.Print("Successfully uploaded file")
	response.Status = http.StatusCreated
	response.Message = "Successfully uploaded file"

	json.NewEncoder(w).Encode(response)
}

func checkErr(err error) {
	if err != nil {
		log.Print(err)
		return
	}
}

func handleCors(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
}
