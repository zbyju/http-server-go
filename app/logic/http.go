package logic

import (
	"fmt"
	"net"
	"strings"
)

func CreateHTTPResponseString(version, statusCode, statusReason, headers, body string) string {
	return fmt.Sprintf("HTTP/%s %s %s\r\n%s\r\n%s", version, statusCode, statusReason, headers, body)
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
