package notifyre_webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/kvizdos/notifyre-client/fax"
)

// GetNotifyreEvent retrieves the FaxWebhookEvent from the context
func GetNotifyreEvent(ctx context.Context) (FaxSentWebhookEvent, bool) {
	event, ok := ctx.Value("event").(FaxSentWebhookEvent)
	return event, ok
}

type NotifyreFaxSentWebhook struct {
	SecretKey string
	ApiKey    string
}

func (webhook NotifyreFaxSentWebhook) ServeHTTP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		userAgent := r.Header.Get("user-agent")
		notifyreSignature := r.Header.Get("notifyre-signature")

		if userAgent != "Notifyre" || notifyreSignature == "" {
			http.Error(w, "You Are Not Allowed", http.StatusForbidden)
			return
		}

		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		if !VerifySignature(notifyreSignature, body, webhook.SecretKey) {
			http.Error(w, "Bad Signature", http.StatusUnauthorized)
			return
		}

		// Parse the Notifyre event
		var event FaxSentWebhookEvent
		if err := json.Unmarshal(body, &event); err != nil {
			log.Println(err.Error())
			http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
			return
		}

		if event.Payload.Status == StatusFailed {
			foundDetails, err := fax.AttemptToGetFailedFaxError(webhook.ApiKey, event.Payload.ID, time.Time(*event.Payload.CreatedDateUtc))

			if err != nil {
				if err.Error() == "fax id not found" {
					event.Payload.AssumedFailureReason = "Unable to detect failure reason"
				} else {
					event.Payload.AssumedFailureReason = fmt.Sprintf("Failed to find failure reason: %s", err.Error())
				}
			} else {
				event.Payload.AssumedFailureReason = foundDetails.FailedMessage
			}
		}

		// Add the parsed event to the request context
		r = r.WithContext(context.WithValue(r.Context(), "event", event))

		next.ServeHTTP(w, r)
	})
}
