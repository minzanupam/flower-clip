package routing

import (
	"net/http"
	"app.flower.clip/src/templates"
)

func loginPageHandler(w http.ResponseWriter, r *http.Request) {
	component := templates.LoginPage()
	component.Render(r.Context(), w)
}

func loginApiHandler(w http.ResponseWriter, r *http.Request) {
}
