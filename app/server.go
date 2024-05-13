package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/logic"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	request := logic.GetHTTPRequest(conn)
	path := logic.GetHTTPRequestPath(request)

	if path == "/" {
		_, err = conn.Write([]byte(logic.CreateHTTPResponseString("1.1", "200", "OK", "", "")))
		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
			os.Exit(1)
		}
	} else if strings.HasPrefix(path, "/echo/") {
		str := strings.TrimPrefix(path, "/echo/")
		res := logic.CreateHTTPResponse(200, map[string]string{"Content-Type": "text/plain"}, str)
		_, err = conn.Write([]byte(res))
		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
			os.Exit(1)
		}
	} else {
		_, err = conn.Write([]byte(logic.CreateHTTPResponseString("1.1", "404", "Not Found", "", "")))
		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
			os.Exit(1)
		}
	}
}
