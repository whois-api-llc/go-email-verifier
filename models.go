package emailverifier

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// unmarshalString parses the JSON-encoded data and returns value as a string
func unmarshalString(raw json.RawMessage) (string, error) {
	var val string
	err := json.Unmarshal(raw, &val)
	if err != nil {
		return "", err
	}
	return val, nil
}

// StringBool is a helper wrapper on bool
type StringBool bool

// UnmarshalJSON decodes true/false values from Email Verification API
func (b *StringBool) UnmarshalJSON(bytes []byte) error {
	str, err := unmarshalString(bytes)
	if err != nil {
		return err
	}

	*b = str == "true" || str == "1"
	return nil
}

// MarshalJSON encodes true/false values to the Email Verification API representation
func (b StringBool) MarshalJSON() ([]byte, error) {
	return []byte(`"` + strconv.FormatBool(bool(b)) + `"`), nil
}

// Time is a helper wrapper on time.Time
type Time time.Time

var emptyTime Time

// UnmarshalJSON decodes time as Email Verification API does
func (t *Time) UnmarshalJSON(b []byte) error {
	str, err := unmarshalString(b)
	if err != nil {
		return err
	}
	if str == "" {
		*t = emptyTime
		return nil
	}
	v, err := time.Parse("2006-01-02 15:04:05 MST", str)
	if err != nil {
		return err
	}
	*t = Time(v)
	return nil
}

// MarshalJSON encodes time as Email Verification API does
func (t Time) MarshalJSON() ([]byte, error) {
	if t == emptyTime {
		return []byte(`""`), nil
	}
	return []byte(`"` + time.Time(t).Format("2006-01-02 15:04:05 MST") + `"`), nil
}

// Audit is part of the Email Verification API response
// It represents dates when data was added and updated in our database
type Audit struct {
	// AuditCreatedDate is the date this data is collected on whoisxmlapi.com
	AuditCreatedDate Time `json:"auditCreatedDate"`

	// AuditUpdatedDate is the date this data is updated on whoisxmlapi.com
	AuditUpdatedDate Time `json:"auditUpdatedDate"`
}

// EvapiResponse is a response of Email Verification API
type EvapiResponse struct {
	// Username is a username
	Username string `json:"username"`

	// Domain is a domain name
	Domain string `json:"domain"`

	// EmailAddress is an email address
	EmailAddress string `json:"emailAddress"`

	// FormatCheck indicates if there are any syntax errors in the email address
	FormatCheck *StringBool `json:"formatCheck"`

	// SmtpCheck indicates if the email address exists and can receive emails by using SMTP connection and
	// email-sending emulation techniques
	SmtpCheck *StringBool `json:"smtpCheck"`

	// DnsCheck ensures that the domain in the email address is a valid domain
	DnsCheck *StringBool `json:"dnsCheck"`

	// FreeCheck indicates if the email address is from a free email provider
	FreeCheck *StringBool `json:"freeCheck"`

	// DisposableCheck tells you whether the email address is disposable
	DisposableCheck *StringBool `json:"disposableCheck"`

	// CatchAllCheck tells you whether the related mail server has a "catch-all" address
	CatchAllCheck *StringBool `json:"catchAllCheck"`

	// MxRecords is a mail servers list
	MxRecords []string `json:"mxRecords"`

	// Audit is a data update dates
	Audit Audit `json:"audit"`
}

// ErrorMessage is an error message
type ErrorMessage struct {
	// Message is an error message
	Message string `json:"Error"`
}

// Error returns error message as a string
func (e ErrorMessage) Error() string {
	return fmt.Sprintf("API error: %s", e.Message)
}
