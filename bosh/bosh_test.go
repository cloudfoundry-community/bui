package bosh_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/lnguyen/bui/bosh"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

var (
	mux           *http.ServeMux
	server        *httptest.Server
	fakeUAAServer *httptest.Server
)

type MockRoute struct {
	Method   string
	Endpoint string
	Output   string
	Redirect string
}

func setup(mock MockRoute, authType string) {
	setupMultiple([]MockRoute{mock}, authType)
}

func setupMultiple(mockEndpoints []MockRoute, authType string) {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	m := martini.New()
	m.Use(render.Renderer())
	r := martini.NewRouter()
	for _, mock := range mockEndpoints {
		method := mock.Method
		endpoint := mock.Endpoint
		output := mock.Output
		redirect := mock.Redirect
		if redirect != "" {
			r.Get(endpoint, func(r render.Render) {
				r.Redirect(redirect)
			})
		}
		if method == "GET" {
			r.Get(endpoint, func() string {
				return output
			})
		} else if method == "POST" {
			r.Post(endpoint, func() string {
				return output
			})
		} else if method == "DELETE" {
			r.Delete(endpoint, func() (int, string) {
				return 204, output
			})
		}
	}
	if authType != "uaa" {
		r.Get("/info", func(r render.Render) {
			r.JSON(200, &Info{
				Name:    "bosh-lite",
				UUID:    "2daf673a-9755-4b4f-aa6d-3632fbed8019",
				Version: "1.3126.0 (00000000)",
				User:    "admin",
				CPI:     "warden_cpi",
				UserAuthenication: UserAuthenication{
					Type: "basic",
				},
			})
		})
	} else {
		r.Get("/info", func(r render.Render) {
			r.JSON(200, &Info{
				Name:    "bosh-lite",
				UUID:    "2daf673a-9755-4b4f-aa6d-3632fbed8012",
				Version: "1.3126.0 (00000000)",
				User:    "admin",
				CPI:     "warden_cpi",
				UserAuthenication: UserAuthenication{
					Type: "uaa",
					Options: struct {
						URL string `json:"url"`
					}{fakeUAAServer.URL},
				},
			})
		})
	}

	m.Action(r.Handle)
	mux.Handle("/", m)
}

func teardown() {
	server.Close()
}
