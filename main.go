package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	notifyre_webhook "github.com/kvizdos/notifyre-client/webhooks"
)

func demoHandler(w http.ResponseWriter, r *http.Request) {
	event, ok := notifyre_webhook.GetNotifyreEvent(r.Context())
	if !ok {
		http.Error(w, "Event not found in context", http.StatusInternalServerError)
		return
	}

	js, _ := json.Marshal(event)

	log.Println(string(js))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Event processed"))
}

func main() {
	godotenv.Load()

	// Initialize NotifyreWebhook middleware
	notifyreMiddleware := notifyre_webhook.NotifyreFaxSentWebhook{
		SecretKey: os.Getenv("NOTIFYRE_WEBHOOK_SECRET"),
		ApiKey:    os.Getenv("NOTIFYRE_KEY"),
	}

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Register the middleware and the demoHandler for "/demo/webhook"
	mux.Handle("/demo/webhook", notifyreMiddleware.ServeHTTP(http.HandlerFunc(demoHandler)))

	// Start the HTTP server
	port := ":8003"
	fmt.Printf("Server running on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, mux))
}
