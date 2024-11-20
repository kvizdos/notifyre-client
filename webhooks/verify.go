package notifyre_webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func VerifySignature(header string, payload []byte, secret string) bool {
	// Extract timestamp and signature
	elements := strings.Split(header, ",")
	var timestamp, signature string
	for _, element := range elements {
		parts := strings.Split(element, "=")
		if len(parts) == 2 {
			switch parts[0] {
			case "t":
				timestamp = parts[1]
			case "v":
				signature = parts[1]
			}
		}
	}

	if timestamp == "" || signature == "" {
		return false
	}

	// Generate signed payload
	signedPayload := timestamp + "." + string(payload)

	// Compute expected signature
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signedPayload))
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	// Compare signatures
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
