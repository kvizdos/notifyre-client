package examples

import (
	"fmt"
	"os"
	"time"

	"github.com/kvizdos/notifyre-client/fax"
)

func ListFaxes() {
	listResp, err := fax.ListSentFaxes(os.Getenv("NOTIFYRE_KEY"), fax.ListParameters{
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
