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

func NewInternalServerError() *HandlerError {
	return &HandlerError{
		StatusCode: response.InternalServerError,
		Message:    "unexpected error occured",
	}
}

func NewBadReqError(message string) *HandlerError {
	return &HandlerError{
		StatusCode: response.BadRequest,
		Message:    message,
	}
}

func NewNotFoundError(message string) *HandlerError {
	return &HandlerError{
		StatusCode: 400,
		Message:    message,
	}
}

func (h *HandlerError) Write(w io.Writer) {
	response.WriteStatusLine(w, h.StatusCode)
	response.WriteHeaders(w, response.GetDefaultHeaders(len(h.Message)))
	fmt.Fprintf(w, "%s", h.Message)
}

type Handler func(w *response.Writer, req *request.Request)
