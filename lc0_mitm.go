package main

import (
	"bufio"
	"io"
	"net/http"
	"os/exec"

	"github.com/gorilla/websocket"      // go get github.com/gorilla/websocket
)

var Upgrader = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024, CheckOrigin: func(r *http.Request) bool {return true}}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe("127.0.0.1:8080", nil)
}

func handler(writer http.ResponseWriter, request * http.Request) {

	conn, err := Upgrader.Upgrade(writer, request, nil)
	if err != nil {
		return
	}

	exec_command := exec.Command("./lc0.exe")

	i_pipe, _ := exec_command.StdinPipe()
	o_pipe, _ := exec_command.StdoutPipe()
	e_pipe, _ := exec_command.StderrPipe()

	go incoming_ws_to_stdin(conn, i_pipe)
	go stdout_to_outgoing_ws(conn, o_pipe)
	go consume_stderr(e_pipe)
}

func incoming_ws_to_stdin(conn * websocket.Conn, stdin io.WriteCloser) {
	for {
		_, p, _ := conn.ReadMessage()
		stdin.Write(p)
	}
}

func stdout_to_outgoing_ws(conn * websocket.Conn, stdout io.ReadCloser) {
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		conn.WriteMessage(websocket.TextMessage, scanner.Bytes())
	}
}

func consume_stderr(stderr io.ReadCloser) {
	// Note that we're not allowed concurrent writes to the conn.
	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {}
}
