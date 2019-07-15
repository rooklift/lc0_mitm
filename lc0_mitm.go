package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os/exec"

	"github.com/gorilla/websocket"      // go get github.com/gorilla/websocket
)

var Upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {return true}}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe("127.0.0.1:8080", nil)
}

func handler(writer http.ResponseWriter, request * http.Request) {

	conn, err := Upgrader.Upgrade(writer, request, nil)
	if err != nil {
		return
	}

	lc0 := exec.Command("./lc0.exe")

	i_pipe, _ := lc0.StdinPipe()
	o_pipe, _ := lc0.StdoutPipe()
	e_pipe, _ := lc0.StderrPipe()

	lc0.Start()

	go incoming_ws_to_stdin(conn, i_pipe, lc0)
	go stdout_to_outgoing_ws(conn, o_pipe)
	go consume_stderr(e_pipe)
}

func incoming_ws_to_stdin(conn * websocket.Conn, stdin io.WriteCloser, lc0 * exec.Cmd) {
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("WS connection closed, killing lc0\n")
			lc0.Process.Kill()
			return
		}
		stdin.Write(p)
		stdin.Write([]byte{'\n'})
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
	for scanner.Scan() {
		// fmt.Printf(scanner.Text() + "\n")
	}
}
