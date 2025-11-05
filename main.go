package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func getLinesChannel(f io.ReadCloser) <-chan string {

	c := make(chan string)
	var byteLine []byte

	go func() {
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
				close(c)
				break
			}
		}
	}()

	return c
}

func main() {

	file, err := os.Open("messages.txt")
	if err != nil {
		panic(err)
	}

	c := getLinesChannel(file)
	for v := range c {
		fmt.Printf("read: %s\n", v)
	}

}
