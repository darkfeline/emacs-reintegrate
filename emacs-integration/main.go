package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os/exec"
	"strings"

	"github.com/coreos/go-systemd/activation"
	"github.com/pkg/browser"
)

var (
	port = flag.String("address", "127.0.0.1:9999", "Proxy listen address.")
)

func main() {
	log.SetPrefix("emacs-integration: ")
	var listener net.Listener
	ls, err := activation.Listeners()
	if err != nil {
		log.Fatal(err)
	}
	if len(ls) >= 1 {
		listener = ls[0]
	} else {
		listener, err = net.Listen("tcp", *port)
		if err != nil {
			log.Fatal(err)
		}
	}

	http.HandleFunc("/browser", handleBrowser)
	http.HandleFunc("/clipboard", handleClipboard)
	http.HandleFunc("/health", handleHealth)
	log.Fatal(http.Serve(listener, nil))
}

func handleBrowser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		log.Printf("POST /browser")
		url, err := requestBody(r)
		if err != nil {
			serverError(w, "Error reading body: %s", err)
			return
		}
		if err := browser.OpenURL(url); err != nil {
			serverError(w, "Error opening URL: %s", err)
			return
		}
	default:
		badMethod(w, "POST")
	}
}

func handleClipboard(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		log.Printf("GET /clipboard")
		s, err := getClipboard()
		if err != nil {
			serverError(w, "Error getting clipboard: %s", err)
			return
		}
		io.WriteString(w, s)
	case "PUT":
		log.Printf("PUT /clipboard")
		s, err := requestBody(r)
		if err != nil {
			serverError(w, "Error reading body: %s", err)
			return
		}
		if err := setClipboard(s); err != nil {
			serverError(w, "Error setting clipboard: %s", err)
			return
		}
	default:
		badMethod(w, "GET", "PUT")
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "okay")
}

func requestBody(r *http.Request) (string, error) {
	var b strings.Builder
	if _, err := io.Copy(&b, r.Body); err != nil {		
		return "", fmt.Errorf("read request body: %w", err)
	}
	return b.String(), nil
}

func serverError(w http.ResponseWriter, s string, a ...interface{}) {
	msg := fmt.Sprintf(s, a...)
	log.Printf("Server error: %s", msg)
	w.WriteHeader(500)
	_, err := io.WriteString(w, msg)
	if err != nil {
		log.Printf("Error sending server error: %s", err)
	}
}

func badMethod(w http.ResponseWriter, method ...string) {
	h := w.Header()
	h["Allow"] = method
	w.WriteHeader(405)
}

func getClipboard() (string, error) {
	c := exec.Command("xsel", "-b", "-o")
	o, err := c.Output()
	if err != nil {
		return "", err
	}
	return string(o), nil

}

func setClipboard(s string) error {
	c := exec.Command("xsel", "-b", "-i")
	c.Stdin = strings.NewReader(s)
	return c.Run()
}
