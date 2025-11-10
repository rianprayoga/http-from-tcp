package main

import (
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func main() {

	w := func(w io.Writer, req *request.Request) *server.HandlerError {
		switch req.RequestLine.RequestTarget {
		case "/yourproblem":
			return &server.HandlerError{
				StatusCode: response.BadRequest,
				Message:    "Your problem is not my problem\n",
			}
		case "/myproblem":
			return &server.HandlerError{
				StatusCode: response.InternalServerError,
				Message:    "Woopsie, my bad\n",
			}
		default:
			m := "All good, frfr\n"
			_, err := w.Write([]byte(m))
			if err != nil {
				return &server.HandlerError{
					StatusCode: response.InternalServerError,
					Message:    "Unhandled error occured\n",
				}
			}
			return nil
		}
	}

	server, err := server.Serve(port, w)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Server gracefully stopped")
}
