package webhooks

import "net/http"

type Webhook struct{}

func (w Webhook) FaxSendPayload(r *http.Request) {

}
