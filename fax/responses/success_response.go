package fax_response

type SuccessResponsePayload struct {
	FaxID      string `json:"faxID"`
	FriendlyID string `json:"friendlyID"`
}

type SuccessResponse struct {
	Payload    SuccessResponsePayload `json:"payload"`
	Success    bool                   `json:"success"`
	StatusCode int                    `json:"statusCode"`
}
