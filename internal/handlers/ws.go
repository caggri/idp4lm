package handlers

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"

	"github.com/caggri/idp4lm/internal/k8s"
	"github.com/gorilla/websocket"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// TerminalPageData holds data for the terminal HTML page
type TerminalPageData struct {
	Namespace string
	PodName   string
}

// TerminalPage renders the HTML page with xterm.js
func TerminalPage(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/terminal-page/"), "/")
	if len(parts) != 2 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	data := TerminalPageData{
		Namespace: parts[0],
		PodName:   parts[1],
	}

	t, err := template.ParseFiles("web/templates/terminal.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}

// WebSocketEndpoint handles the WS connection and bridges it to Kubernetes pod exec
func WebSocketEndpoint(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/terminal/"), "/")
	if len(parts) != 2 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	namespace := parts[0]
	podName := parts[1]

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("WebSocket upgrade error: %v\n", err)
		return
	}
	defer ws.Close()

	req := k8s.Client.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Command: []string{"/bin/sh"},
			Stdin:   true,
			Stdout:  true,
			Stderr:  true,
			TTY:     true,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(k8s.Config, "POST", req.URL())
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("\r\nError connecting to pod: %v\r\n", err)))
		return
	}

	// Create a bridge between WebSocket and Executor IO
	wsReader := &wsReader{ws: ws}
	wsWriter := &wsWriter{ws: ws}

	err = exec.StreamWithContext(r.Context(), remotecommand.StreamOptions{
		Stdin:  wsReader,
		Stdout: wsWriter,
		Stderr: wsWriter,
		Tty:    true,
	})

	if err != nil {
		fmt.Printf("Stream error: %v\n", err)
	}
}

type wsReader struct {
	ws *websocket.Conn
}

func (r *wsReader) Read(p []byte) (n int, err error) {
	_, msg, err := r.ws.ReadMessage()
	if err != nil {
		if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
			return 0, io.EOF
		}
		return 0, err
	}
	n = copy(p, msg)
	return n, nil
}

type wsWriter struct {
	ws *websocket.Conn
}

func (w *wsWriter) Write(p []byte) (n int, err error) {
	err = w.ws.WriteMessage(websocket.TextMessage, p)
	if err != nil {
		if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
			return 0, io.EOF
		}
		return 0, err
	}
	return len(p), nil
}
