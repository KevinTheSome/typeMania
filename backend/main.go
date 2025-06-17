package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "github.com/mattn/go-sqlite3"
)

const ServerPort = ":8080"
var db *sql.DB

func getScores() []Score {
	var scores []Score
	rows, err := db.Query("select * from scores")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var score Score
		if err := rows.Scan(&score.Id, &score.Name, &score.Wpm, &score.Errors, &score.Characters, &score.Accuracy, &score.WordCount, &score.Time); err != nil {
			log.Fatal(err)
		}
		scores = append(scores, score)
	}
	return scores
}

func addScores(score Score) string {
	log.Print(score)
	sqlStatement := `
		INSERT INTO scores (name, wpm, errors, characters, accuracy, wordCount, time)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	res,err := db.Exec(sqlStatement, score.Name, score.Wpm, score.Errors, score.Characters, score.Accuracy, score.WordCount, score.Time)
	if err != nil {
		log.Fatal(err)
		return err.Error()
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
		return err.Error()
	}
	return string(id)

}

func getLb(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	send, err := json.Marshal(getScores())
	if err != nil {
		log.Fatal("json.Marshal Fucking blow up!")
	}
	if send == nil {
		send = []byte("null")
	}

	w.Write([]byte(send))

}

func addScore(w http.ResponseWriter, r *http.Request) {
	var tempScore Score
	err := json.NewDecoder(r.Body).Decode(&tempScore)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	send := addScores(tempScore)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(send))
}

func main() {

	dbc, err := sql.Open("sqlite3", "./db.db")
	db = dbc
	if err != nil {
		log.Fatal(err)
	}

	_, err = os.Stat("./db.db")
	if os.IsNotExist(err) {
		log.Print("db.db does not exist. Creating it.")
		_, err = db.Exec(`
		create table scores (
		id integer not null primary key,
		name text,
		wpm integer,
		errors integer,
		characters integer,
		accuracy integer,
		wordCount integer,
		time integer
		)
		`)
		if err != nil {
			log.Fatal(err)
		}
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
  	}))
	r.Get("/", getLb)
	r.Post("/", addScore)

	log.Print("Server running on port", ServerPort)
	http.ListenAndServe(ServerPort, r)
	defer db.Close()
}