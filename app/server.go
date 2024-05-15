package main

import (
	"flag"

	"github.com/codecrafters-io/http-server-starter-go/app/logic"
)

func main() {
	dir := flag.String("directory", "", "the directory to search for text.txt")
	flag.Parse()

	server := logic.Server{Directory: *dir}
	server.Listen("0.0.0.0:4221")
}
