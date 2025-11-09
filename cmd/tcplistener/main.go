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
		fmt.Printf("- %s: %s\n", "Method", req.RequestLine.Method)
		fmt.Printf("- %s: %s\n", "Target", req.RequestLine.RequestTarget)
		fmt.Printf("- %s: %s\n", "Version", req.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for k, v := range req.Headers {
			fmt.Printf("- %s: %s\n", k, v)
		}
		fmt.Println("Body:")
		fmt.Printf("%s", string(req.Body))

	}
}
