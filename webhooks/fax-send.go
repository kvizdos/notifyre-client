package webhooks

type FaxEvent string

const (
	FAX_RECEIVED FaxEvent = "fax_received"
	FAX_SENT     FaxEvent = "fax_sent"
	SMS_RECEIVED FaxEvent = "sms_received"
	SMS_SENT     FaxEvent = "sms_sent"
	MMS_RECEIVED FaxEvent = "mms_received"
)

type FaxSendWebhook struct {
	Event     FaxEvent `json:"Event"`
	Timestamp int64    `json:"Timestamp"`
	// Payload
}

type FaxSendWebhookPayload struct {
}
