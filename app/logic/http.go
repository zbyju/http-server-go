package logic

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
)

func CreateHTTPResponse(statusCode int, headers map[string]string, body string) string {
	bodyLen := len(body)
	headers["Content-Length"] = strconv.Itoa(bodyLen)

	statusText := http.StatusText(statusCode)
	statusCodeStr := strconv.Itoa(statusCode)
	return CreateHTTPResponseString("1.1", statusCodeStr, statusText, CreateHTTPResponseHeader(headers), body)
}

func CreateHTTPResponseString(version, statusCode, statusReason, headers, body string) string {
	return fmt.Sprintf("HTTP/%s %s %s\r\n%s\r\n%s", version, statusCode, statusReason, headers, body)
}

func CreateHTTPResponseHeader(headers map[string]string) string {
	h := ""
	for k, v := range headers {
		h += k + ": " + v + "\r\n"
	}
	return h
}

func GetHTTPRequestPath(request string) string {
	lines := strings.Split(request, "\r\n")
	firstLine := strings.ReplaceAll(lines[0], "  ", " ")
	splitRequestLine := strings.Split(firstLine, " ")
	return splitRequestLine[1]
}

func GetHTTPRequest(conn net.Conn) string {
	buf := make([]byte, 512)
	conn.Read(buf)
	return string(buf[:])
}

func ParseHTTPRequestParts(request string) (string, string, string) {
	split := strings.Split(request, "\r\n\r\n")
	body := split[1]

	lines := strings.Split(split[0], "\r\n")
	requestLine := lines[0]
	headers := strings.Join(lines[1:], "\r\n")

	return requestLine, headers, body
}

func ParseHTTPHeaders(headerStr string) map[string]string {
	lines := strings.Split(headerStr, "\r\n")
	headers := make(map[string]string, 0)
	for _, l := range lines {
		split := strings.Split(l, ": ")
		headers[split[0]] = split[1]
	}
	return headers
}
