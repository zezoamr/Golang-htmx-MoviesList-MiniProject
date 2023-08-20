package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/", rootRoute)
	log.Fatal(http.ListenAndServe(":8000", nil))
	fmt.Println("server is running on port 8000")

}

func rootRoute(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello world \n")
	io.WriteString(w, "some request details "+r.Host+" "+r.Method+" "+r.Proto)
}
