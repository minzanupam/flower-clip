package routing

import (
	"log"
	"net/http"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
	return
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func StartServer() error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", rootHandler)
	log.Println("http://localhost:4000")
	return http.ListenAndServe(":4000", LoggingMiddleware(mux))
}
