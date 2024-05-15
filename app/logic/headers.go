package logic

import "strings"

type HTTPHeaders map[string]string

func ParseHTTPHeaders(headerStr string) map[string]string {
	lines := strings.Split(headerStr, "\r\n")
	headers := make(map[string]string, 0)
	for _, l := range lines {
		split := strings.Split(l, ": ")
		headers[split[0]] = split[1]
	}
	return headers
}

func (headers HTTPHeaders) GetHeader(key string) string {
	return headers[key]
}

func (headers HTTPHeaders) SetHeader(key, value string) {
	headers[key] = value
}

func (headers HTTPHeaders) ToString() string {
	h := ""
	for k, v := range headers {
		h += k + ": " + v + "\r\n"
	}
	return h
}
