package emailverifier

import (
	"encoding/json"
	"testing"
)

//TestTime tests JSON encoding/parsing functions for the time values
func TestTime(t *testing.T) {
	tests := []struct {
		name   string
		decErr string
		encErr string
	}{
		{
			name:   `"2006-01-02 15:04:05 EST"`,
			decErr: "",
			encErr: "",
		},
		{
			name:   `"2006-01-02 12:04:05 UTC"`,
			decErr: "",
			encErr: "",
		},
		{
			name:   `"2006-01-02T15:04:05-07:00"`,
			decErr: `parsing time "2006-01-02T15:04:05-07:00" as "2006-01-02 15:04:05 MST": cannot parse "T15:04:05-07:00" as " "`,
			encErr: "",
		},
		{
			name:   `""`,
			decErr: "",
			encErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var v Time

			err := json.Unmarshal([]byte(tt.name), &v)
			checkErr(t, err, tt.decErr)
			if tt.decErr != "" {
				return
			}

			bb, err := json.Marshal(v)
			checkErr(t, err, tt.encErr)
			if tt.encErr != "" {
				return
			}

			if string(bb) != tt.name {
				t.Errorf("got = %v, want %v", string(bb), tt.name)
			}
		})
	}
}

// TestStringBool tests JSON encoding/parsing functions for the bool values
func TestStringBool(t *testing.T) {
	tests := []struct {
		name   string
		want   string
		decErr string
		encErr string
	}{
		{
			name:   `"true"`,
			want:   `"true"`,
			decErr: "",
			encErr: "",
		},
		{
			name:   `"false"`,
			want:   `"false"`,
			decErr: "",
			encErr: "",
		},
		{
			name:   `"1"`,
			want:   `"true"`,
			decErr: "",
			encErr: "",
		},
		{
			name:   `"2"`,
			want:   `"false"`,
			decErr: "",
			encErr: "",
		},
		{
			name:   `""`,
			want:   `"false"`,
			decErr: "",
			encErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var v StringBool

			err := json.Unmarshal([]byte(tt.name), &v)
			checkErr(t, err, tt.decErr)
			if tt.decErr != "" {
				return
			}

			bb, err := json.Marshal(v)
			checkErr(t, err, tt.encErr)
			if tt.encErr != "" {
				return
			}

			if string(bb) != tt.want {
				t.Errorf("got = %v, want %v", string(bb), tt.name)
			}
		})
	}
}

// checkErr checks for an error
func checkErr(t *testing.T, err error, want string) {
	if (err != nil || want != "") && (err == nil || err.Error() != want) {
		t.Errorf("error = %v, wantErr %v", err, want)
	}
}
