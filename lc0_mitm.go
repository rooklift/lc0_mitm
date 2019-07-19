package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/gorilla/websocket"      // go get github.com/gorilla/websocket
)

var Upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {return true}}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe("127.0.0.1:8080", nil)
}

func handler(writer http.ResponseWriter, request * http.Request) {

	fmt.Printf("Incoming connection...\n")

	conn, err := Upgrader.Upgrade(writer, request, nil)
	if err != nil {
		return
	}

	var lc0 * exec.Cmd

	if _, err := os.Stat("./lc0.exe"); err == nil {
		lc0 = exec.Command("./lc0.exe")
	} else if _, err := os.Stat("./lc0"); err == nil {
		lc0 = exec.Command("./lc0")
	} else {
		fmt.Printf("Could not find lc0.exe or lc0\n")
		return
	}

	i_pipe, _ := lc0.StdinPipe()
	o_pipe, _ := lc0.StdoutPipe()
	e_pipe, _ := lc0.StderrPipe()

	lc0.Start()

	i_pipe.Write([]byte("uci\n"))
	i_pipe.Write([]byte("setoption name VerboseMoveStats value true\n"))
	i_pipe.Write([]byte("setoption name LogLiveStats value true\n"))
	i_pipe.Write([]byte("setoption name MultiPV value 500\n"))
	i_pipe.Write([]byte("setoption name SmartPruningFactor value 0\n"))
	i_pipe.Write([]byte("setoption name ScoreType value centipawn\n"))
	i_pipe.Write([]byte("ucinewgame\n"))

	go incoming_ws_to_stdin(conn, i_pipe, lc0)
	go stdout_to_outgoing_ws(conn, o_pipe)
	go consume_stderr(e_pipe)
}

func incoming_ws_to_stdin(conn * websocket.Conn, stdin io.WriteCloser, lc0 * exec.Cmd) {
	for {
		_, b, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("Connection closed, killing lc0.\n")
			lc0.Process.Kill()
			return
		}

		if bytes.Contains(b, []byte("setoption")) {		// Disallow setoption from frontend
			continue
		}

		stdin.Write(b)
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
