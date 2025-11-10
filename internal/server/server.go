package server

import (
	"bytes"
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"io"
	"net"
	"strconv"
	"sync/atomic"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Server struct {
	listener net.Listener
	closed   atomic.Bool
	handler  Handler
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func Serve(port int, h Handler) (*Server, error) {

	l, err := net.Listen("tcp", fmt.Sprintf(":%s", strconv.Itoa(port)))
	if err != nil {
		return nil, err
	}

	s := &Server{
		listener: l,
		handler:  h,
	}
	s.closed.Store(false)

	go s.listen()

	return s, nil

}

func (s *Server) Close() error {
	err := s.listener.Close()
	if err != nil {
		return err
	}
	s.closed.Store(true)
	return nil
}

func (s *Server) listen() {

	for !s.closed.Load() {
		conn, _ := s.listener.Accept()
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)
	if err != nil {

	}

	var b bytes.Buffer
	httpErr := s.handler(&b, req)
	if httpErr != nil {
		writeErrorResponse(conn, httpErr)
	}

	body := b.String()
	response.WriteStatusLine(conn, response.Ok)
	response.WriteHeaders(conn, response.GetDefaultHeaders(len(body)))
	fmt.Fprintf(conn, "%s", body)

}

func writeErrorResponse(w io.Writer, err *HandlerError) {
	if err != nil {
		response.WriteStatusLine(w, err.StatusCode)
		response.WriteHeaders(w, response.GetDefaultHeaders(len(err.Message)))
		fmt.Fprintf(w, "%s", err.Message)
		return
	}
}
