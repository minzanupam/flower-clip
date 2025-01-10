package routing

import (
	"app.flower.clip/src/templates"
	"net/http"
)

func loginPageHandler(w http.ResponseWriter, r *http.Request) {
	component := templates.LoginPage()
	component.Render(r.Context(), w)
}

func (s *Service) loginApiHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	row := s.DB.QueryRow(`SELECT id, password FROM users WHERE email = ?`, email)
	var dbPassword string
	err := row.Scan(&dbPassword)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("email not found: user does not exist"))
		return
	}
	if password != dbPassword {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("incorrect username or password"))
		return
	}

	// create and authenticate
}

func signupPageHandler(w http.ResponseWriter, r *http.Request) {
	component := templates.SignupPage()
	component.Render(r.Context(), w)
}

func (s *Service) signupApiHandler(w http.ResponseWriter, r *http.Request) {
}
