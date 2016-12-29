package api

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/cloudfoundry-community/bui/bosh"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/kr/pty"
	uuid "github.com/satori/go.uuid"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1,
	WriteBufferSize: 1,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type wsPty struct {
	Cmd  *exec.Cmd // pty builds on os.exec
	Pty  *os.File  // a pty is simply an os.File
	Args []string
}

type SSHRequest struct {
	Command        string `json:"command"`
	DeploymentName string `json:"deployment_name"`
	Target         Target `json:"target"`
	Params         Params `json:"params"`
}

type Target struct {
	Job     string   `json:"job"`
	Indexes []string `json:"indexes"`
	Ids     []string `json:"ids"`
}

type Params struct {
	User      string `json:"user"`
	Password  string `json:"password"`
	PublicKey string `json:"public_key"`
}

func (wp *wsPty) Start() {
	var err error
	wp.Cmd = exec.Command("ssh", wp.Args...)
	wp.Pty, err = pty.Start(wp.Cmd)
	if err != nil {
		log.Fatalf("Failed to start command: %s\n", err)
	}
}

func (wp *wsPty) Stop() {
	wp.Pty.Close()
	wp.Cmd.Wait()
}

func (b BOSHHandler) ssh(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	name := vars["name"]
	vm := strings.Split(vars["vm"], "-")
	auth := getAuthInfo(b.CookieSession, w, req)
	privKey, pubKey, err := makeSSHKeyPair()
	if err != nil {
		log.Printf("Websocket upgrade failed: %s\n", err)
		return
	}
	user := uuid.NewV4()
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	username := "bosh_" + strings.Replace(user.String(), "-", "", -1)[0:16]
	privKeyPath := createPrivateKey(privKey)
	defer os.Remove(privKeyPath)

	sshRequest := bosh.SSHRequest{
		Command:        "setup",
		DeploymentName: name,
		Target: bosh.Target{
			Job: vm[0],
			Ids: []string{vm[1]},
		},
		Params: map[string]string{
			"user":       username,
			"password":   "",
			"public_key": string(pubKey),
		},
	}
	response, err := b.BoshClient.SSH(sshRequest, auth)
	if err != nil {
		log.Printf("Websocket upgrade failed: %s\n", err)
	}
	fmt.Println(response)

	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Printf("Websocket upgrade failed: %s\n", err)
	}
	defer conn.Close()

	wp := wsPty{
		Args: []string{"-oStrictHostKeyChecking=no", "-i", privKeyPath, username + "@" + response[0].IP},
	}
	// TODO: check for errors, return 500 on fail
	wp.Start()

	// copy everything from the pty master to the websocket
	// using base64 encoding for now due to limitations in term.js
	go func() {
		buf := make([]byte, 128)
		// TODO: more graceful exit on socket close / process exit
		for {
			n, err := wp.Pty.Read(buf)
			if err != nil {
				log.Printf("Failed to read from pty master: %s", err)
				return
			}

			out := make([]byte, base64.StdEncoding.EncodedLen(n))
			base64.StdEncoding.Encode(out, buf[0:n])

			err = conn.WriteMessage(websocket.TextMessage, out)

			if err != nil {
				log.Printf("Failed to send %d bytes on websocket: %s", n, err)
				return
			}
		}
	}()

	// read from the web socket, copying to the pty master
	// messages are expected to be text and base64 encoded
	for {
		mt, payload, err := conn.ReadMessage()
		if err != nil {
			if err != io.EOF {
				log.Printf("conn.ReadMessage failed: %s\n", err)
				return
			}
		}

		switch mt {
		case websocket.BinaryMessage:
			log.Printf("Ignoring binary message: %q\n", payload)
		case websocket.TextMessage:
			buf := make([]byte, base64.StdEncoding.DecodedLen(len(payload)))
			_, err := base64.StdEncoding.Decode(buf, payload)
			if err != nil {
				log.Printf("base64 decoding of payload failed: %s\n", err)
			}
			wp.Pty.Write(buf)
		default:
			log.Printf("Invalid message type %d\n", mt)
			return
		}
	}
	// TODO: Clean up SSH
	/*sshCleanup := bosh.SSHRequest{
		Command:        "cleanup",
		DeploymentName: name,
		Target: bosh.Target{
			Job: vm[0],
			Ids: []string{vm[1]},
		},
		Params: map[string]string{
			"user_regex": "^" + username,
		},
	}
	response, err = b.BoshClient.SSH(sshCleanup, auth)
	if err != nil {
		log.Printf("Websocket upgrade failed: %s\n", err)
		return
	}
	wp.Stop()*/
}
