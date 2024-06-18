package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/kvizdos/notifyre-client/fax"
)

func main() {
	apiKey := "<<API KEY>>"

	// Read file from disk
	filePath := "test.pdf"
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Encode file data to base64
	encodedFileData := base64.StdEncoding.EncodeToString(fileData)

	resp, err := fax.Send(fax.Payload{
		Faxes: fax.Fax{
			Recipients: []fax.Recipient{
				{
					Type:  fax.FAX_NUMBER,
					Value: "+14342324045",
				},
			},
			SendFrom:        "+14342324045",
			ClientReference: "This is a test fax.",
			Subject:         "This is a test fax.",
			IsHighQuality:   false,
			Documents: []fax.Document{
				{
					Filename: "test.pdf",
					Data:     encodedFileData,
				},
			},
		},
	}, apiKey)

	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", resp)

	listResp, err := fax.ListSentFaxes(apiKey, fax.ListParameters{
		FromDate: time.Now().AddDate(0, 0, -30).UTC(),
		ToDate:   time.Now().UTC(),
		Sort:     fax.ASCENDING,
		Limit:    10,
		Skip:     0,
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", listResp)
}
