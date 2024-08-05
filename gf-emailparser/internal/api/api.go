package api

import (
	"net/http"

	jobs "github.com/matty271828/flightalerts/gf-emailparser/internal/jobs"
)

type API struct {
	Jobs *jobs.Jobs
}

// NewAPI is used to initialise a new instance of the API service.
func NewAPI(
	jobs *jobs.Jobs,
) (*API, error) {
	return &API{
		Jobs: jobs,
	}, nil
}

// respond is used to construct a header and write a response to an
// endpoint call.
func (a *API) respond(w http.ResponseWriter, statusCode int, message string) error {
	// Set the response status code to 202 Accepted
	w.WriteHeader(statusCode)

	// Write the response body with the message
	_, err := w.Write([]byte(message))
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
	return err
}

// ReadEmails is used to action a job to read all unread flight alerts
// currently in the inbox and extract their contents.
func (a *API) ReadEmails(w http.ResponseWriter, r *http.Request) {
	err := a.Jobs.ReadEmailsJob()
	if err != nil {
		http.Error(w, "Failed to run ReadEmails", http.StatusInternalServerError)
	}

	a.respond(w, http.StatusAccepted, "ReadEmails request accepted")
}
