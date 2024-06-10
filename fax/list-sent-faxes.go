package fax

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type ListStatusType string

const (
	ACCEPTED    ListStatusType = "accepted"
	SUCCESSFUL  ListStatusType = "successful"
	IN_PROGRESS ListStatusType = "in_progress"
	FAILED      ListStatusType = "failed"
	QUEUED      ListStatusType = "queued"
)

type ListOrdering string

const (
	ASCENDING  ListOrdering = "asc"
	DESCENDING ListOrdering = "desc"
)

type ListParameters struct {
	StatusType ListStatusType `json:"statusType"`
	FromDate   time.Time      `json:"fromDate"`
	ToDate     time.Time      `json:"toDate"`
	Sort       ListOrdering   `json:"sort"`
	Limit      int            `json:"limit"`
	Skip       int            `json:"skip"`
}

func (p *ListParameters) ToQueryString() string {
	v := url.Values{}
	if p.StatusType != "" {
		v.Set("statusType", string(p.StatusType))
	}
	v.Set("fromDate", fmt.Sprintf("%d", p.FromDate.Unix()))
	v.Set("toDate", fmt.Sprintf("%d", p.ToDate.Unix()))
	v.Set("sort", string(p.Sort))
	v.Set("limit", fmt.Sprintf("%d", p.Limit))
	v.Set("skip", fmt.Sprintf("%d", p.Skip))
	return v.Encode()
}

type SentFaxesListItem struct {
	ID                  string         `json:"id"`
	FriendlyID          string         `json:"friendlyID"`
	RecipientID         string         `json:"recipientID"`
	FromNumber          string         `json:"fromNumber"`
	ToNumber            string         `json:"to"`
	Reference           string         `json:"reference"`
	CreatedDateUtc      int64          `json:"createdDateUtc"`
	QueuedDateUtc       int64          `json:"queuedDateUtc,omitempty"`
	LastModifiedDateUtc int64          `json:"lastModifiedDateUtc"`
	HighQuality         bool           `json:"highQuality"`
	Pages               int            `json:"pages"`
	Status              ListStatusType `json:"status"`
	FailedMessage       string         `json:"failedMessage"`
}

type SentFaxesListResponse struct {
	Success    bool     `json:"success"`
	StatusCode int      `json:"statusCode"`
	Message    string   `json:"message"`
	Errors     []string `json:"errors"`
	Payload    struct {
		Faxes []SentFaxesListItem `json:"faxes"`
	}
}

func ListSentFaxes(apiKey string, payload ListParameters) (SentFaxesListResponse, error) {
	// Create request
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.notifyre.com/fax/send?%s", payload.ToQueryString()), nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	// Set headers
	req.Header = http.Header{
		"x-api-token":  {apiKey},
		"Content-Type": {"application/json"},
	}

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	var apiResponse SentFaxesListResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return SentFaxesListResponse{}, err
	}

	return apiResponse, nil
}
