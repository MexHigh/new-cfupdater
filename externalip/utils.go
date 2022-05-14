package externalip

import (
	"net/http"
	"strings"
	"time"
)

// parseTraceResponseBody parses the reponse of https://1.1.1.1/cdn-cgi/trace
// (or its IPv6 counterpart) into a map[string]string
func parseTraceResponseBody(body []byte) map[string]string {
	bodyString := string(body)
	bodySlice := strings.Split(bodyString, "\n")
	kvs := make(map[string]string)
	for _, line := range bodySlice {
		if !strings.Contains(line, "=") {
			continue
		}
		splat := strings.Split(line, "=")
		kvs[splat[0]] = splat[1]
	}
	return kvs
}

// newExternalIPClient generates a new HTTP client with a timeout
// suitable to be used to check for the current IP
func newExternalIPClient(timeout int) *http.Client {
	c := http.DefaultClient
	c.Timeout = time.Duration(timeout) * time.Second
	return c
}
