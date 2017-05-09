package api

import (
	"io/ioutil"

	"github.com/cloudfoundry-community/bui/bosh"
	"github.com/gorilla/sessions"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Addr              string `yaml:"listen_addr"`
	BoshAddr          string `yaml:"bosh_addr"`
	UAA               UAA    `yaml:"uaa"`
	WebRoot           string `yaml:"web_root"`
	CookieSecret      string `yaml:"cookie_secret"`
	SkipSSLValidation bool   `yaml:"skip_ssl_validation"`
}

type UAA struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
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
	boshConfig.SkipSslValidation = config.SkipSSLValidation
	boshConfig.UAA.ClientID = config.UAA.ClientID
	boshConfig.UAA.ClientSecret = config.UAA.ClientSecret
	boshClient, err := bosh.NewClient(boshConfig)
	if err != nil {
		return err
	}

	ws := WebServer{
		Addr:          config.Addr,
		WebRoot:       config.WebRoot,
		Api:           a,
		CookieSession: sessions.NewCookieStore([]byte(config.CookieSecret)),
		BOSHClient:    boshClient,
	}
	a.Web = &ws
	return nil
}
