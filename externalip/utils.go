package externalip

import (
	"net"
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

// newExternalIPClient generates a new HTTP client with settings
// suitable to be used to check for the current IP
func newExternalIPClient(timeout int) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   time.Duration(timeout) * time.Second,
				KeepAlive: time.Duration(timeout) * time.Second,
			}).DialContext,
			MaxIdleConns:          10,
			IdleConnTimeout:       time.Duration(timeout) * time.Second,
			TLSHandshakeTimeout:   time.Duration(timeout) * time.Second,
			ExpectContinueTimeout: 2 * time.Second,
		},
	}
}
