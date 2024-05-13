package logic

import "fmt"

func HttpResponseString(version, statusCode, statusReason, headers, body string) string {
	return fmt.Sprintf("HTTP/%s %s %s\r\n%s\r\n%s", version, statusCode, statusReason, headers, body)
}
