package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/coreos/go-log/log"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/lnguyen/bui/bosh"
)

type BOSHHandler struct {
	CookieSession *sessions.CookieStore
	BoshClient    *bosh.Client
}

type ErrorResponse struct {
	Error       string `json:"error,omitempty"`
	Description string `json:"description"`
}

func (b BOSHHandler) sessions(w http.ResponseWriter, req *http.Request) {
	session, _ := b.CookieSession.Get(req, "auth")
	fmt.Println(session)
	http.Redirect(w, req, "/#/dashboard", http.StatusFound)
}

func (b BOSHHandler) currentUser(w http.ResponseWriter, req *http.Request) {
	session, err := b.CookieSession.Get(req, "auth")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(session.Values["username"])
	if session.Values["username"] == nil {
		b.respond(w, http.StatusForbidden, map[string]string{
			"error": "currently not logged in",
		})
		return
	}
	b.respond(w, http.StatusOK, map[string]string{
		"username": session.Values["username"].(string),
	})
}

func (b BOSHHandler) login(w http.ResponseWriter, req *http.Request) {
	session, err := b.CookieSession.Get(req, "auth")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.ParseForm()
	username := req.PostFormValue("username")
	password := req.PostFormValue("password")
	r := b.BoshClient.NewRequest("GET", "/releases")
	resp, err := b.BoshClient.DoAuthRequest(r, username, password, "")
	if err != nil {
		b.respond(w, http.StatusInternalServerError, ErrorResponse{
			Description: err.Error(),
		})
		return
	}
	if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("unauthorized")
		b.respond(w, http.StatusUnauthorized, ErrorResponse{
			Description: "Unauthorized",
		})
		return
	}
	session.Values["auth_type"] = req.PostFormValue("auth_type")
	session.Values["username"] = username
	session.Values["password"] = password
	session.Values["token"] = req.PostFormValue("token")
	session.Save(req, w)
	http.Redirect(w, req, "/#/dashboard", http.StatusFound)
}

func (b BOSHHandler) info(w http.ResponseWriter, req *http.Request) {
	r := b.BoshClient.NewRequest("GET", "/info")
	resp, err := b.BoshClient.DoRequest(r)
	if err != nil {
		b.respond(w, http.StatusInternalServerError, ErrorResponse{
			Description: err.Error(),
		})
		return
	}
	var data []byte
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		b.respond(w, http.StatusInternalServerError, ErrorResponse{
			Description: err.Error(),
		})
		return
	}
	w.WriteHeader(resp.StatusCode)
	w.Write(data)
}

func (b BOSHHandler) releases(w http.ResponseWriter, req *http.Request) {
	session, err := b.CookieSession.Get(req, "auth")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	username := session.Values["username"].(string)
	password := session.Values["password"].(string)
	token := session.Values["token"].(string)
	r := b.BoshClient.NewRequest("GET", "/releases")
	resp, err := b.BoshClient.DoAuthRequest(r, username, password, token)
	if err != nil {
		b.respond(w, http.StatusInternalServerError, ErrorResponse{
			Description: err.Error(),
		})
		return
	}
	var data []byte
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		b.respond(w, http.StatusInternalServerError, ErrorResponse{
			Description: err.Error(),
		})
		return
	}
	w.WriteHeader(resp.StatusCode)
	w.Write(data)
}

func (b BOSHHandler) stemcells(w http.ResponseWriter, req *http.Request) {
	session, err := b.CookieSession.Get(req, "auth")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	username := session.Values["username"].(string)
	password := session.Values["password"].(string)
	token := session.Values["token"].(string)
	r := b.BoshClient.NewRequest("GET", "/stemcells")
	resp, err := b.BoshClient.DoAuthRequest(r, username, password, token)
	if err != nil {
		b.respond(w, http.StatusInternalServerError, ErrorResponse{
			Description: err.Error(),
		})
		return
	}
	var data []byte
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		b.respond(w, http.StatusInternalServerError, ErrorResponse{
			Description: err.Error(),
		})
		return
	}
	w.WriteHeader(resp.StatusCode)
	w.Write(data)
}

func (b BOSHHandler) deployment(w http.ResponseWriter, req *http.Request) {
	session, err := b.CookieSession.Get(req, "auth")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	username := session.Values["username"].(string)
	password := session.Values["password"].(string)
	token := session.Values["token"].(string)
	vars := mux.Vars(req)
	name := vars["name"]

	r := b.BoshClient.NewRequest("GET", "/deployments/"+name)
	resp, err := b.BoshClient.DoAuthRequest(r, username, password, token)
	if err != nil {
		b.respond(w, http.StatusInternalServerError, ErrorResponse{
			Description: err.Error(),
		})
		return
	}
	var data []byte
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		b.respond(w, http.StatusInternalServerError, ErrorResponse{
			Description: err.Error(),
		})
		return
	}
	w.WriteHeader(resp.StatusCode)
	w.Write(data)
}

func (b BOSHHandler) deployments(w http.ResponseWriter, req *http.Request) {
	session, err := b.CookieSession.Get(req, "auth")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	username := session.Values["username"].(string)
	password := session.Values["password"].(string)
	token := session.Values["token"].(string)
	r := b.BoshClient.NewRequest("GET", "/deployments")
	resp, err := b.BoshClient.DoAuthRequest(r, username, password, token)
	if err != nil {
		b.respond(w, http.StatusInternalServerError, ErrorResponse{
			Description: err.Error(),
		})
		return
	}
	var data []byte
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		b.respond(w, http.StatusInternalServerError, ErrorResponse{
			Description: err.Error(),
		})
		return
	}
	w.WriteHeader(resp.StatusCode)
	w.Write(data)
}

func (b BOSHHandler) running_tasks(w http.ResponseWriter, req *http.Request) {
	session, err := b.CookieSession.Get(req, "auth")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	username := session.Values["username"].(string)
	password := session.Values["password"].(string)
	token := session.Values["token"].(string)

	r := b.BoshClient.NewRequest("GET", "/tasks?state=queued,processing,cancelling&verbose=2")
	resp, err := b.BoshClient.DoAuthRequest(r, username, password, token)
	if err != nil {
		b.respond(w, http.StatusInternalServerError, ErrorResponse{
			Description: err.Error(),
		})
		return
	}
	var data []byte
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		b.respond(w, http.StatusInternalServerError, ErrorResponse{
			Description: err.Error(),
		})
		return
	}
	w.WriteHeader(resp.StatusCode)
	w.Write(data)
}

func (b BOSHHandler) respond(w http.ResponseWriter, status int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	encoder := json.NewEncoder(w)
	err := encoder.Encode(response)
	if err != nil {
		log.Errorf("unable to encode response %s", "")
	}
}
