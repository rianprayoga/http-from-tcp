package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const port = 42069

func main() {

	nw := func(w *response.Writer, req *request.Request) {
		switch req.RequestLine.RequestTarget {
		case "/yourproblem":
			Urproblem(w, req)
			return
		case "/myproblem":
			MyProblem(w, req)
			return

		case "/video":
			GetVideo(w, req)
			return
		default:
			if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin") {
				HttpBin(w, req)
				return
			}

			body := fmt.Sprintf(template, "200 OK", "Success!", "Your request was an absolute banger.")
			w.WriteStatusLine(response.Ok)
			h := response.GetDefaultHeaders(len(body))
			h.Set("Content-Type", "text/html")
			w.WriteHeaders(h)
			w.WriteBody([]byte(body))
		}

	}

	server, err := server.Serve(port, nw)
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
