package emailverifier

import (
	"net/url"
	"reflect"
	"testing"
)

//TestOptions tests the Options functions
func TestOptions(t *testing.T) {

	tests := []struct {
		name   string
		values url.Values
		option Option
		want   string
	}{
		{
			name:   "output format",
			values: url.Values{},
			option: OptionOutputFormat("JSON"),
			want:   "outputFormat=JSON",
		},
		{
			name:   "hard refresh",
			values: url.Values{},
			option: OptionHardRefresh(1),
			want:   "_hardRefresh=1",
		},
		{
			name:   "validate DNS",
			values: url.Values{},
			option: OptionValidateDNS(0),
			want:   "validateDNS=0",
		},
		{
			name:   "validate SMTP",
			values: url.Values{},
			option: OptionValidateSMTP(1),
			want:   "validateSMTP=1",
		},
		{
			name:   "check CatchAll",
			values: url.Values{},
			option: OptionCheckCatchAll(0),
			want:   "checkCatchAll=0",
		},
		{
			name:   "check Free",
			values: url.Values{},
			option: OptionCheckFree(1),
			want:   "checkFree=1",
		},
		{
			name:   "check Disposable",
			values: url.Values{},
			option: OptionCheckDisposable(0),
			want:   "checkDisposable=0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.option(tt.values)
			if got := tt.values.Encode(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Option() = %v, want %v", got, tt.want)
			}
		})
	}
}
