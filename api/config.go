package api

import (
	"io/ioutil"

	"github.com/cloudfoundry-community/bui/bosh"
	"github.com/cloudfoundry-community/bui/uaa"
	"github.com/gorilla/sessions"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Addr         string `yaml:"listen_addr"`
	BoshAddr     string `yaml:"bosh_addr"`
	UAA          UAA    `yaml:"uaa"`
	WebRoot      string `yaml:"web_root"`
	CookieSecret string `yaml:"cookie_secret"`
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

	var uaaClient *uaa.Client

	boshConfig := bosh.DefaultConfig()
	boshConfig.BOSHAddress = config.BoshAddr
	boshClient, err := bosh.NewClient(boshConfig)
	if err != nil {
		return err
	}
	boshInfo, err := boshClient.GetInfo()
	if err != nil {
		return err
	}
	if boshInfo.UserAuthenication.Type == "uaa" {
		uaaConfig := uaa.DefaultConfig()
		uaaConfig.Address = boshInfo.UserAuthenication.Options.URL
		uaaConfig.ClientID = config.UAA.ClientID
		uaaConfig.ClientSecret = config.UAA.ClientSecret
		uaaClient, err = uaa.NewClient(uaaConfig)
		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	ws := WebServer{
		Addr:          config.Addr,
		WebRoot:       config.WebRoot,
		Api:           a,
		CookieSession: sessions.NewCookieStore([]byte(config.CookieSecret)),
		BOSHClient:    boshClient,
		UAAClient:     uaaClient,
	}
	a.Web = &ws
	return nil
}
