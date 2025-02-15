package routing

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"app.flower.clip/src/templates"
	"github.com/gorilla/sessions"
	"github.com/michaeljs1990/sqlitestore"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) loginPageHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := authenticate(r, s.Store)
	if err != nil {
		log.Println(err)
	}
	authenticated := false
	if userID != 0 {
		authenticated = true
	}
	component := templates.LoginPage(authenticated)
	component.Render(r.Context(), w)
}

func (s *Service) loginApiHandler(w http.ResponseWriter, r *http.Request) {
	req_email := r.FormValue("email")
	req_password := r.FormValue("password")
	if req_email == "" {
		log.Println("null email")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("email is null"))
		return
	}
	if req_password == "" {
		log.Println("null password")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("password is null"))
		return
	}
	row := s.DB.QueryRow(`SELECT id, password FROM users WHERE email = ?`, req_email)
	var hashedPasword string
	var userID int
	err := row.Scan(&userID, &hashedPasword)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("email not found: user does not exist"))
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(hashedPasword), []byte(req_password)); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("incorrect email or password"))
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
	session.Values["user_id"] = userID
	if err = session.Save(r, w); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("server error"))
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
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
	createdAt := time.Now().Format(time.RFC3339)
	res, err := s.DB.Exec(`INSERT INTO users (fullname, email, password,
		created_at) VALUES (?, ?, ?, ?)`, req_fullname, req_email,
		string(passwordHash), createdAt)
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

	http.Redirect(w, r, "/", http.StatusFound)
}

func authenticate(r *http.Request, store *sqlitestore.SqliteStore) (userID int, err error) {
	session, err := store.Get(r, "auth-store")
	if err != nil {
		log.Println(err)
		return
	}
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		err = fmt.Errorf("incorrect type for user_id")
		log.Println(err)
		return
	}
	return
}

func (s *Service) profilePageHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := authenticate(r, s.Store)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("please login before continuing"))
		return
	}
	row := s.DB.QueryRow(`SELECT fullname, email FROM users WHERE id = ?`, userID)
	var userFullname, userEmail string
	if err = row.Scan(&userFullname, &userEmail); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("please login before continuing"))
		return
	}
	component := templates.ProfilePage(userFullname, userEmail)
	component.Render(r.Context(), w)
}

func (s *Service) logoutApiHandler(w http.ResponseWriter, r *http.Request) {
	defer http.Redirect(w, r, "/", http.StatusFound)

	session, err := s.Store.Get(r, "auth-store")
	if err != nil {
		log.Println(err)
		return
	}
	if s.Store.Delete(r, w, session); err != nil {
		log.Println(err)
	}
}
