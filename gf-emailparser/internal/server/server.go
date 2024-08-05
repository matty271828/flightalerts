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
	ReadEmailsURL = "/ReadEmails"
)

// NewServer is used to initialise a new instance of the server
func NewServer(
	api *api.API,
) (*Server, error) {
	return &Server{
		API: api,
	}, nil
}

func Start() {
	http.HandleFunc(CallBackURL, internalGoogle.HandleGoogleCallback)
	log.Println("Starting web server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// RegisterEndpoints is used to register endpoints onto the web server
func (s *Server) RegisterEndpoints() {
	http.HandleFunc(ReadEmailsURL, s.API.ReadEmails)
}
