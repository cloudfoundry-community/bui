package uaa_test

import (
	"net/http"
	"net/http/httptest"

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

func setup(mock MockRoute) {
	setupMultiple([]MockRoute{mock})
}

func setupMultiple(mockEndpoints []MockRoute) {
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

	m.Action(r.Handle)
	mux.Handle("/", m)
}

func teardown() {
	server.Close()
}
