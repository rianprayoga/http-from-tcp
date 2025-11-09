package request

import (
	"bytes"
	"errors"
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"slices"
	"strconv"
)

var ErrWrongFormat = fmt.Errorf("mismatch format in request line")
var ErrParsedAlready = fmt.Errorf("data  already parsed")
var ErrWrongHttpMethod = fmt.Errorf("incorect http method")

const CRLF = "\r\n"

type ParserState int

const (
	requestLineState ParserState = iota
	headersState
	bodyState
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

func (r *Request) parseLines(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.state != doneState {
		n, err := r.parse(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		totalBytesParsed += n
		if n == 0 {
			break
		}
	}
	return totalBytesParsed, nil
}

func (r *Request) parse(data []byte) (int, error) {

	switch r.state {
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
			r.state = bodyState
		}
		return n, nil
	case bodyState:
		v, ok := r.Headers.Get("Content-Length")

		if v == nil {
			r.state = doneState
			return 0, nil
		}

		bodySize, err := strconv.Atoi(*v)
		if err != nil {
			return 0, err
		}

		if !ok || bodySize == 0 {
			r.state = doneState
			return 0, nil
		}

		r.Body = append(r.Body, data...)

		if len(r.Body) > bodySize {
			return 0, fmt.Errorf("reported content-lentgh was %d but body still parsed", bodySize)
		}

		if len(r.Body) == bodySize {
			r.state = doneState
			return len(data), nil
		}

		return len(data), nil

	default:
		return 0, fmt.Errorf("unknown state")

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
				if r.state != doneState {
					return nil, fmt.Errorf("incomplete request")
				}
				break
			}
			return nil, err
		}

		buffLen += n

		if buffLen >= cap(buffer) {
			tmp := make([]byte, len(buffer)*2)
			copy(tmp, buffer[:])
			buffer = tmp
		}

		n, err = r.parseLines(buffer[:buffLen])
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
