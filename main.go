package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
)

const (
	envVarDomain  = "GOGETMETA_DOMAIN"
	envVarAddress = "GOGETMETA_ADDRESS"
)

var (
	domain  = "example.org"
	address = ":8080"
)

func main() {
	loadConfig()
	fmt.Printf("gogetmeta\nDomain: %v\nListening on address: %v\n", domain, address)
	http.HandleFunc("/", handleRequest)
	http.ListenAndServe(address, nil)
}

func loadConfig() {
	s := os.Getenv(envVarDomain)
	if s != "" {
		domain = s
	}
	s = os.Getenv(envVarAddress)
	if s != "" {
		address = s
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if !isGoGetRequest(r.URL.Query()) {
		http.Error(w, "Expected request from go-get", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, `<html><head><meta name="go-import" content="%v%v git http://%v%v.git"></head><body>go get meta information</body></html>`, domain, r.URL.Path, domain, r.URL.Path)
}

func isGoGetRequest(query url.Values) bool {
	a := query["go-get"]
	if len(a) == 0 {
		return false
	}
	for _, s := range a {
		if s == "1" {
			return true
		}
	}
	return false
}
