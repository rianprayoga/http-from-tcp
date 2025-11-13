package main

import (
	"errors"
	"fmt"
	"httpfromtcp/internal/headers"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"io"
	"net/http"
	"strings"
)

const template = "<html><head><title>%s</title>" +
	"</head><body><h1>%s</h1><p>%s</p>" +
	"</body></html>"

func Urproblem(w *response.Writer, req *request.Request) {

	body := fmt.Sprintf(
		template,
		"400 Bad Request", "Bad Request",
		"Your request honestly kinda sucked.")
	w.WriteStatusLine(response.BadRequest)
	h := response.GetDefaultHeaders(len(body))
	h.Set("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody([]byte(body))

}

func MyProblem(w *response.Writer, req *request.Request) {

	body := fmt.Sprintf(
		template,
		"500 Internal Server Error", "Internal Server Error",
		"Okay, you know what? This one is on me.")
	w.WriteStatusLine(response.InternalServerError)
	h := response.GetDefaultHeaders(len(body))
	h.Set("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody([]byte(body))
}

func HttpBin(w *response.Writer, req *request.Request) {
	target := req.RequestLine.RequestTarget
	if after, ok := strings.CutPrefix(target, "/httpbin"); ok {

		res, err := http.Get(fmt.Sprintf("https://httpbin.org%s", after))
		if err != nil {
			httpErr := server.NewInternalServerError()
			httpErr.Write(w.IoWriter)
			return
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			httpErr := server.NewInternalServerError()
			httpErr.Write(w.IoWriter)
			return
		}

		buffer := make([]byte, 32)

		w.WriteStatusLine(http.StatusOK)

		h := headers.NewHeaders()
		h.Set("Content-Type", "text/plain")
		h.Set("Transfer-Encoding", "chunked")

		w.WriteHeaders(h)

		for {
			n, err := res.Body.Read(buffer)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				httpErr := server.NewInternalServerError()
				httpErr.Write(w.IoWriter)
				return
			}

			_, err = w.WriteChunkedBody(buffer[:n])
			if err != nil {
				httpErr := server.NewInternalServerError()
				httpErr.Write(w.IoWriter)
				return
			}

		}

		_, err = w.WriteChunkedBodyDone()
		if err != nil {
			httpErr := server.NewInternalServerError()
			httpErr.Write(w.IoWriter)
			return
		}

		return
	}
}
