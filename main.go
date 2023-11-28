package main

import (
	"log"
	"net/http"
)

func baseHandler(w http.ResponseWriter, r *http.Request) {
    path := r.URL.Path
    if path == "/" || path == "/home" {
	http.ServeFile(w, r, "templates/index.html")
    } else {
	http.ServeFile(w, r, "templates/error.html")
    }
}

func radioHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "templates/radio.html")
}

func clipboardHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "templates/clipboard.html")
}

func main() {
    http.HandleFunc("/", baseHandler)
    http.HandleFunc("/radio", radioHandler)
    http.HandleFunc("/clipboard", clipboardHandler)
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
    log.Print("Listening at localhost:8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
