package routing

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"app.flower.clip/src/templates"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	component := templates.IndexPage()
	component.Render(r.Context(), w)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

type Service struct {
	DB *sql.DB
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
	s := Service{
		DB: db,
	}
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	mux.HandleFunc("GET /login", loginPageHandler)
	mux.HandleFunc("POST /login", s.loginApiHandler)
	mux.HandleFunc("GET /signup", signupPageHandler)
	mux.HandleFunc("POST /signup", s.signupApiHandler)
	mux.HandleFunc("GET /", rootHandler)
	log.Println("http://localhost:4000")
	return http.ListenAndServe(":4000", LoggingMiddleware(mux))
}
