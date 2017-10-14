package tinyurl

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mizoR/go-tinyurl/model"
)

type Server struct {
	Port int
}

type TinyurlForm struct {
	Url string `json:"url"`
}

func NewServer() *Server {
	s := &Server{
		Port: 8080,
	}

	return s
}

func (s Server) Start() error {
	log.Print("== Initialize DB")

	db, err := sql.Open("sqlite3", ":memory:")

	if err != nil {
		log.Fatal(err)

		return err
	}

	_, err = db.Exec(
		`CREATE TABLE "tinyurls" (
			"id"   INTEGER PRIMARY KEY AUTOINCREMENT,
			"slug" VARCHAR(255)  NOT NULL UNIQUE,
			"url"  VARCHAR(1024) NOT NULL
		)`,
	)

	if err != nil {
		return err
	}

	log.Print("== Initialized DB successfully")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")

		rootHandler(w, r, db)
	})

	http.HandleFunc("/api/tinyurls", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case "GET":
			listTinyurlsHandler(w, r, db)
		case "POST":
			createTinyurlHandler(w, r, db)
		default:
			writeResponse(w, 405, "{}")
		}
	})

	http.ListenAndServe(":"+strconv.Itoa(s.Port), nil)

	return nil
}

func rootHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.URL.Path == "/" {
		writeResponse(w, 200, "Tinyurl")

		return
	}

	slug := r.URL.Path[1:]

	tinyurl, err := model.FindTinyurlBySlug(db, slug)

	if err != nil {
		writeResponse(w, 500, "Internal server error")

		return
	}

	if tinyurl == nil {
		writeResponse(w, 404, "Not found")

		return
	}

	w.Header().Set("Location", tinyurl.Url)
	writeResponse(w, 302, "Redirecting...")
}

func listTinyurlsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	entities, err := model.FindAllTinyurls(db)

	if err != nil {
		writeResponse(w, 500, "{}")

		return
	}

	bytes, err := json.Marshal(entities)

	if err != nil {
		writeResponse(w, 500, "{}")

		return
	}

	writeResponse(w, 200, string(bytes))
}

func createTinyurlHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Header.Get("Content-Type") != "application/json" {
		writeResponse(w, 400, "{}")

		return
	}

	body := make([]byte, r.ContentLength)

	length, err := r.Body.Read(body)

	if err != nil && err != io.EOF {
		writeResponse(w, 400, "{}")

		return
	}

	var form TinyurlForm

	err = json.Unmarshal(body[:length], &form)

	if err != nil {
		writeResponse(w, 400, "{}")

		return
	}

	slug := random(5)

	tinyurl, err := model.CreateTinyurl(db, slug, form.Url)

	if err != nil {
		writeResponse(w, 500, "{}")

		return
	}

	bytes, err := json.Marshal(tinyurl)

	if err != nil {
		writeResponse(w, 500, "{}")

		return
	}

	writeResponse(w, 201, string(bytes))
}

func writeResponse(w http.ResponseWriter, code int, body string) {
	w.WriteHeader(code)
	fmt.Fprintf(w, body)
}

func random(n int) string {
	const alphabets = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, n)

	for i := range b {
		j := rand.Int63() % int64(len(alphabets))

		b[i] = alphabets[int(j)]
	}

	return string(b)
}
