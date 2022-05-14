package externalip

import (
	"errors"
	"io"
	"net/http"
)

func GetIPv4() (string, error) {
	resp, err := http.Get("https://1.1.1.1/cdn-cgi/trace")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	kvs := parseTraceResponseBody(body)
	if ip, ok := kvs["ip"]; ok {
		return ip, nil
	} else {
		return "", errors.New("no 'ip' field in trace response")
	}
}
