package examples

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/kvizdos/notifyre-client/fax"
)

func SendFax(apiKey string) {
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
					Value: "+143423240451",
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
}
