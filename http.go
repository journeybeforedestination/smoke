package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"text/template"

	"github.com/google/uuid"
	"github.com/journeybeforedestination/smoke/fhir"

	"github.com/gorilla/mux"
)

var templates = template.Must(template.ParseFiles("tmpl/home.html", "tmpl/test.html", "tmpl/conformance.json"))

type server struct{}

// NewHTTPServer creates an http server
func NewHTTPServer(addr string) *http.Server {
	s := &server{}

	r := mux.NewRouter()

	// setup handlers
	r.HandleFunc("/", s.handleRoot).Methods("GET")
	r.HandleFunc("/test", s.handleTest).Methods("GET")
	r.HandleFunc("/launch", s.handleLaunch).Methods("POST")

	r.HandleFunc("/.well-known/smart-configuration", s.handleConfig).Methods("GET")

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

type launchAndMeta struct {
	Launch fhir.Launch
	Meta   fhir.Conformance
}

// TODO: I would like to make this a seperate server and compose them together with docker
func (s *server) handleTest(w http.ResponseWriter, r *http.Request) {
	iss, err := url.QueryUnescape(r.FormValue("iss"))
	if err != nil {
		w.Write([]byte("invalid iss"))
	}
	l, err := url.QueryUnescape(r.FormValue("launch"))
	if err != nil {
		w.Write([]byte("invalid launch"))
	}
	launch := fhir.Launch{Launch: l, ISS: iss}

	resp, err := http.Get("http://" + launch.ISS + "/.well-known/smart-configuration")
	if err != nil {
		w.Write([]byte("Error fetching metadata"))
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		w.Write([]byte("Error reading metadata"))
		return
	}

	var conformance fhir.Conformance
	err = json.Unmarshal([]byte(body), &conformance)
	if err != nil {
		w.Write([]byte("Error parsing metadata"))
		return
	}

	both := launchAndMeta{Launch: launch, Meta: conformance}
	templates.ExecuteTemplate(w, "test.html", both)
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

	l := uuid.NewString()
	iss := r.Host
	launch := fhir.Launch{Launch: l, ISS: iss}
	launchURL = fmt.Sprintf("%s?%s", launchURL, launch.Params())

	iFrame := fmt.Sprintf(`<iframe id="launch-iframe" src="%s"></iframe>`, launchURL)

	w.Write([]byte(iFrame))
	//templates.ExecuteTemplate(w, "iframe.html", launchURL)
}

// handle well known config
func (s *server) handleConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	templates.ExecuteTemplate(w, "conformance.json", nil)
}
