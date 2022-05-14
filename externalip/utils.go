package externalip

import (
	"strings"
)

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
