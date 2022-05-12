package main

import (
	"context"
	"errors"
	emailverifier "github.com/whois-api-llc/go-email-verifier"
	"log"
	"os"
	"strconv"
	"time"
)

func GetData(apikey string) {
	client := emailverifier.NewBasicClient(apikey)
	/*
		client := emailverifier.NewClient(apikey, emailverifier.ClientParams{
			HTTPClient: &http.Client{
				Transport: nil,
				Timeout:   10 * time.Second,
			},
		})
	*/
	// Get parsed Email Verification API response as a model instance
	evapiResp, resp, err := client.EvapiService.Get(context.Background(), "support@whoisxmlapi.com",
		// this option is ignored, as the inner parser works with JSON only
		emailverifier.OptionOutputFormat("XML"),
		// this option results the catchAll check is omitted
		emailverifier.OptionCheckCatchAll(1))

	if err != nil {
		// Handle error message returned by server
		var apiErr *emailverifier.ErrorMessage
		if errors.As(err, &apiErr) {
			log.Println(apiErr.Message)
		}
		log.Fatal(err)
	}

	if !*evapiResp.FormatCheck {
		log.Printf("\"%s\" is invalid email address", evapiResp.EmailAddress)
	}

	if evapiResp.DnsCheck == nil || !*evapiResp.DnsCheck {
		log.Printf("\"%s\" is invalid domain name", evapiResp.Domain)
	}

	//Some values are not always returned and need to be validated before printing
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

func main() {
	apikey := os.Getenv("APIKEY")
	if apikey == "" {
		log.Fatal("Empty API KEY")
	}

	GetData(apikey)
	//GetRawData(apikey)

}
