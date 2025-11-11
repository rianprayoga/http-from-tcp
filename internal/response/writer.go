package response

import (
	"fmt"
	"httpfromtcp/internal/headers"

	"io"
)

type WriterState int

const (
	StatusLineState WriterState = iota
	HeadersState
	BodyState
	Done
)

type Writer struct {
	// State    WriterState
	IoWriter io.Writer
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	// if w.State != StatusLineState {
	// 	return fmt.Errorf("can't write status line")
	// }

	err := WriteStatusLine(w.IoWriter, statusCode)
	if err != nil {
		return fmt.Errorf("failed to write status line")
	}
	// w.State = HeadersState

	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	// if w.State != HeadersState {
	// 	return fmt.Errorf("can't write headers")
	// }

	err := WriteHeaders(w.IoWriter, headers)
	if err != nil {
		return fmt.Errorf("failed to write headers")
	}

	// w.State = BodyState
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	// if w.State != BodyState {
	// 	return 0, fmt.Errorf("can't write body")
	// }

	n, err := w.IoWriter.Write(p)
	if err != nil {
		return 0, fmt.Errorf("failed to write body")
	}

	// w.State = Done
	return n, nil
}
