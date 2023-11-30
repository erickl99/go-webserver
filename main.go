package main

import (
	"bufio"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

type entry struct {
	lexeme  string
	reading string
}

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
	fp, err := os.Open("data/clipboard")
	if err != nil {
		log.Print("Error occurred reading the file!")
		http.ServeFile(w, r, "templates/error.html")
		return
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	var clipboardList []string
	for scanner.Scan() {
		clipboardList = append(clipboardList, scanner.Text())
	}
	tp, err := template.ParseFiles("templates/clipboard.html")
	if err != nil {
		log.Print("Error occurred parsing the template!")
		http.ServeFile(w, r, "templates/error.html")
		return
	}
	tp.Execute(w, clipboardList)
}

func reviewHandler(w http.ResponseWriter, r *http.Request) {
	fp, err := os.Open("data/decks")
	if err != nil {
		log.Print("An error occurred opening the file!")
		http.ServeFile(w, r, "templates/error.html")
		return
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	decks := make(map[string][]entry)
	current_deck := ""
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) > 0 && line[0] == '[' {
			current_deck = line[1 : len(line)-1]
			var entries []entry
			decks[current_deck] = entries
		} else if current_deck != "" {
			tokens := strings.Split(line, " ")
			if len(tokens) == 2 {
				ent := entry{lexeme: tokens[0], reading: tokens[1]}
				entries := decks[current_deck]
				entries = append(entries, ent)
				decks[current_deck] = entries
			}
		}
	}
	tp, err := template.ParseFiles("templates/review.html")
	if err != nil {
		log.Print("An error occurred reading the template!")
		http.ServeFile(w, r, "templates/error.html")
		return
	}
	tp.Execute(w, decks)
}

func main() {
	http.HandleFunc("/", baseHandler)
	http.HandleFunc("/radio", radioHandler)
	http.HandleFunc("/clipboard", clipboardHandler)
	http.HandleFunc("/review", reviewHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	log.Print("Listening at localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
