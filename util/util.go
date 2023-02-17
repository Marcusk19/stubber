// module to map metadata to movie
package util

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"

	"github.com/Marcusk19/stubber/data"
	tmdb "github.com/cyruzin/golang-tmdb"
)

var tmdbClient *tmdb.Client

func connect() {

	client, err := tmdb.Init(os.Getenv("TMDB_KEY"))
	if err != nil {
		log.Print("[util] [ERROR] " + err.Error())
	}
	tmdbClient = client
}

func GenerateMetadata(name string, movie_id int) int64 {
	// use movie title to search with tmdb and create entries in metadata table
	connect()
	result, err := tmdbClient.GetSearchMovies(name, nil)
	HandleError(err)

	movie := result.SearchMoviesResults.Results[0]

	id := movie.ID
	release := movie.ReleaseDate
	poster := movie.PosterPath

	db := data.Connect()
	defer db.Close()

	_, db_err := db.Query("INSERT INTO metadata (id, movie_id, release_date, poster) VALUES($1, $2, $3, $4)", id, movie_id, release, poster)
	HandleError(db_err)

	return id
}

func GetReleaseDate(id int64) string {
	// lookup release date in metadata table and return it
	var releasedate string
	db := data.Connect()
	defer db.Close()

	rows := db.QueryRow("SELECT release_date FROM metadata WHERE id=$1", id)
	err := rows.Scan(&releasedate)
	HandleError(err)

	return releasedate
}

func GetPosterPath(movie_id int) string {
	// lookup poster path in metadata table and return it
	var posterpath string
	db := data.Connect()
	defer db.Close()

	rows := db.QueryRow("SELECT poster FROM metadata WHERE movie_id=$1", movie_id)
	err := rows.Scan(&posterpath)
	HandleError(err)

	return posterpath
}

func HandleError(err error) bool {
	// handle error - returns true if error is pressent, false otherwise
	_, file, line, _ := runtime.Caller(1)
	if err != nil {
		log.Print("[ERROR] " + file + " line " + strconv.Itoa(line) + " " + err.Error())
		return true
	}
	return false
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	// helper function for uploading a file
	log.Print("Uploading file...")

	// Parse multipart form
	// 10 << 20 specifies max upload of 10 MB files
	r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile("myFile")
	HandleError(err)
	defer file.Close()
	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile("temp-data", "data-*.csv")
	HandleError(err)
	defer tempFile.Close()

	// read contents of file into a byte array
	fileBytes, err := ioutil.ReadAll(file)
	HandleError(err)
	// write this file to temp file
	tempFile.Write(fileBytes)

	log.Print("Successfully uploaded file")
}
