package request

import (
	"bytes"
	"errors"
	"io"
	"slices"
)

var ErrWrongFormat = errors.New("mismatch format in request line")
var ErrWrongHttpMethod = errors.New("mismatch format in request line")

type Request struct {
	RequestLine RequestLine
	Headers     map[string]string
	Body        []byte
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	b, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	part := bytes.Split(b, []byte("\r\n"))
	rl, err := parseRequestLine(part[0])
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *rl,
	}, nil
}

func parseRequestLine(b []byte) (*RequestLine, error) {
	rlPart := bytes.Split(b, []byte(" "))
	if len(rlPart) != 3 {
		return nil, ErrWrongFormat
	}

	// check http method
	method := rlPart[0]
	ok := isValidHttpMethod(method)
	if !ok {
		return nil, ErrWrongHttpMethod
	}

	// check http version
	version := bytes.Split(rlPart[2], []byte("/"))
	if len(version) != 2 || string(version[1]) != "1.1" {
		return nil, ErrWrongFormat
	}

	return &RequestLine{
		Method:        string(method),
		HttpVersion:   string(version[1]),
		RequestTarget: string(rlPart[1]),
	}, nil
}

func isValidHttpMethod(b []byte) bool {
	valid := []string{
		"GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH",
	}

	return slices.Contains(valid, string(b))
}
