package logic

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
)

type HTTPResponse struct {
	Version    string
	StatusCode int
	Headers    HTTPHeaders
	Body       string
}

func (response HTTPResponse) ToString() string {
	statusCodeString := strconv.Itoa(response.StatusCode)
	statusText := http.StatusText(response.StatusCode)
	if response.Headers == nil {
		response.Headers = HTTPHeaders{}
	}

	response.compress()

	response.Headers["Content-Length"] = strconv.Itoa(len(response.Body))
	headers := response.Headers.ToString()
	return fmt.Sprintf("HTTP/%s %s %s\r\n%s\r\n%s",
		response.Version, statusCodeString, statusText, headers, response.Body)
}

func (response HTTPResponse) Send(conn net.Conn) error {
	_, err := conn.Write([]byte(response.ToString()))
	return err
}

func (response *HTTPResponse) compress() error {
	if response.Headers.GetHeader("Content-Encoding") == "gzip" {
		compressed, err := GzipString(response.Body)
		if err != nil {
			return err
		}
		response.Body = compressed
	}
	return nil
}
