// server/main.go
package server

import (
	"log"
	"net/http"

	"github.com/matty271828/flightalerts/gf-emailparser/internal/api"
	internalGoogle "github.com/matty271828/flightalerts/gf-emailparser/internal/google"
)

type Server struct {
	API *api.API
}

const (
	CallBackURL   = "/callback"
	ReadEmailsURL = "/reademails"
)

// NewServer is used to initialise a new instance of the server
func NewServer(
	api *api.API,
) error {
	server := Server{
		API: api,
	}

	server.RegisterEndpoints()
	log.Println("Successfully initialised full web server on :8080")
	return nil
}

func Start() {
	http.HandleFunc(CallBackURL, internalGoogle.HandleGoogleCallback)
	log.Println("Starting web server on :8080")
	http.ListenAndServe(":8080", nil)
}

// RegisterEndpoints is used to register endpoints onto the web server
func (s *Server) RegisterEndpoints() {
	http.HandleFunc(ReadEmailsURL, s.API.ReadEmails)
}
