package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"
)

type Film struct {
	Title    string
	Director string
}

func main() {

	http.HandleFunc("/assets/css/bootstrap.min.css", serveCSS)
	http.HandleFunc("/", rootRoute)
	http.HandleFunc("/add-film/", addFilm)
	log.Fatal(http.ListenAndServe(":8000", nil))
	fmt.Println("server is running on port 8000")

}

func rootRoute(w http.ResponseWriter, r *http.Request) {
	//enableCors(&w)
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
	htmlStr := fmt.Sprintf("<li class='list-group-item bg-primary text-white'> %s - %s </li>", Title, Director)
	tmpl, _ := template.New("newfilm").Parse(htmlStr)
	tmpl.Execute(w, nil)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func serveCSS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css")
	http.ServeFile(w, r, "assets/css/bootstrap.min.css")
}
