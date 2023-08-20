package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

type Film struct {
	Title    string
	Director string
}

func main() {

	http.HandleFunc("/assets/css/bootstrap.min.css", serveCSS)
	http.HandleFunc("/", rootRoute)
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

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func serveCSS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css")
	http.ServeFile(w, r, "assets/css/bootstrap.min.css")
}
