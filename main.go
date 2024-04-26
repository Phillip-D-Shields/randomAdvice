package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

const (
	ADVICE_URL = "https://api.adviceslip.com/advice"
)

// Advice represents the structure of the advice data
type Advice struct {
	Slip struct {
		ID     int    `json:"id"`
		Advice string `json:"advice"`
	} `json:"slip"`
	Count int
}

var db *sql.DB

func main() {
	// database
	var err error
	db, err = sql.Open("sqlite3", "advice.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS advice (
			id INTEGER PRIMARY KEY,
			advice TEXT,
			count INTEGER DEFAULT 0
	)`)
	if err != nil {
		log.Fatal(err)
	}

	// routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})
	http.HandleFunc("/advice", getCheekyScoops)

	// serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// server
	log.Println("Server is running on http://localhost:8080")

	// ! =======================================
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getCheekyScoops(w http.ResponseWriter, r *http.Request) {
	// hit advice url
	resp, err := http.Get(ADVICE_URL)
	if err != nil {
		log.Println("Error fetching advice:", err)
		http.Error(w, "Oops! Something went wrong.", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		http.Error(w, "Oops! Something went wrong.", http.StatusInternalServerError)
		return
	}

	// parse response body
	var advice Advice
	err = json.Unmarshal(body, &advice)
	if err != nil {
		log.Println("Error parsing JSON:", err)
		http.Error(w, "Oops! Something went wrong.", http.StatusInternalServerError)
		return
	}

	var count int
	err = db.QueryRow("SELECT count FROM advice WHERE id = ?", advice.Slip.ID).Scan(&count)
	if err != nil && err != sql.ErrNoRows {
		log.Println("Error querying database:", err)
		http.Error(w, "Oops! Something went wrong.", http.StatusInternalServerError)
		return
	}

	if count == 0 {
		// Insert advice into the database
		_, err = db.Exec("INSERT INTO advice (id, advice, count) VALUES (?, ?, ?)", advice.Slip.ID, advice.Slip.Advice, 1)
		if err != nil {
			log.Println("Error inserting advice:", err)
			http.Error(w, "Oops! Something went wrong.", http.StatusInternalServerError)
			return
		}
	} else {
		// Update the count for the existing advice
		_, err = db.Exec("UPDATE advice SET count = count + 1 WHERE id = ?", advice.Slip.ID)
		if err != nil {
			log.Println("Error updating advice count:", err)
			http.Error(w, "Oops! Something went wrong.", http.StatusInternalServerError)
			return
		}
		count++
	}

	advice.Count = count

	if advice.Count == 0 {
		// return the advice string
		fmt.Fprintf(w, "%s \n <br /> <span class='text-primary'>first time displayed</span>", advice.Slip.Advice)
	} else {
		// return the advice string
		fmt.Fprintf(w, "%s \n <br /> <span class='text-primary'>(displayed %d times)</span>", advice.Slip.Advice, advice.Count)
	}
}
