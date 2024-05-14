package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/codecrafters-io/http-server-starter-go/app/logic"
)

type DirectoryHolder struct {
	Directory string
}

var instance *DirectoryHolder
var once sync.Once

func GetInstance() *DirectoryHolder {
	once.Do(func() {
		dir := flag.String("directory", "", "the directory to search for text.txt")
		flag.Parse()

		if *dir == "" {
			fmt.Println("Please provide a directory using the --directory flag.")
			os.Exit(1)
		}
		instance = &DirectoryHolder{Directory: *dir}
	})
	return instance
}
func (d *DirectoryHolder) SetDirectory(dir string) {
	d.Directory = dir
}
func (d *DirectoryHolder) GetDirectory() string {
	return d.Directory
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	request := logic.GetHTTPRequest(conn)
	path := logic.GetHTTPRequestPath(request)
	_, headersStr, _ := logic.ParseHTTPRequestParts(request)
	headers := logic.ParseHTTPHeaders(headersStr)

	if path == "/" {
		_, err := conn.Write([]byte(logic.CreateHTTPResponseString("1.1", "200", "OK", "", "")))
		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
			os.Exit(1)
		}
	} else if strings.HasPrefix(path, "/echo/") {
		str := strings.TrimPrefix(path, "/echo/")
		res := logic.CreateHTTPResponse(200, map[string]string{"Content-Type": "text/plain"}, str)
		_, err := conn.Write([]byte(res))
		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
			os.Exit(1)
		}
	} else if strings.HasPrefix(path, "/files/") {
		str := strings.TrimPrefix(path, "/files/")
		dir := GetInstance().Directory
		filePath := filepath.Join(dir, str)

		// Check if the file exists
		var res string
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				res = logic.CreateHTTPResponse(404, map[string]string{}, "")
			} else {
				res = logic.CreateHTTPResponse(500, map[string]string{}, "")
			}
		} else {
			res = logic.CreateHTTPResponse(200, map[string]string{"Content-Type": "application/octet-stream"}, string(content))
		}

		_, err = conn.Write([]byte(res))
		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
			os.Exit(1)
		}
	} else if strings.HasPrefix(path, "/user-agent") {
		res := logic.CreateHTTPResponse(200, map[string]string{"Content-Type": "text/plain"}, headers["User-Agent"])
		_, err := conn.Write([]byte(res))
		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
			os.Exit(1)
		}
	} else {
		_, err := conn.Write([]byte(logic.CreateHTTPResponseString("1.1", "404", "Not Found", "", "")))
		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
			os.Exit(1)
		}
	}
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(conn)
	}
}
