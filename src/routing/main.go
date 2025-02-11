package routing

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"app.flower.clip/src/templates"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/michaeljs1990/sqlitestore"
)

func (s *Service) rootHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := authenticate(r, s.Store)
	if err != nil {
		log.Println(err)
	}
	authenticated := false
	if userID != 0 {
		authenticated = true
	}
	component := templates.IndexPage(authenticated)
	component.Render(r.Context(), w)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

type Service struct {
	DB    *sql.DB
	Store *sqlitestore.SqliteStore
}

func StartServer() error {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}
	db_name := os.Getenv("DATABASE_URL")
	if db_name == "" {
		log.Fatal("DATABASE_URL not found in .env")
	}
	mux := http.NewServeMux()
	db, err := sql.Open("sqlite3", db_name)
	if err != nil {
		log.Fatal(err)
	}
	sessionSecret := os.Getenv("SESSION_SECRET")
	sessionDb := os.Getenv("SESSION_DB")
	store, err := sqlitestore.NewSqliteStore(sessionDb, "sessions", "/", 3600, []byte(sessionSecret))
	if err != nil {
		panic(err)
	}
	s := Service{
		DB:    db,
		Store: store,
	}
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	mux.HandleFunc("GET /login", s.loginPageHandler)
	mux.HandleFunc("POST /login", s.loginApiHandler)
	mux.HandleFunc("GET /signup", signupPageHandler)
	mux.HandleFunc("POST /signup", s.signupApiHandler)
	mux.HandleFunc("GET /profile", s.profilePageHandler)
	mux.HandleFunc("GET /", s.rootHandler)
	log.Println("http://localhost:4000")
	return http.ListenAndServe(":4000", LoggingMiddleware(mux))
}
