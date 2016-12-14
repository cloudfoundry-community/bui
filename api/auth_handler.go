package api

import (
	"net/http"

	"github.com/gorilla/sessions"
)

type authHandler struct {
	CookieSession *sessions.CookieStore
	handler       http.Handler
}

func AuthHandler(cookieSession *sessions.CookieStore, h http.Handler) http.Handler {
	return authHandler{cookieSession, h}
}

func (h authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session, err := h.CookieSession.Get(r, "auth")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if session.Values["auth_type"] == nil {
		http.Redirect(w, r, "/#!/login", http.StatusFound)
		return
	}
	h.handler.ServeHTTP(w, r)
}
