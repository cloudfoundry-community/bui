package bosh

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
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

// request is used to help build up a request
type request struct {
	method string
	url    string
	header map[string]string
	params url.Values
	body   io.Reader
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
	client := &Client{
		config: *config,
	}

	return client, nil
}

// NewRequest is used to create a new request
func (c *Client) NewRequest(method, path string) *request {
	r := &request{
		method: method,
		url:    c.config.BOSHAddress + path,
		params: make(map[string][]string),
		header: make(map[string]string),
	}
	return r
}

// DoAuthRequest runs a request with our client
func (c *Client) DoAuthRequest(r *request, username, password, token string) (*http.Response, error) {
	req, err := r.toHTTP()
	if err != nil {
		return nil, err
	}
	for key, value := range r.header {
		req.Header.Add(key, value)
	}
	if token != "" {
		req.Header.Add("Authorization", "Bearer "+token)
	} else {
		req.SetBasicAuth(username, password)
	}
	req.Header.Add("User-Agent", "bui")
	resp, err := c.config.HttpClient.Do(req)
	return resp, err
}

func (c *Client) DoRequest(r *request) (*http.Response, error) {
	req, err := r.toHTTP()
	if err != nil {
		return nil, err
	}
	for key, value := range r.header {
		req.Header.Add(key, value)
	}
	req.Header.Add("User-Agent", "bui")
	resp, err := c.config.HttpClient.Do(req)
	return resp, err
}

// toHTTP converts the request to an HTTP request
func (r *request) toHTTP() (*http.Request, error) {

	// Check if we should encode the body
	if r.body == nil && r.obj != nil {
		if b, err := encodeBody(r.obj); err != nil {
			return nil, err
		} else {
			r.body = b
		}
	}

	// Create the HTTP request
	return http.NewRequest(r.method, r.url, r.body)
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
