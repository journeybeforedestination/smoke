package launcher

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"text/template"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/journeybeforedestination/smoke/fhir"
)

var templates = template.Must(template.ParseFiles("../../launcher/tmpl/home.html", "../../launcher/tmpl/conformance.json"))

type fhirServer struct{}

// NewFhirServer creates a fhir server that can test app launch
func NewFhirServer(addr string) *http.Server {
	s := &fhirServer{}

	r := mux.NewRouter()

	// setup handlers
	r.HandleFunc("/", s.handleRoot).Methods("GET")
	r.HandleFunc("/test", s.handleTest).Methods("GET")
	r.HandleFunc("/launch", s.handleLaunch).Methods("POST")
	r.HandleFunc("/authorize", s.handleAuth).Methods("GET")

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
func (s *fhirServer) handleRoot(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "home.html", nil)
}

// TODO: I would like to make this a seperate server and compose them together with docker
// TODO: figure out how I want to log data here
func (s *fhirServer) handleTest(w http.ResponseWriter, r *http.Request) {
	launch, err := fhir.ParseLaunch(r)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	resp, err := http.Get("http://" + launch.ISS + "/.well-known/smart-configuration")
	if err != nil {
		w.Write([]byte("Error fetching metadata"))
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

// handleLaunch handles sumbission of a launch
func (s *fhirServer) handleLaunch(w http.ResponseWriter, r *http.Request) {
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
func (s *fhirServer) handleConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	templates.ExecuteTemplate(w, "conformance.json", nil)
}

// handleAuth performs a very limited test OAuth2 check
func (s *fhirServer) handleAuth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Authorizing...</h1>"))
}
