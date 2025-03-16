package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	_ "modernc.org/sqlite"
)

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type RequestData struct {
	URL string `json:"url"`
}

func generateShortCode() string {
	code := make([]byte, 6)
	for i := range code {
		code[i] = chars[rand.Intn(len(chars))]
	}

	return string(code)
}

func handleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data RequestData

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	longURL := data.URL
	shortCode := generateShortCode()

	saveShortURL(shortCode, longURL)

	response := map[string]string{"short_url": "http://localhost:8080/redirect/" + shortCode}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

var db *sql.DB

func redirectURL(w http.ResponseWriter, r *http.Request) {
	shortCode := r.URL.Path[len("/redirect/"):]

	var longURL string
	row := db.QueryRow("SELECT long_url FROM short_urls WHERE short_url = ?", shortCode)
	err := row.Scan(&longURL)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, longURL, http.StatusSeeOther)
}

func saveShortURL(shortCode, longURL string) {
	// Prevent sql injection
	// TODO

	_, err := db.Exec("INSERT INTO short_urls (short_url, long_url) VALUES (?, ?)", shortCode, longURL)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var err error
	db, err = sql.Open("sqlite", "shortener.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the table if it doesn't exist
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS short_urls (short_url TEXT PRIMARY KEY, long_url TEXT)")
	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	http.HandleFunc("/shorten", handleShorten)
	http.HandleFunc("/redirect/", redirectURL)

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
