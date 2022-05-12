package emailverifier

import (
	"net/url"
	"strconv"
	"strings"
)

// Option adds parameters to the query
type Option func(v url.Values)

var _ = []Option{
	OptionOutputFormat("JSON"),
	OptionHardRefresh(0),
	OptionValidateDNS(0),
	OptionValidateSMTP(0),
	OptionCheckCatchAll(0),
	OptionCheckFree(0),
	OptionCheckDisposable(0),
}

// OptionOutputFormat to set Response output format JSON | XML. Default: JSON.
func OptionOutputFormat(outputFormat string) Option {
	return func(v url.Values) {
		v.Set("outputFormat", strings.ToUpper(outputFormat))
	}
}

// OptionHardRefresh to set parameter for getting fresh data. Default: 0.
func OptionHardRefresh(value int) Option {
	return func(v url.Values) {
		v.Set("_hardRefresh", strconv.Itoa(value))
	}
}

// OptionValidateDNS to set parameter for checking the email address with DNS. Default: 1.
func OptionValidateDNS(value int) Option {
	return func(v url.Values) {
		v.Set("validateDNS", strconv.Itoa(value))
	}
}

// OptionValidateSMTP to set parameter for checking the email address with SMTP. Default: 1.
func OptionValidateSMTP(value int) Option {
	return func(v url.Values) {
		v.Set("validateSMTP", strconv.Itoa(value))
	}
}

// OptionCheckCatchAll to set parameter for checking if the email provider has a catch-all email address. Default: 1.
func OptionCheckCatchAll(value int) Option {
	return func(v url.Values) {
		v.Set("checkCatchAll", strconv.Itoa(value))
	}
}

// OptionCheckFree to set parameter for checking whether the email provider is a free one. Default: 1.
func OptionCheckFree(value int) Option {
	return func(v url.Values) {
		v.Set("checkFree", strconv.Itoa(value))
	}
}

// OptionCheckDisposable to set parameter for checking if the address is disposable. Default: 1.
func OptionCheckDisposable(value int) Option {
	return func(v url.Values) {
		v.Set("checkDisposable", strconv.Itoa(value))
	}
}
