package fax_response

import "fmt"

type FailureResponse struct {
	Type       string `json:"type"`
	Title      string `json:"title"`
	StatusCode int    `json:"status"`
	TraceID    string `json:"traceId"`
}

func (f FailureResponse) Successful() bool {
	return false
}

func (f FailureResponse) Error() string {
	return fmt.Sprintf("%s (%d) - %s - More info: %s", f.Title, f.StatusCode, f.TraceID, f.Type)
}
