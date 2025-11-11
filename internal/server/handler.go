package server

import (
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"io"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

func (h *HandlerError) Write(w io.Writer) {
	response.WriteStatusLine(w, h.StatusCode)
	response.WriteHeaders(w, response.GetDefaultHeaders(len(h.Message)))
	fmt.Fprintf(w, "%s", h.Message)
}

type Handler func(w *response.Writer, req *request.Request)
