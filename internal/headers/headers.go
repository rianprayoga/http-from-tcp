package headers

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

var ErrInvalidFormat = fmt.Errorf("invalid header format")
var ErrInvalidKey = fmt.Errorf("invalid Key format")

const CRLF = "\r\n"

type Headers map[string]string

func NewHeaders() Headers {
	return map[string]string{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {

	eol := bytes.Index(data, []byte(CRLF))
	if eol == -1 {
		return 0, false, nil
	}

	if eol == 0 {
		return len(CRLF), true, nil
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
	ok, err := regexp.Match("^[a-zA-Z0-9!#$%&'*+.^_`|~-]+$", key)
	if err != nil || !ok {
		return 0, false, ErrInvalidKey
	}

	value := string(bytes.Trim(line[i+1:], " "))

	lowercaseKey := string(bytes.ToLower(key))

	exstValue, ok := h[lowercaseKey]
	if ok {
		h[lowercaseKey] = exstValue + ", " + value
	} else {
		h[lowercaseKey] = value
	}

	return eol + 2, false, nil
}

func (h Headers) Get(key string) (value *string, ok bool) {

	lowercaseKey := strings.ToLower(key)
	exstValue, ok := h[lowercaseKey]
	if ok {
		return &exstValue, ok
	}

	return nil, false
}
