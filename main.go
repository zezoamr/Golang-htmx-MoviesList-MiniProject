package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

type Film struct {
	Title    string
	Director string
}

type omdbResponse struct {
	Title  string `json:"Title"`
	Year   string `json:"Year"`
	Poster string `json:"Poster"`
	Plot   string `json:"Plot"`
}

type FilmCard struct {
	Title  string
	Year   string
	Poster string
	Plot   string
}

func main() {

	http.HandleFunc("/assets/css/bootstrap.min.css", serveCSS)
	http.HandleFunc("/add-film/", addFilm)
	http.HandleFunc("/get-film-info/", getFilmInfo)
	http.HandleFunc("/", rootRoute)
	log.Fatal(http.ListenAndServe(":8000", nil))
	fmt.Println("server is running on port 8000")

}

func rootRoute(w http.ResponseWriter, r *http.Request) {
	//enableCors(&w)
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	films := map[string][]Film{
		"Films": {
			{Title: "The Godfather", Director: "Francis Ford Coppola"},
			{Title: "Blade Runner", Director: "Ridley Scott"},
			{Title: "The Thing", Director: "John Carpenter"},
		},
	}
	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, films)
}

func addFilm(w http.ResponseWriter, r *http.Request) {
	//enableCors(&w)
	time.Sleep(1 * time.Second)
	Title := r.PostFormValue("title")
	Director := r.PostFormValue("director")
	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.ExecuteTemplate(w, "film-list-element", Film{Title: Title, Director: Director})
}

func getFilmInfo(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	if title == "" {
		http.Error(w, "Missing title parameter", http.StatusBadRequest)
		return
	}

	film, err := fetchFilm(title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.ExecuteTemplate(w, "film-info-element", FilmCard{Poster: film.Poster, Title: film.Title, Year: film.Year, Plot: film.Plot})
}

func fetchFilm(title string) (*omdbResponse, error) {
	/*resp, err := http.Get(fmt.Sprintf("https://www.omdbapi.com/?t=%s", title))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var film omdbResponse
	err = json.Unmarshal(body, &film)
	if err != nil {
		return nil, err
	}
	return &film, nil*/

	// Seed the random number generator
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate a random integer between 0 and 99
	num := rng.Intn(100) + 1900

	// Convert the integer to a string
	numStr := strconv.Itoa(num)

	return &omdbResponse{
		Title:  "placeholder - " + numStr,
		Year:   numStr,
		Poster: "https://i.ytimg.com/vi/lgwT6tDniko/maxresdefault.jpg",
		Plot:   "place holder wow what a fantastic movie",
	}, nil
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func serveCSS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css")
	http.ServeFile(w, r, "assets/css/bootstrap.min.css")
}
