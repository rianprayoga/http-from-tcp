package response

import (
	"fmt"
	"httpfromtcp/internal/headers"

	"io"
)

type Writer struct {
	IoWriter io.Writer
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {

	err := WriteStatusLine(w.IoWriter, statusCode)
	if err != nil {
		return fmt.Errorf("failed to write status line")
	}

	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {

	err := WriteHeaders(w.IoWriter, headers)
	if err != nil {
		return fmt.Errorf("failed to write headers")
	}

	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {

	n, err := w.IoWriter.Write(p)
	if err != nil {
		return 0, fmt.Errorf("failed to write body")
	}

	return n, nil
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {

	hex := fmt.Sprintf("%x", len(p))
	b := fmt.Sprintf("%s%s%s%s", hex, CRLF, p, CRLF)
	n, err := w.IoWriter.Write([]byte(b))
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	b := fmt.Sprintf("%d%s%s", 0, CRLF, CRLF)
	n, err := w.IoWriter.Write(([]byte(b)))
	if err != nil {
		return 0, err
	}
	return n, nil
}
