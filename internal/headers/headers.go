package headers

import (
	"bytes"
	"fmt"
)

var ErrInvalidFormat = fmt.Errorf("invalid header format")

type Headers map[string]string

func NewHeaders() Headers {
	return map[string]string{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {

	eol := bytes.Index(data, []byte("\r\n"))
	if eol == -1 {
		return 0, false, nil
	}

	if eol == 0 {
		return 0, true, nil
	}

	line := data[:eol]
	i := bytes.Index(line, []byte(":"))
	if i == -1 {
		return 0, false, ErrInvalidFormat
	}

	if bytes.Contains(line[i-1:i+1], []byte(" ")) {
		return 0, false, ErrInvalidFormat

	}

	key := bytes.Trim(line[:i], " ")
	value := bytes.Trim(line[i+1:], " ")

	h[string(key)] = string(value)

	return eol + 2, false, nil
}
