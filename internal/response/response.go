package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
)

type StatusCode int16

const (
	Ok                  StatusCode = 200
	BadRequest          StatusCode = 400
	InternalServerError StatusCode = 500
)

const CRLF = "\r\n"

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.Headers{}

	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")

	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {

	for k, v := range headers {
		_, err := fmt.Fprintf(w, "%s: %s%s", k, v, CRLF)
		if err != nil {
			return err
		}
	}

	_, err := fmt.Fprintf(w, "%s", CRLF)
	if err != nil {
		return err
	}

	return nil
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {

	switch statusCode {
	case Ok:
		fmt.Fprintf(w, "HTTP/1.1 200 Ok%s", CRLF)
	case BadRequest:
		fmt.Fprintf(w, "HTTP/1.1 400 Bad Request%s", CRLF)
	case InternalServerError:
		fmt.Fprintf(w, "HTTP/1.1 500 Internal Server Error%s", CRLF)
	default:
		fmt.Fprintf(w, "HTTP/1.1 %d%s", statusCode, CRLF)
	}

	return nil

}
