package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"net"
)

func main() {

	l, err := net.Listen("tcp", ":42069")
	if err != nil {
		panic(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}
		req, err := request.RequestFromReader(conn)
		if err != nil {
			panic(err)
		}

		fmt.Println("Request line:")
		fmt.Println("- Method: ", req.RequestLine.Method)
		fmt.Println("- Target: ", req.RequestLine.RequestTarget)
		fmt.Println("- Version: ", req.RequestLine.HttpVersion)

	}
}
