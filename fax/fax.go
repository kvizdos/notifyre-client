package fax

type Document struct {
	Filename string `json:"Filename"`
	Data     string `json:"Data"`
}

type RecipientType string

const (
	FAX_NUMBER RecipientType = "fax_number"
	CONTACT    RecipientType = "contact"
	GROUP      RecipientType = "group"
)

type Recipient struct {
	Type  RecipientType `json:"Type"`
	Value string        `json:"Value"`
}

type Fax struct {
	Recipients      []Recipient `json:"Recipients"`
	SendFrom        string      `json:"SendFrom"`
	ClientReference string      `json:"ClientReference"`
	Subject         string      `json:"Subject"`
	IsHighQuality   bool        `json:"IsHighQuality"`
	Documents       []Document  `json:"Documents"`
}

type Payload struct {
	Faxes Fax `json:"Faxes"`
}
