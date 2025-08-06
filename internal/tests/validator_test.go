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
	handler.UrlValidator(v) // Регистрируем кастомную валидацию

	tests := []struct {
		input   string
		wantErr bool
	}{
		{"https://openai.com", false},
		{"https://google.com", false},
		{"https://github.com", false},
		{"http://example.com", false},
		{"example.com", true},
		{"http://", true},
		{"", true},
		{"not a url", true},
	}

	for i, tt := range tests {
		data := testURLStruct{URL: tt.input}
		err := v.Struct(data)
		if (err != nil) != tt.wantErr {
			t.Errorf("Test %d: URL=%q - expected error=%v, got error=%v", i+1, tt.input, tt.wantErr, err != nil)
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
func TestValidate(t *testing.T) {
	v := validator.New()

	type User struct {
		Password string `validate:"required,min=8,max=24"`
	}

	user := User{Password: "Zxc1234"}
	err := v.Struct(user)
	if err != nil {
		fmt.Println("Validation error:", err)
	} else {
		fmt.Println("Validation passed")
	}
}
