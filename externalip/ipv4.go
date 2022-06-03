package externalip

import (
	"errors"
	"io"
	"log"
	"net/url"
)

func GetIPv4(timeout int) (string, error) {
	c := newExternalIPClient(timeout)
	resp, err := c.Get("https://1.1.1.1/cdn-cgi/trace")
	if err != nil {
		if urlErr, ok := err.(*url.Error); ok && urlErr.Timeout() {
			log.Println("1.1.1.1 timed out, falling back to 1.0.0.1")
			resp, err = c.Get("https://1.0.0.1/cdn-cgi/trace")
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
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
