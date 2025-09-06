package auth

import "testing"

func TestHasPassword(t *testing.T) {
	testCases := []struct {
		name         string
		passwordText string
		hasError     bool
	}{
		{
			name:         "happy path",
			passwordText: "password",
			hasError:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := HashPassword(tc.passwordText)
			if err != nil {
				t.Fatalf("unexpected error while hashing the password: %s", err)
			}
			if tc.hasError {
				if err == nil {
					t.Error("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("got unexpected error: %s", err)
				}
			}
			if got == tc.passwordText {
				t.Error("output is the same as the input")
			}
		})
	}
}

func TestCheckPasswordHash(t *testing.T) {
	testCases := []struct {
		name         string
		passwordText string
		hasError     bool
	}{
		{
			name:         "happy path",
			passwordText: "password",
			hasError:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hashPassword, err := HashPassword(tc.passwordText)
			if err != nil {
				t.Fatalf("unexpected error while hashing the password: %s", err)
			}

			err = CheckPasswordHash(hashPassword, tc.passwordText)
			if !tc.hasError {
				return
			}
			if err == nil {
				t.Error("expected error but got nil")
			}

		})
	}
}
