package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func main() {

	template := "<html><head><title>%s</title>" +
		"</head><body><h1>%s</h1><p>%s</p>" +
		"</body></html>"

	nw := func(w *response.Writer, req *request.Request) {

		Urproblem(fmt.Sprintf(template, "400 Bad Request", "Bad Request", "Your request honestly kinda sucked."), w, req)
		MyProblem(fmt.Sprintf(template, "500 Internal Server Error", "Internal Server Error", "Okay, you know what? This one is on me."), w, req)
		HttpBin(w, req)

		body := fmt.Sprintf(template, "200 OK", "Success!", "Your request was an absolute banger.")
		w.WriteStatusLine(response.Ok)
		h := response.GetDefaultHeaders(len(body))
		h.Set("Content-Type", "text/html")
		w.WriteHeaders(h)
		w.WriteBody([]byte(body))

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
