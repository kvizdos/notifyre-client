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

	v.Set("statusType", string(p.StatusType))

	v.Set("fromDate", fmt.Sprintf("%d", p.FromDate.Unix()))
	v.Set("toDate", fmt.Sprintf("%d", p.ToDate.Unix()))
	if p.Sort == "" {
		p.Sort = ASCENDING
	}
	v.Set("sort", string(p.Sort))

	v.Set("limit", fmt.Sprintf("%d", p.Limit))
	v.Set("skip", fmt.Sprintf("%d", p.Skip))

	log.Println(v.Encode())
	return v.Encode()
}

type SentFaxesListItem struct {
	ID                  string         `json:"ID"`
	FriendlyID          string         `json:"FriendlyID"`
	RecipientID         string         `json:"RecipientID"`
	FromNumber          string         `json:"FromNumber"`
	ToNumber            string         `json:"To"`
	Reference           string         `json:"Reference"`
	CreatedDateUtc      int64          `json:"CreatedDateUtc"`
	QueuedDateUtc       int64          `json:"QueuedDateUtc,omitempty"`
	LastModifiedDateUtc int64          `json:"LastModifiedDateUtc"`
	HighQuality         bool           `json:"HighQuality"`
	Pages               int            `json:"Pages"`
	Status              ListStatusType `json:"Status"`
	FailedMessage       string         `json:"FailedMessage"`
}

type SentFaxesListResponse struct {
	Success    bool     `json:"Success"`
	StatusCode int      `json:"StatusCode"`
	Message    string   `json:"Message"`
	Errors     []string `json:"Errors"`
	Payload    struct {
		Faxes []SentFaxesListItem `json:"Faxes"`
	} `json:"Payload"`
}

func ListSentFaxes(apiKey string, payload ListParameters) (SentFaxesListResponse, error) {
	if payload.FromDate.IsZero() || payload.ToDate.IsZero() {
		return SentFaxesListResponse{}, fmt.Errorf("requires payload.FromDate and payload.ToDate")
	}

	if payload.Limit == 0 {
		payload.Limit = 20
	}

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

	if apiResponse.StatusCode != 200 {
		return apiResponse, fmt.Errorf("notifyre bad status: %d %s", apiResponse.StatusCode, apiResponse.Message)
	}

	return apiResponse, nil
}
