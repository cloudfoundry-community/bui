package api

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/lnguyen/bui/bosh"

	"golang.org/x/crypto/ssh"
)

func makeSSHKeyPair() ([]byte, []byte, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	privKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privKey),
	}

	privKeyBuf := bytes.NewBufferString("")

	err = pem.Encode(privKeyBuf, privKeyPEM)
	if err != nil {
		return nil, nil, err
	}

	pub, err := ssh.NewPublicKey(&privKey.PublicKey)
	if err != nil {
		return nil, nil, err
	}

	return privKeyBuf.Bytes(), ssh.MarshalAuthorizedKey(pub), nil
}

func createPrivateKey(key []byte) string {
	tmpfile, err := ioutil.TempFile("", "privKey")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Write(key); err != nil {
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}
	return tmpfile.Name()
}

func getAuthInfo(cookie *sessions.CookieStore, w http.ResponseWriter, req *http.Request) bosh.Auth {
	session, err := cookie.Get(req, "auth")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return bosh.Auth{}
	}
	username := session.Values["username"].(string)
	password := session.Values["password"].(string)
	token := session.Values["token"].(string)
	auth := bosh.Auth{
		Username: username,
		Password: password,
		Token:    token,
	}
	return auth
}
