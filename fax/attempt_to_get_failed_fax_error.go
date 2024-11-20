package fax

import (
	"fmt"
	"time"
)

/*
**THIS API WILL KILL YOU AT SCALE**

Sadly, Notifyre does *not* have an API to query Status directly given a fax ID.

Support recommended the following:
- Query the List API by Status = Failing
- Include the CreatedDateUtc as the FromDate and add about 10 minutes for ToDate
- Find the Fax ID in the returned list

This feels gross. I requested a new "forId" field to be added to this endpoint, or better, a direct API route for this purpose.
*/

func AttemptToGetFailedFaxError(apiKey string, lookingForFaxID string, createdAtTimestamp time.Time) (*SentFaxesListItem, error) {
	// Fetch the list of failed faxes
	responses, err := ListSentFaxes(apiKey, ListParameters{
		StatusType: "failed",
		FromDate:   createdAtTimestamp.Add(-1 * time.Minute),
		ToDate:     createdAtTimestamp.Add(10 * time.Minute),
		Limit:      20, /* This can't possibly help with the "Murder at Scale" problem, but I don't want to search through hundreds of possible events.. */
		Skip:       0,
	})

	// Handle errors from the API request
	if err != nil {
		return nil, fmt.Errorf("failed to list sent faxes: %w", err)
	}

	// Check for API-level errors
	if !responses.Success {
		if len(responses.Errors) > 0 {
			return nil, fmt.Errorf("API request failed with errors: %v", responses.Errors)
		}
		return nil, fmt.Errorf("API request failed without specific errors")
	}

	// Look for the fax ID in the response payload
	for _, response := range responses.Payload.Faxes {
		if response.ID == lookingForFaxID {
			return &response, nil
		}
	}

	return nil, fmt.Errorf("fax id not found")
}
