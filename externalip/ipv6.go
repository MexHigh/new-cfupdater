package externalip

import (
	"errors"
	"io"
)

func GetIPv6(timeout int) (string, error) {
	c := newExternalIPClient(timeout)
	resp, err := c.Get("https://[2606:4700:4700::1111]/cdn-cgi/trace")
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
