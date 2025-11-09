package request

import (
	"bytes"
	"errors"
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"slices"
)

var ErrWrongFormat = fmt.Errorf("mismatch format in request line")
var ErrParsedAlready = fmt.Errorf("data  already parsed")
var ErrWrongHttpMethod = fmt.Errorf("mismatch format in request line")
var CRLF = "\r\n"

type ParserState int

const (
	requestLineState ParserState = iota
	headersState
	doneState
)
const bufferSize = 8

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
	state       ParserState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *Request) parse(data []byte) (int, error) {

	switch state := r.state; state {
	case doneState:
		return 0, ErrParsedAlready
	case requestLineState:
		rl, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}
		r.RequestLine = *rl
		r.state = headersState
		return n, nil
	case headersState:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.state = doneState
		}
		return n, nil
	default:
		return 0, nil

	}

}

func RequestFromReader(reader io.Reader) (*Request, error) {

	r := &Request{
		state:   requestLineState,
		Headers: headers.NewHeaders(),
	}

	buffer := make([]byte, bufferSize)
	buffLen := 0

	for r.state != doneState {

		n, err := reader.Read(buffer[buffLen:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				r.parse(buffer[:buffLen])
				r.state = doneState
				continue
			}
			return nil, err
		}

		buffLen += n

		if buffLen >= cap(buffer) {
			tmp := make([]byte, len(buffer)*2)
			copy(tmp, buffer[:])
			buffer = tmp
		}

		n, err = r.parse(buffer[:buffLen])
		if err != nil {
			return nil, err
		}

		if n != 0 {
			tmp := make([]byte, len(buffer))
			copy(tmp, buffer[n:])
			buffer = tmp
			buffLen -= n
		}
	}

	return r, nil
}

func parseRequestLine(b []byte) (*RequestLine, int, error) {

	i := bytes.Index(b, []byte("\r\n"))
	if i == -1 {
		return nil, 0, nil
	}

	rlPart := bytes.Split(b[:i], []byte(" "))
	if len(rlPart) != 3 {
		return nil, 0, ErrWrongFormat
	}

	method := rlPart[0]
	ok := isValidHttpMethod(method)
	if !ok {
		return nil, 0, ErrWrongHttpMethod
	}

	version := bytes.Split(rlPart[2], []byte("/"))
	if len(version) != 2 || string(version[1]) != "1.1" {
		return nil, 0, ErrWrongFormat
	}

	return &RequestLine{
		Method:        string(method),
		HttpVersion:   string(version[1]),
		RequestTarget: string(rlPart[1]),
	}, i + len(CRLF), nil
}

func isValidHttpMethod(b []byte) bool {
	valid := []string{
		"GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH",
	}

	return slices.Contains(valid, string(b))
}
