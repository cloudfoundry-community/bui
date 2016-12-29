package api

import (
	"io/ioutil"

	"github.com/cloudfoundry-community/bui/bosh"
	"github.com/gorilla/sessions"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Addr         string `yaml:"listen_addr"`
	BoshAddr     string `yaml:"bosh_addr"`
	WebRoot      string `yaml:"web_root"`
	CookieSecret string `yaml:"cookie_secret"`
}

func (a *Api) ReadConfig(path string) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var config Config
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		return err
	}

	if config.Addr == "" {
		config.Addr = ":9304"
	}

	if config.WebRoot == "" {
		config.WebRoot = "/usr/share/bui/webui"
	}

	if config.CookieSecret == "" {
		config.CookieSecret = "something-secret"
	}

	if config.BoshAddr == "" {
		config.BoshAddr = "https://192.168.50.4:25555"
	}

	boshConfig := bosh.DefaultConfig()
	boshConfig.BOSHAddress = config.BoshAddr
	boshClient, err := bosh.NewClient(boshConfig)
	if err != nil {
		return err
	}

	ws := WebServer{
		Addr:          config.Addr,
		WebRoot:       config.WebRoot,
		Api:           a,
		CookieSession: sessions.NewCookieStore([]byte(config.CookieSecret)),
		BoshClient:    boshClient,
	}
	a.Web = &ws
	return nil
}
