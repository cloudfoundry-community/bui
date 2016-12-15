package bosh

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

//Client used to communicate with BOSH
type Client struct {
	config Config
}

//Config is used to configure the creation of a client
type Config struct {
	BOSHAddress       string
	HttpClient        *http.Client
	SkipSslValidation bool
}

type Auth struct {
	Username string
	Password string
	Token    string
}

// request is used to help build up a request
type request struct {
	method string
	url    string
	Header map[string]string
	params url.Values
	Body   io.Reader
	obj    interface{}
}

//DefaultConfig configuration for client
func DefaultConfig() *Config {
	return &Config{
		BOSHAddress:       "https://192.168.50.4:25555",
		HttpClient:        http.DefaultClient,
		SkipSslValidation: true,
	}
}

// NewClient returns a new client
func NewClient(config *Config) (*Client, error) {
	// bootstrap the config
	defConfig := DefaultConfig()

	if len(config.BOSHAddress) == 0 {
		config.BOSHAddress = defConfig.BOSHAddress
	}

	config.HttpClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: config.SkipSslValidation,
			},
		},
	}
	config.HttpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) > 10 {
			return fmt.Errorf("stopped after 10 redirects")
		}
		req.URL.Host = strings.TrimPrefix(config.BOSHAddress, req.URL.Scheme+"://")
		req.Header.Add("User-Agent", "bui")
		req.Header.Add("Authorization", via[0].Header.Get("Authorization"))
		req.Header.Del("Referer")
		return nil
	}
	client := &Client{
		config: *config,
	}

	return client, nil
}

// GetInfo returns BOSH Info
func (c *Client) GetInfo() (info Info, err error) {
	r := c.NewRequest("GET", "/info")
	resp, err := c.DoRequest(r)

	if err != nil {
		log.Printf("Error requesting info %v", err)
		return
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading info request %v", resBody)
		return
	}
	err = json.Unmarshal(resBody, &info)
	if err != nil {
		log.Printf("Error unmarshalling info %v", err)
		return
	}
	return
}

// NewRequest is used to create a new request
func (c *Client) NewRequest(method, path string) *request {
	r := &request{
		method: method,
		url:    c.config.BOSHAddress + path,
		params: make(map[string][]string),
		Header: make(map[string]string),
	}
	return r
}

// DoAuthRequest runs a request with our client
func (c *Client) DoAuthRequest(r *request, auth Auth) ([]byte, error) {
	req, err := r.toHTTP()
	if err != nil {
		return nil, err
	}
	for key, value := range r.Header {
		req.Header.Add(key, value)
	}
	if auth.Token != "" {
		req.Header.Add("Authorization", "Bearer "+auth.Token)
	} else {
		req.SetBasicAuth(auth.Username, auth.Password)
	}
	req.Header.Add("User-Agent", "bui")
	resp, err := c.config.HttpClient.Do(req)
	if err != nil {
		return []byte{}, err
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	return resBody, err
}

func (c *Client) DoRequest(r *request) (*http.Response, error) {
	req, err := r.toHTTP()
	if err != nil {
		return nil, err
	}
	for key, value := range r.Header {
		req.Header.Add(key, value)
	}
	req.Header.Add("User-Agent", "bui")
	resp, err := c.config.HttpClient.Do(req)
	return resp, err
}

// toHTTP converts the request to an HTTP request
func (r *request) toHTTP() (*http.Request, error) {

	// Check if we should encode the body
	if r.Body == nil && r.obj != nil {
		if b, err := encodeBody(r.obj); err != nil {
			return nil, err
		} else {
			r.Body = b
		}
	}

	// Create the HTTP request
	return http.NewRequest(r.method, r.url, r.Body)
}

// decodeBody is used to JSON decode a body
func decodeBody(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	return dec.Decode(out)
}

// encodeBody is used to encode a request body
func encodeBody(obj interface{}) (io.Reader, error) {
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	if err := enc.Encode(obj); err != nil {
		return nil, err
	}
	return buf, nil
}
