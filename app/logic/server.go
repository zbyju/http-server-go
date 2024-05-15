package logic

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
)

type Server struct {
	Directory string
}

func (server Server) Listen(host string) {
	l, err := net.Listen("tcp", host)
	if err != nil {
		fmt.Println("Failed to bind")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go server.handleConnection(conn)
	}
}

func (server Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	request, err := ParseRequest(conn)
	if err != nil {
		response := HTTPResponse{Version: "1.1", StatusCode: 500}
		response.Send(conn)
	}

	if request.Path == "/" {
		response := HTTPResponse{Version: "1.1", StatusCode: 200}
		response.Send(conn)
	} else if isPrefixed, ending := request.PathStarts("/echo/"); isPrefixed {
		headers := HTTPHeaders{"Content-Type": "text/plain"}
		if request.Headers.GetHeader("Accept-Encoding") == "gzip" {
			headers.SetHeader("Content-Encoding", "gzip")
		}

		response := HTTPResponse{Version: "1.1", StatusCode: 200, Headers: headers, Body: ending}
		response.Send(conn)
	} else if isPrefixed, ending := request.PathStarts("/files/"); request.Method == "POST" && isPrefixed {
		filePath := filepath.Join(server.Directory, ending)

		var response HTTPResponse
		content := []byte(request.Body)

		err := os.WriteFile(filePath, content, 0644)
		if err != nil {
			response = HTTPResponse{Version: "1.1", StatusCode: 500}
		} else {
			response = HTTPResponse{Version: "1.1", StatusCode: 201}
		}
		response.Send(conn)
	} else if isPrefixed, ending := request.PathStarts("/files/"); isPrefixed {
		filePath := filepath.Join(server.Directory, ending)

		var response HTTPResponse
		content, err := os.ReadFile(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				response = HTTPResponse{Version: "1.1", StatusCode: 404}
			} else {
				response = HTTPResponse{Version: "1.1", StatusCode: 500}
			}
		} else {
			response = HTTPResponse{Version: "1.1", StatusCode: 200, Headers: HTTPHeaders{"Content-Type": "application/octet-stream"}, Body: string(content)}
		}
		response.Send(conn)
	} else if isPrefixed, _ := request.PathStarts("/user-agent"); isPrefixed {
		userAgent := request.Headers.GetHeader("User-Agent")
		response := HTTPResponse{Version: "1.1", StatusCode: 200, Headers: HTTPHeaders{"Content-Type": "text/plain"}, Body: userAgent}
		response.Send(conn)
	} else {
		response := HTTPResponse{Version: "1.1", StatusCode: 404}
		response.Send(conn)
	}
}
