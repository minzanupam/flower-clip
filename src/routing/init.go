package routing

import (
	"log"
	"net/http"

	"app.flower.clip/src/templates"
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

func StartServer() error {
	mux := http.NewServeMux()
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	mux.HandleFunc("GET /login", loginPageHandler)
	mux.HandleFunc("POST /login", loginApiHandler)
	mux.HandleFunc("GET /", rootHandler)
	log.Println("http://localhost:4000")
	return http.ListenAndServe(":4000", LoggingMiddleware(mux))
}
