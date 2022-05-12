package example

import (
	"context"
	"errors"
	emailverifier "github.com/whois-api-llc/go-email-verifier"
	"log"
	"strconv"
	"time"
)

func GetData(apikey string) {
	client := emailverifier.NewBasicClient(apikey)

	// Get parsed Email Verification API response as a model instance
	evapiResp, resp, err := client.EvapiService.Get(context.Background(), "support@whoisxmlapi.com",
		// this option is ignored, as the inner parser works with JSON only
		emailverifier.OptionOutputFormat("XML"),
		// this option results in the catchAll check being omitted
		emailverifier.OptionCheckCatchAll(0))

	if err != nil {
		// Handle error message returned by server
		var apiErr *emailverifier.ErrorMessage
		if errors.As(err, &apiErr) {
			log.Println(apiErr.Message)
		}
		log.Fatal(err)
	}

	//Some values are not always returned and need to be validated before using
	if evapiResp.CatchAllCheck != nil {
		log.Printf("emailAddress: %s, catchAll: %s\n",
			evapiResp.EmailAddress,
			strconv.FormatBool(bool(*evapiResp.CatchAllCheck)))
	}

	if evapiResp.SmtpCheck != nil {
		log.Printf("emailAddress: %s, audit.updatedDate: %s, smtpCheck: %s\n",
			evapiResp.EmailAddress,
			time.Time(evapiResp.Audit.AuditUpdatedDate).Format("2006-01-02 15:04:05 MST"),
			strconv.FormatBool(bool(*evapiResp.SmtpCheck)))
	}

	log.Println("raw response is always in JSON format. Most likely you don't need it.")
	log.Printf("raw response: %s\n", string(resp.Body))
}

func GetRawData(apikey string) {
	client := emailverifier.NewBasicClient(apikey)

	// Get raw API response
	resp, err := client.EvapiService.GetRaw(context.Background(), "support@whoisxmlapi.com",
		emailverifier.OptionOutputFormat("JSON"),
		emailverifier.OptionHardRefresh(1))

	if err != nil {
		// Handle error message returned by server
		log.Fatal(err)
	}

	log.Println(string(resp.Body))
}
