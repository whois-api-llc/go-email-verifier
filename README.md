[![go-email-verifier license](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)
[![go-email-verifier made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](https://pkg.go.dev/github.com/whois-api-llc/go-email-verifier)
[![go-email-verifier test](https://github.com/whois-api-llc/go-email-verifier/workflows/Test/badge.svg)](https://github.com/whois-api-llc/go-email-verifier/actions/)

# Overview

The client library for
[Email Verification API](https://emailverification.whoisxmlapi.com)
in Go language.

The minimum go version is 1.17.

# Installation

The library is distributed as a Go module

```bash
go get github.com/whois-api-llc/go-email-verifier
```

# Examples

Full API documentation available [here](https://emailverification.whoisxmlapi.com/api/documentation/making-requests)

You can find all examples in `example` directory.

## Create a new client

To start making requests you need the API Key. 
You can find it on your profile page on [whoisxmlapi.com](https://whoisxmlapi.com/).
Using the API Key you can create Client.

Most users will be fine with `NewBasicClient` function. 
```go
client := emailverifier.NewBasicClient(apiKey)
```

If you want to set custom `http.Client` to use proxy then you can use `NewClient` function.
```go
transport := &http.Transport{Proxy: http.ProxyURL(proxyUrl)}

client := emailverifier.NewClient(apiKey, emailverifier.ClientParams{
    HTTPClient: &http.Client{
        Transport: transport,
        Timeout:   20 * time.Second,
    },
})
```

## Make basic requests

Email Verification API performs a comprehensive validation of email addresses in real-time and conveniently. 

```go

// Make request to get parsed Email Verification API response
evapiResp, _, err := client.EvapiService.Get(context.Background(), "support@whoisxmlapi.com")
if err != nil {
    log.Fatal(err)
}

// Check if an email address is valid
if !*evapiResp.FormatCheck {
    log.Printf("\"%s\" is invalid email address", evapiResp.EmailAddress)
}

// Make request to get raw Email Verification API data
resp, err := client.EvapiService.GetRaw(ctx, "whoisxmlapi.com")
if err != nil {
    log.Fatal(err)
}

log.Println(string(resp.Body))


```
