package uaa

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

//Client used to communicate with BOSH
type Client struct {
	config Config
}

//Config is used to configure the creation of a client
type Config struct {
	ClientID          string
	ClientSecret      string
	Address           string
	HttpClient        *http.Client
	SkipSslValidation bool
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
		Address:           "https://192.168.50.4:8443",
		HttpClient:        http.DefaultClient,
		SkipSslValidation: true,
	}
}

// NewClient returns a new client
func NewClient(config *Config) (*Client, error) {
	// bootstrap the config
	defConfig := DefaultConfig()

	if len(config.Address) == 0 {
		config.Address = defConfig.Address
	}

	if len(config.ClientID) == 0 {
		config.ClientID = defConfig.ClientID
	}

	if len(config.ClientSecret) == 0 {
		config.ClientSecret = defConfig.ClientSecret
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

// GetInfo returns UAA Info
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

func (c *Client) GetPasswordToken(username, password string) (tokenResp TokenResp, err error) {
	data := url.Values{
		"grant_type": {"password"},
		"username":   {username},
		"password":   {password},
	}
	r := c.NewRequest("POST", fmt.Sprintf("/oauth/token?%s", data.Encode()))

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
	err = json.Unmarshal(resBody, &tokenResp)
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
		url:    c.config.Address + path,
		params: make(map[string][]string),
		Header: make(map[string]string),
	}
	return r
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
	req.Header.Add("Authorization", buildAuthorizationHeader(c.config.ClientID, c.config.ClientSecret))
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

func buildAuthorizationHeader(client, clientSecret string) string {
	data := []byte(fmt.Sprintf("%s:%s", client, clientSecret))
	encodedBasicAuth := base64.StdEncoding.EncodeToString(data)
	headerString := fmt.Sprintf("Basic %s", encodedBasicAuth)

	return headerString
}
