package tester

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/gorilla/mux"
	"github.com/journeybeforedestination/smoke/fhir"
)

type testServer struct{}

// NewFhirServer creates a fhir server that can test app launch
func NewTestServer(addr string) *http.Server {
	s := &testServer{}

	r := mux.NewRouter()

	// setup handlers
	r.HandleFunc("/launch", s.handleLaunch).Methods("GET")

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

// TODO: I would like to make this a seperate server and compose them together with docker
// TODO: figure out how I want to log data here
func (s *testServer) handleLaunch(w http.ResponseWriter, r *http.Request) {
	_, err := fhir.ParseLaunch(r)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	// TODO: don't really want to hard code this...
	resp, err := http.Get("http://launcher:8080/.well-known/smart-configuration")
	if err != nil {
		w.Write([]byte("Error fetching metadata: " + err.Error()))
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
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

	http.Redirect(w, r, conformance.AuthEndpoint, http.StatusFound)
	// both := launchAndMeta{Launch: launch, Meta: conformance}
	// templates.ExecuteTemplate(w, "test.html", both)
}
