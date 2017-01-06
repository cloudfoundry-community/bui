package api

import (
	"net/http"

	"github.com/cloudfoundry-community/bui/bosh"
	"github.com/cloudfoundry-community/bui/uaa"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/starkandwayne/goutils/log"
)

type WebServer struct {
	Addr          string
	WebRoot       string
	Api           *Api
	CookieSession *sessions.CookieStore
	BOSHClient    *bosh.Client
	UAAClient     *uaa.Client
}

// Setup webserver
func (ws *WebServer) Setup() error {
	log.Debugf("Configuring WebServer...")

	boshHandler := BOSHHandler{
		CookieSession: ws.CookieSession,
		BOSHClient:    ws.BOSHClient,
		UAAClient:     ws.UAAClient,
	}

	r := mux.NewRouter()
	r.Handle("/info2", AuthHandler(ws.CookieSession, http.HandlerFunc(boshHandler.info))).Methods("GET")
	r.HandleFunc("/user", boshHandler.currentUser).Methods("GET")
	r.HandleFunc("/login", boshHandler.login).Methods("POST")
	r.HandleFunc("/info", boshHandler.info).Methods("GET")
	r.Handle("/releases", AuthHandler(ws.CookieSession, http.HandlerFunc(boshHandler.releases))).Methods("GET")
	r.Handle("/stemcells", AuthHandler(ws.CookieSession, http.HandlerFunc(boshHandler.stemcells))).Methods("GET")
	r.Handle("/deployments", AuthHandler(ws.CookieSession, http.HandlerFunc(boshHandler.deployments))).Methods("GET")
	r.Handle("/deployments/{name}", AuthHandler(ws.CookieSession, http.HandlerFunc(boshHandler.deployment))).Methods("GET")
	r.Handle("/deployments/{name}/vms", AuthHandler(ws.CookieSession, http.HandlerFunc(boshHandler.deployment_vms))).Methods("GET")
	r.Handle("/deployments/{name}/vms/{vm}/ssh", AuthHandler(ws.CookieSession, http.HandlerFunc(boshHandler.ssh)))
	r.Handle("/tasks/running", AuthHandler(ws.CookieSession, http.HandlerFunc(boshHandler.running_tasks))).Methods("GET")
	r.HandleFunc("/sessions", boshHandler.sessions).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(ws.WebRoot)))
	http.Handle("/", r)
	return nil
}

// Start webserver
func (ws *WebServer) Start() {
	err := ws.Setup()
	if err != nil {
		panic("Could not set up WebServer for B-ui: " + err.Error())
	}
	log.Debugf("Starting WebServer on '%s'...", ws.Addr)
	err = http.ListenAndServe(ws.Addr, nil)
	if err != nil {
		log.Errorf("HTTP API failed %s", err.Error())
		panic("Cannot setup WebServer, aborting. Check logs for details.")
	}
}
