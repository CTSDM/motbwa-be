package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestParseLoginRequest(t *testing.T) {
	type testCases []struct {
		p   any
		err error
	}

	happyCases := testCases{
		{
			p: parameters{
				Username: "user",
				Password: "test",
			},
			err: nil,
		},
		{
			p: parameters{
				Username: "test",
				Password: "user",
			},
			err: nil,
		},
	}

	t.Run("Return the correct inputs", func(t *testing.T) {
		cfg := CfgAPI{}
		for _, c := range happyCases {
			want := c.p
			sending, err := json.Marshal(want)
			if err != nil {
				t.Fatal("Something went wrong while marshaling")
			}
			body := io.NopCloser(bytes.NewBuffer(sending))
			r := &http.Request{
				Body: body,
			}
			got, err := cfg.parseLoginRequest(r)
			if err != c.err {
				t.Errorf("Got err: %s", err)
			}
			if *got != want {
				t.Errorf("Got %q, want %q", got, want)
			}
		}
	})

	t.Run("Identifies outputs with different values", func(t *testing.T) {
		cfg := CfgAPI{}
		want := parameters{
			Username: "hello",
			Password: "byebye",
		}
		for _, c := range happyCases {
			sending, err := json.Marshal(c.p)
			if err != nil {
				t.Fatal("Something went wrong while marshaling")
			}
			body := io.NopCloser(bytes.NewBuffer(sending))
			r := &http.Request{
				Body: body,
			}
			got, err := cfg.parseLoginRequest(r)
			if err != c.err {
				t.Errorf("Got err: %s", err)
			}
			if *got == want {
				t.Errorf("Got %q, want %q", got, want)
			}
		}
	})

	type wrongParam struct {
		Email string `json:"email"`
	}

	wrongJSONCase := testCases{
		{
			p: wrongParam{
				Email: "a@a.com",
			},
			err: nil,
		},
	}

	t.Run("Identifies different shapes", func(t *testing.T) {
		cfg := CfgAPI{}
		want := parameters{
			Username: "hello",
			Password: "byebye",
		}
		for _, c := range wrongJSONCase {
			sending, err := json.Marshal(c.p)
			if err != nil {
				t.Fatal("Something went wrong while marshaling")
			}
			body := io.NopCloser(bytes.NewBuffer(sending))
			r := &http.Request{
				Body: body,
			}
			got, err := cfg.parseLoginRequest(r)
			if err != c.err {
				t.Errorf("Got err: %s", err)
			}
			if *got == want {
				t.Errorf("Got %q, want %q", got, want)
			}
		}
	})
}
