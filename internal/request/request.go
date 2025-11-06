package request

import (
	"bytes"
	"errors"
	"io"
	"slices"
)

var ErrWrongFormat = errors.New("mismatch format in request line")
var ErrParsedAlready = errors.New("data  already parsed")
var ErrWrongHttpMethod = errors.New("mismatch format in request line")

type ParserState int

const (
	initialized ParserState = iota
	done
)
const bufferSize = 8

var stateName = map[ParserState]string{
	initialized: "initialized",
	done:        "done",
}

type Request struct {
	RequestLine RequestLine
	Headers     map[string]string
	Body        []byte
	state       ParserState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *Request) parse(data []byte) (int, error) {

	if r.state == done {
		return 0, ErrParsedAlready
	}

	rl, n, err := parseRequestLine(data)
	if err != nil {
		return 0, err
	}

	if n == 0 {
		return 0, nil
	}
	r.RequestLine = *rl
	r.state = done
	return n, nil

}

func RequestFromReader(reader io.Reader) (*Request, error) {

	r := &Request{
		state: initialized,
	}

	buf := make([]byte, 0, bufferSize)

	for {
		n, err := reader.Read(buf[len(buf):cap(buf)])
		buf = buf[:len(buf)+n]
		if len(buf) == cap(buf) {
			buf = append(buf, 0)[:len(buf)]
		}

		if err == io.EOF {
			r.state = done
			break
		}

		_, err = r.parse(buf)
		if err != nil {
			return nil, err
		}

		if r.state == done {
			break
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

	// check http method
	method := rlPart[0]
	ok := isValidHttpMethod(method)
	if !ok {
		return nil, 0, ErrWrongHttpMethod
	}

	// check http version
	version := bytes.Split(rlPart[2], []byte("/"))
	if len(version) != 2 || string(version[1]) != "1.1" {
		return nil, 0, ErrWrongFormat
	}

	return &RequestLine{
		Method:        string(method),
		HttpVersion:   string(version[1]),
		RequestTarget: string(rlPart[1]),
	}, i, nil
}

func isValidHttpMethod(b []byte) bool {
	valid := []string{
		"GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH",
	}

	return slices.Contains(valid, string(b))
}
