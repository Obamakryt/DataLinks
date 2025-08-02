package tests

import (
	"DataLinks/internal/api/rest/handler"
	"fmt"
	"github.com/go-playground/validator/v10"
	"testing"
)

type testURLStruct struct {
	URL string `validate:"url"`
}

type testHTTPStruct struct {
	URL string `validate:"http"`
}

func TestUrlValidator(t *testing.T) {
	v := validator.New()
	handler.UrlValidator(v)

	tests := []struct {
		input   string
		wantErr bool
	}{
		{"google.com", false},
		{"https://google.com", false},
		{"https://doesnotexist.badtld", true},
		{"bad", true},
		{"", true},
	}

	for _, test := range tests {
		s := testURLStruct{URL: test.input}
		err := v.Struct(s)
		fmt.Println(err, test.input)
		if (err != nil) != test.wantErr {
			t.Errorf("input: %q, expected error: %v, got error: %v", test.input, test.wantErr, err)
		}
	}
}

func TestHTTPValidator(t *testing.T) {
	v := validator.New()
	handler.HTTPValidator(v)

	tests := []struct {
		input   string
		wantErr bool
	}{
		{"http://example.com", true},
		{"https://example.com", false},
		{"example.com", false},
	}

	for _, test := range tests {
		s := testHTTPStruct{URL: test.input}
		err := v.Struct(s)
		if (err != nil) != test.wantErr {
			t.Errorf("input: %q, expected error: %v, got error: %v", test.input, test.wantErr, err)
		}
	}
}
