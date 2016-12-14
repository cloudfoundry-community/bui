package api

type Api struct {
	Web *WebServer /* Webserver that gets spawned to handle http requests */

}

var Version = "(development)"

func NewApi() *Api {
	return &Api{}
}

func (a *Api) Run() error {
	a.Web.Start()
	return nil
}
