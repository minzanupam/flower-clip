package routing

import (
	"log"
	"net/http"

	"app.flower.clip/src/templates"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
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
	req_fullname := r.FormValue("fullname")
	req_email := r.FormValue("email")
	req_password := r.FormValue("password")
	if req_fullname == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error fullname not found"))
		return
	}
	if req_email == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error fullname not found"))
		return
	}
	if req_password == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error fullname not found"))
		return
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req_password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("server error"))
		return
	}
	res, err := s.DB.Exec(`INSERT INTO users (fullname, email, password) VALUES (?, ?, ?)`, req_fullname, req_email, string(passwordHash))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("server error"))
		return
	}

	var session *sessions.Session
	session, err = s.Store.Get(r, "auth-store")
	if err != nil {
		log.Println(err)
		session, err = s.Store.New(r, "auth-store")
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("server error"))
			return
		}
	}
	userid, err := res.LastInsertId()
	session.Values["user_id"] = int(userid)
	if err = session.Save(r, w); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("server error"))
		return
	}
}
