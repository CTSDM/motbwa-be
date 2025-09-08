package auth

import (
	"crypto/rand"
	"testing"

	"net/http"
)

func TestGetHeaderValueTokenAPI(t *testing.T) {
	testCases := []struct {
		name       string
		headerName string
		prefix     string
		hasError   bool
	}{
		{
			name:       "happy path",
			headerName: "Auth",
			prefix:     "Bearer",
			hasError:   false,
		},
		{
			name:       "no header",
			headerName: "",
			prefix:     "",
			hasError:   true,
		},
		{
			name:       "",
			headerName: "",
			prefix:     "",
			hasError:   true,
		},
	}

	for _, tc := range testCases {
		header := make(http.Header)
		tokenValue := rand.Text()
		headerValue := tc.prefix + " " + tokenValue
		if tc.headerName != "" {
			header.Set(tc.headerName, headerValue)
		}

		got, err := GetHeaderValueTokenAPI(header, tc.headerName)
		if tc.hasError {
			if err == nil {
				t.Error("expected some error but got nil")
			}
			return
		}
		if got != tokenValue {
			t.Errorf("got %s, want %s", got, tokenValue)
		}
	}
}
