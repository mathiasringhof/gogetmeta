package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

func TestIsGoGetRequestWorks(t *testing.T) {
	vals := url.Values(map[string][]string{"go-get": []string{"1"}})
	if !isGoGetRequest(vals) {
		t.Fatalf("Failed with query vals: %+v", vals)
	}
}

func TestRequestWithoutGoGetParameterGetsError(t *testing.T) {
	withHandleRequestResult(t, "/something", func(rec *httptest.ResponseRecorder) {
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("Expected request to return status bad request, but got %v", rec.Code)
		}
	})
}

func TestRequestWithGoGetReturnsProperMetaTag(t *testing.T) {
	withHandleRequestResult(t, "/something?go-get=1", func(rec *httptest.ResponseRecorder) {
		if rec.Code != http.StatusOK {
			t.Errorf("Expected request to return status OK, but got %v", rec.Code)
		}
		expected := fmt.Sprintf(`<meta name="go-import" content="%v/something git http://%v/something.git">`, domain, domain)
		if !strings.Contains(rec.Body.String(), expected) {
			t.Fatalf("Response did not contain expected meta tag. Expected: %v\nResponse:\n%v", expected, rec.Body.String())
		}
	})
}

func TestLoadConfigReadsDomainFromEnvironment(t *testing.T) {
	d := "mydomain.org"
	os.Setenv(envVarDomain, d)
	loadConfig()
	if domain != d {
		t.Fatalf("Expected domain to be the same as environment variable %v (value = %v), but got %v", envVarDomain, d, domain)
	}
}

func TestLoadConfigReadsPortFromEnvironment(t *testing.T) {
	a := ":9991"
	os.Setenv(envVarAddress, a)
	loadConfig()
	if address != a {
		t.Fatalf("Expected address to be the same as environment variable %v (value = %v), but got %v", envVarAddress, a, address)
	}
}

func withHandleRequestResult(t *testing.T, URL string, closure func(*httptest.ResponseRecorder)) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}
	handleRequest(rec, req)
	closure(rec)
}
