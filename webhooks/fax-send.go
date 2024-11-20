package notifyre_webhook

import (
	"encoding/json"
	"time"
)

type FaxEvent string

const (
	FAX_RECEIVED FaxEvent = "fax_received"
	FAX_SENT     FaxEvent = "fax_sent"
	SMS_RECEIVED FaxEvent = "sms_received"
	SMS_SENT     FaxEvent = "sms_sent"
	MMS_RECEIVED FaxEvent = "mms_received"
)

// FaxStatus represents the status of the fax document
type FaxStatus string

const (
	StatusAccepted   FaxStatus = "accepted"
	StatusSuccessful FaxStatus = "successful"
	StatusFailed     FaxStatus = "failed"
	StatusInProgress FaxStatus = "in_progress"
	StatusQueued     FaxStatus = "queued"
)

// WebhookEvent represents the structure of the webhook event
type FaxSentWebhookEvent struct {
	Event     FaxEvent   `json:"Event"`
	Timestamp UnixTime   `json:"Timestamp"`
	Payload   FaxPayload `json:"Payload"`
}

// UnixTime is a custom type for handling Unix timestamps in JSON
type UnixTime time.Time

// UnmarshalJSON converts a Unix timestamp into a time.Time
func (ut *UnixTime) UnmarshalJSON(data []byte) error {
	// Parse the JSON number (Unix timestamp)
	var unixTimestamp int64
	if err := json.Unmarshal(data, &unixTimestamp); err != nil {
		return err
	}

	// Convert the Unix timestamp to time.Time
	*ut = UnixTime(time.Unix(unixTimestamp, 0))
	return nil
}

func (ut UnixTime) MarshalJSON() ([]byte, error) {
	return time.Time(ut).MarshalJSON()
}

// FaxPayload represents the payload of a fax event
type FaxPayload struct {
	ID                   string    `json:"ID"`
	FriendlyID           string    `json:"FriendlyID"`
	RecipientID          string    `json:"RecipientID"`
	FromNumber           string    `json:"FromNumber"`
	To                   string    `json:"To"`
	Reference            string    `json:"Reference"`
	CreatedDateUtc       *UnixTime `json:"CreatedDateUtc"`
	QueuedDateUtc        *UnixTime `json:"QueuedDateUtc"`
	LastModifiedDateUtc  *UnixTime `json:"LastModifiedDateUtc"`
	HighQuality          bool      `json:"HighQuality"`
	Pages                int       `json:"Pages"`
	Status               FaxStatus `json:"Status"`
	AssumedFailureReason string    `json:"AssumedFailureReason"` /* Filled using the rough approach documented in readme.md */
}
