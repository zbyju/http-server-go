package logic

import (
	"net"
	"strings"
)

func readHTTPRequestString(conn net.Conn) (string, error) {
	buf := make([]byte, 2048)
	_, err := conn.Read(buf)
	if err != nil {
		return "", err
	}
	str := string(buf[:])
	endIndex := strings.Index(str, "\x00")
	str = str[0:endIndex]
	return str, nil
}

func parseHTTPRequestString(request string) (string, string, string) {
	split := strings.Split(request, "\r\n\r\n")
	body := split[1]

	lines := strings.Split(split[0], "\r\n")
	requestLine := lines[0]
	requestLine = strings.ReplaceAll(requestLine, "  ", " ")
	headers := strings.Join(lines[1:], "\r\n")

	return requestLine, headers, body
}

func parseHTTPRequestPath(requestLine string) string {
	splitRequestLine := strings.Split(requestLine, " ")
	return splitRequestLine[1]
}

func parseHTTPRequestMethod(requestLine string) string {
	splitRequestLine := strings.Split(requestLine, " ")
	return splitRequestLine[0]
}

func ParseRequest(conn net.Conn) (HTTPRequest, error) {
	requestString, err := readHTTPRequestString(conn)
	if err != nil {
		return HTTPRequest{}, err
	}
	requestLine, headerString, body := parseHTTPRequestString(requestString)
	method := parseHTTPRequestMethod(requestLine)
	path := parseHTTPRequestPath(requestLine)
	headers := ParseHTTPHeaders(headerString)

	return HTTPRequest{Method: method, Path: path, Headers: headers, Body: body}, nil
}

type HTTPRequest struct {
	Method  string
	Path    string
	Headers HTTPHeaders
	Body    string
}

func (request HTTPRequest) PathStarts(str string) (bool, string) {
	if strings.HasPrefix(request.Path, str) {
		rest := strings.TrimPrefix(request.Path, str)
		return true, rest
	}
	return false, ""
}
