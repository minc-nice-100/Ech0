package caphttp

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDecodeJSON_StrictUnknownFieldRejected(t *testing.T) {
	req := httptest.NewRequest("POST", "/x", strings.NewReader(`{"secret":"a","response":"b","x":1}`))
	var out struct {
		Secret   string `json:"secret"`
		Response string `json:"response"`
	}
	if err := decodeJSON(req, &out, decodeOptions{Strict: true}); err == nil {
		t.Fatalf("expected strict decode to reject unknown field")
	}
}

func TestDecodeJSON_NonStrictUnknownFieldAllowed(t *testing.T) {
	req := httptest.NewRequest("POST", "/x", strings.NewReader(`{"token":"a","solutions":[1],"x":1}`))
	var out struct {
		Token     string `json:"token"`
		Solutions []int  `json:"solutions"`
	}
	if err := decodeJSON(req, &out, decodeOptions{Strict: false}); err != nil {
		t.Fatalf("expected non-strict decode success, got %v", err)
	}
}

func TestDecodeJSON_MultiplePayloadRejected(t *testing.T) {
	req := httptest.NewRequest("POST", "/x", strings.NewReader(`{"token":"a","solutions":[1]}{"x":1}`))
	var out struct {
		Token     string `json:"token"`
		Solutions []int  `json:"solutions"`
	}
	if err := decodeJSON(req, &out, decodeOptions{Strict: false}); err == nil {
		t.Fatalf("expected decode to reject multiple payloads")
	}
}
