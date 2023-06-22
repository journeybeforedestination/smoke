package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"text/template"

	"github.com/gorilla/mux"
)

var templates = template.Must(template.ParseFiles("tmpl/home.html", "tmpl/test.html", "tmpl/iFrame.html"))

type server struct{}

// NewHTTPServer creates an http server
func NewHTTPServer(addr string) *http.Server {
	s := &server{}

	r := mux.NewRouter()

	// setup handlers
	r.HandleFunc("/", s.handleRoot).Methods("GET")
	r.HandleFunc("/test", s.handleTest).Methods("GET")
	r.HandleFunc("/launch", s.handleLaunch).Methods("POST")

	// register middleware
	var handler http.Handler = r
	handler = logRequestHandler(handler)

	return &http.Server{
		Addr:    addr,
		Handler: handler,
	}
}

// logRequestHandler is a middleware that writes an http log for each request
func logRequestHandler(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)

		requestDump, err := httputil.DumpRequest(r, true)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Failed to dump request:", err)
			return
		}

		log.Println("Received request:")
		log.Println(string(requestDump))
	}
	return http.HandlerFunc(fn)
}

// handleRoot handles the root path "/"
func (s *server) handleRoot(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "home.html", nil)
}

func (s *server) handleTest(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "test.html", nil)
}

// handleLaunch handles sumbission of a launch
func (s *server) handleLaunch(w http.ResponseWriter, r *http.Request) {
	launchURL := r.FormValue("launch-url")

	// Check if the launch URL is valid
	_, err := url.ParseRequestURI(launchURL)
	if err != nil {
		w.Write([]byte("Invalid URL"))
		return
	}
	iFrame := fmt.Sprintf(`<iframe id="launch-iframe" src="%s"></iframe>`, launchURL)

	w.Write([]byte(iFrame))
	//templates.ExecuteTemplate(w, "iframe.html", launchURL)
}
