package emailverifier

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

const (
	pathEvapiResponseOK         = "/Evapi/ok"
	pathEvapiResponseError      = "/Evapi/error"
	pathEvapiResponse500        = "/Evapi/500"
	pathEvapiResponsePartial1   = "/Evapi/partial"
	pathEvapiResponsePartial2   = "/Evapi/partial2"
	pathEvapiResponseUnparsable = "/Evapi/unparsable"
)

const apiKey = "at_LoremIpsumDolorSitAmetConsect"

// dummyServer is the sample of the Email Verification API server for testing
func dummyServer(resp, respUnparsable string, respErr string) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var response string

		response = resp

		switch req.URL.Path {
		case pathEvapiResponseOK:
		case pathEvapiResponseError:
			w.WriteHeader(400)
			response = respErr
		case pathEvapiResponse500:
			w.WriteHeader(500)
			response = respUnparsable
		case pathEvapiResponsePartial1:
			response = response[:len(response)-10]
		case pathEvapiResponsePartial2:
			w.Header().Set("Content-Length", strconv.Itoa(len(response)))
			response = response[:len(response)-10]
		case pathEvapiResponseUnparsable:
			response = respUnparsable
		default:
			panic(req.URL.Path)
		}
		_, err := w.Write([]byte(response))
		if err != nil {
			panic(err)
		}
	}))

	return server
}

// newAPI returns new Email Verification API client for testing
func newAPI(apiServer *httptest.Server, link string) *Client {

	apiURL, err := url.Parse(apiServer.URL)
	if err != nil {
		panic(err)
	}
	apiURL.Path = link

	params := ClientParams{
		HTTPClient:   apiServer.Client(),
		EvapiBaseURL: apiURL,
	}

	return NewClient(apiKey, params)
}

// TestEvapiAPIGet tests the Get function
func TestEvapiGet(t *testing.T) {

	checkResultRec := func(res *EvapiResponse) bool {
		return res != nil
	}

	ctx := context.Background()

	const resp = `{"username":"support","domain":"whoisxmlapi.com","emailAddress":"support@whoisxmlapi.com",
"formatCheck":"true","smtpCheck":"true","dnsCheck":"true","freeCheck":"false","disposableCheck":"false",
"catchAllCheck":"true","mxRecords":["alt1.aspmx.l.google.com.","aspmx2.googlemail.com.","aspmx.l.google.com.",
"aspmx3.googlemail.com.","alt2.aspmx.l.google.com."],"audit":{"auditCreatedDate":"2022-04-03 05:02:37 UTC",
"auditUpdatedDate":"2022-04-03 05:02:37 UTC"}}`

	const respUnparsable = `<?xml version="1.0" encoding="utf-8"?><>`

	const errResp = `{"ErrorMessage":{"Error":"test error message"}}`

	server := dummyServer(resp, respUnparsable, errResp)
	defer server.Close()

	type args struct {
		ctx     context.Context
		options string
	}
	tests := []struct {
		name    string
		path    string
		args    args
		want    bool
		wantErr string
	}{
		{
			name: "successful request",
			path: pathEvapiResponseOK,
			args: args{
				ctx:     ctx,
				options: "support@whoisxmlapi.com",
			},
			want:    true,
			wantErr: "",
		},
		{
			name: "non 200 status code",
			path: pathEvapiResponse500,
			args: args{
				ctx:     ctx,
				options: "support@whoisxmlapi.com",
			},
			want:    false,
			wantErr: "cannot parse response: invalid character '<' looking for beginning of value",
		},
		{
			name: "partial response 1",
			path: pathEvapiResponsePartial1,
			args: args{
				ctx:     ctx,
				options: "support@whoisxmlapi.com",
			},
			want:    false,
			wantErr: "cannot parse response: unexpected EOF",
		},
		{
			name: "partial response 2",
			path: pathEvapiResponsePartial2,
			args: args{
				ctx:     ctx,
				options: "support@whoisxmlapi.com",
			},
			want:    false,
			wantErr: "cannot read response: unexpected EOF",
		},
		{
			name: "could not process request",
			path: pathEvapiResponseError,
			args: args{
				ctx:     ctx,
				options: "support@whoisxmlapi.com",
			},
			want:    false,
			wantErr: "API error: test error message",
		},
		{
			name: "unparsable response",
			path: pathEvapiResponseUnparsable,
			args: args{
				ctx:     ctx,
				options: "support@whoisxmlapi.com",
			},
			want:    false,
			wantErr: "cannot parse response: invalid character '<' looking for beginning of value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			api := newAPI(server, tt.path)

			gotRec, _, err := api.Get(tt.args.ctx, tt.args.options)
			if (err != nil || tt.wantErr != "") && (err == nil || err.Error() != tt.wantErr) {
				t.Errorf("Evapi.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want {
				if !checkResultRec(gotRec) {
					t.Errorf("Evapi.Get() got = %v, expected something else", gotRec)
				}
			} else {
				if gotRec != nil {
					t.Errorf("Evapi.Get() got = %v, expected nil", gotRec)
				}
			}

		})
	}
}

// TestEvapiAPIGetRaw tests the RawData function
func TestEvapiGetRaw(t *testing.T) {

	checkResultRaw := func(res []byte) bool {
		return len(res) != 0
	}

	ctx := context.Background()

	const resp = `{"username":"support","domain":"whoisxmlapi.com","emailAddress":"support@whoisxmlapi.com",
"formatCheck":"true","smtpCheck":"true","dnsCheck":"true","freeCheck":"false","disposableCheck":"false",
"catchAllCheck":"true","mxRecords":["alt1.aspmx.l.google.com.","aspmx2.googlemail.com.","aspmx.l.google.com.",
"aspmx3.googlemail.com.","alt2.aspmx.l.google.com."],"audit":{"auditCreatedDate":"2022-04-03 05:02:37 UTC",
"auditUpdatedDate":"2022-04-03 05:02:37 UTC"}}`

	const respUnparsable = `<?xml version="1.0" encoding="utf-8"?><>`

	const errResp = `{"ErrorMessage":{"Error":"test error message"}}`

	server := dummyServer(resp, respUnparsable, errResp)
	defer server.Close()

	type args struct {
		ctx     context.Context
		options string
	}
	tests := []struct {
		name    string
		path    string
		args    args
		wantErr string
	}{
		{
			name: "successful request",
			path: pathEvapiResponseOK,
			args: args{
				ctx:     ctx,
				options: "support@whoisxmlapi.com",
			},
			wantErr: "",
		},
		{
			name: "non 200 status code",
			path: pathEvapiResponse500,
			args: args{
				ctx:     ctx,
				options: "support@whoisxmlapi.com",
			},
			wantErr: "API failed with status code: 500",
		},
		{
			name: "partial response 1",
			path: pathEvapiResponsePartial1,
			args: args{
				ctx:     ctx,
				options: "support@whoisxmlapi.com",
			},
			wantErr: "",
		},
		{
			name: "partial response 2",
			path: pathEvapiResponsePartial2,
			args: args{
				ctx:     ctx,
				options: "support@whoisxmlapi.com",
			},
			wantErr: "cannot read response: unexpected EOF",
		},
		{
			name: "unparsable response",
			path: pathEvapiResponseUnparsable,
			args: args{
				ctx:     ctx,
				options: "support@whoisxmlapi.com",
			},
			wantErr: "",
		},
		{
			name: "could not process request",
			path: pathEvapiResponseError,
			args: args{
				ctx:     ctx,
				options: "support@whoisxmlapi.com",
			},
			wantErr: "API failed with status code: 400",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			api := newAPI(server, tt.path)

			resp, err := api.GetRaw(tt.args.ctx, tt.args.options)
			if (err != nil || tt.wantErr != "") && (err == nil || err.Error() != tt.wantErr) {
				t.Errorf("Evapi.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !checkResultRaw(resp.Body) {
				t.Errorf("Evapi.Get() got = %v, expected something else", string(resp.Body))
			}
		})
	}
}
