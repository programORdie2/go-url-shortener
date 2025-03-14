package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
)

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var (
	shortURLs = make(map[string]string)
	mutex     = new(sync.Mutex)
)

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

	mutex.Lock()
	shortURLs[shortCode] = longURL
	mutex.Unlock()

	response := map[string]string{"short_url": "http://localhost:8080/redirect/" + shortCode}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func redirectURL(w http.ResponseWriter, r *http.Request) {
	shortCode := r.URL.Path[len("/redirect/"):]

	mutex.Lock()
	longURL, ok := shortURLs[shortCode]
	mutex.Unlock()

	if !ok {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, longURL, http.StatusSeeOther)
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	http.HandleFunc("/shorten", handleShorten)
	http.HandleFunc("/redirect/", redirectURL)

	fmt.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
