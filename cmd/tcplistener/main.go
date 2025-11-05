package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
)

func getLinesChannel(f io.ReadCloser) <-chan string {

	c := make(chan string)
	var byteLine []byte

	go func() {

		defer func() {
			close(c)
		}()

		for {
			b := make([]byte, 8)
			n, _ := f.Read(b)

			i := bytes.Index(b, []byte("\n"))

			if i != -1 {
				byteLine = append(byteLine, b[:i]...)
				c <- string(byteLine)

				byteLine = nil
				byteLine = append(byteLine, b[i+1:]...)
				continue
			}

			byteLine = append(byteLine, b[:]...)

			if n == 0 {
				c <- string(byteLine)
				break
			}
		}
	}()

	return c
}

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
		c := getLinesChannel(conn)

		for v := range c {
			fmt.Printf("%s\n", v)
		}
	}
}
