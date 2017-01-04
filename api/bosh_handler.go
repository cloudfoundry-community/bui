package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cloudfoundry-community/bui/bosh"
	"github.com/cloudfoundry-community/bui/uaa"
	"github.com/coreos/go-log/log"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type BOSHHandler struct {
	CookieSession *sessions.CookieStore
	BOSHClient    *bosh.Client
	UAAClient     *uaa.Client
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
	if session.Values["username"] == nil {
		b.respond(w, http.StatusForbidden, map[string]string{
			"error": "currently not logged in",
		})
		return
	}
	b.respond(w, http.StatusOK, map[string]string{
		"name": session.Values["username"].(string),
	})
}

func (b BOSHHandler) login(w http.ResponseWriter, req *http.Request) {
	var token string
	session, err := b.CookieSession.Get(req, "auth")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.ParseForm()
	username := req.PostFormValue("username")
	password := req.PostFormValue("password")

	info, err := b.BOSHClient.GetInfo()
	if err != nil {
		b.respond(w, http.StatusInternalServerError, ErrorResponse{
			Description: err.Error(),
		})
		return
	}
	if info.UserAuthenication.Type == "uaa" {
		tokenResp, err := b.UAAClient.GetPasswordToken(username, password)
		if err != nil {
			b.respond(w, http.StatusInternalServerError, ErrorResponse{
				Description: err.Error(),
			})
			return
		}
		token = tokenResp.AccessToken
	}

	auth := bosh.Auth{
		Username: username,
		Password: password,
		Token:    token,
	}
	r := b.BOSHClient.NewRequest("GET", "/releases")
	resp, err := b.BOSHClient.DoAuthRequestRaw(r, auth)
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
	session.Values["auth_type"] = info.UserAuthenication.Type
	session.Values["username"] = username
	session.Values["password"] = password
	session.Values["token"] = token
	session.Save(req, w)
	http.Redirect(w, req, "/#/dashboard", http.StatusFound)
}

func (b BOSHHandler) info(w http.ResponseWriter, req *http.Request) {
	info, err := b.BOSHClient.GetInfo()
	if err != nil {
		b.respond(w, http.StatusInternalServerError, ErrorResponse{
			Description: err.Error(),
		})
		return
	}
	b.respond(w, http.StatusOK, info)
}

func (b BOSHHandler) releases(w http.ResponseWriter, req *http.Request) {
	auth := getAuthInfo(b.CookieSession, w, req)
	releases, err := b.BOSHClient.GetReleases(auth)
	if err != nil {
		b.respond(w, http.StatusInternalServerError, ErrorResponse{
			Description: err.Error(),
		})
		return
	}
	b.respond(w, http.StatusOK, releases)
}

func (b BOSHHandler) stemcells(w http.ResponseWriter, req *http.Request) {
	auth := getAuthInfo(b.CookieSession, w, req)
	stemcells, err := b.BOSHClient.GetStemcells(auth)
	if err != nil {
		b.respond(w, http.StatusInternalServerError, ErrorResponse{
			Description: err.Error(),
		})
		return
	}
	b.respond(w, http.StatusOK, stemcells)
}

func (b BOSHHandler) deployment(w http.ResponseWriter, req *http.Request) {
	auth := getAuthInfo(b.CookieSession, w, req)
	vars := mux.Vars(req)
	name := vars["name"]

	deployment, err := b.BOSHClient.GetDeployment(name, auth)
	if err != nil {
		b.respond(w, http.StatusInternalServerError, ErrorResponse{
			Description: err.Error(),
		})
		return
	}
	b.respond(w, http.StatusOK, deployment)
}

func (b BOSHHandler) deployments(w http.ResponseWriter, req *http.Request) {
	auth := getAuthInfo(b.CookieSession, w, req)
	deployments, err := b.BOSHClient.GetDeployments(auth)
	if err != nil {
		b.respond(w, http.StatusInternalServerError, ErrorResponse{
			Description: err.Error(),
		})
		return
	}
	b.respond(w, http.StatusOK, deployments)
}

func (b BOSHHandler) deploymentVMs(w http.ResponseWriter, req *http.Request) {
	auth := getAuthInfo(b.CookieSession, w, req)
	vars := mux.Vars(req)
	name := vars["name"]
	deploymentVMs, err := b.BOSHClient.GetDeploymentVMs(name, auth)
	if err != nil {
		b.respond(w, http.StatusInternalServerError, ErrorResponse{
			Description: err.Error(),
		})
		return
	}
	b.respond(w, http.StatusOK, deploymentVMs)
}

func (b BOSHHandler) running_tasks(w http.ResponseWriter, req *http.Request) {
	auth := getAuthInfo(b.CookieSession, w, req)

	tasks, err := b.BOSHClient.GetRunningTasks(auth)
	if err != nil {
		b.respond(w, http.StatusInternalServerError, ErrorResponse{
			Description: err.Error(),
		})
		return
	}
	b.respond(w, http.StatusOK, tasks)
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
