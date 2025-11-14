package main

import (
	"crypto/sha256"
	"fmt"
	"httpfromtcp/internal/headers"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
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

		w.WriteStatusLine(http.StatusOK)

		h := headers.NewHeaders()
		h.Set("Content-Type", "text/plain")
		h.Set("Transfer-Encoding", "chunked")
		h.Set("Trailer", "X-Content-SHA256, X-Content-Length")
		w.WriteHeaders(h)

		body := []byte{}

		for {
			buffer := make([]byte, 32)
			n, err := res.Body.Read(buffer)
			if err != nil {
				break
			}

			body = append(body, buffer[:n]...)
			w.WriteChunkedBody(buffer[:n])
		}

		w.WriteChunkedBodyDone()

		sha := sha256.Sum256(body)
		shaStr := ""
		for _, s := range sha {
			shaStr += fmt.Sprintf("%x", s)
		}

		t := headers.NewHeaders()
		t.Set("X-Content-SHA256", shaStr)
		t.Set("X-Content-Length", fmt.Sprintf("%d", len(body)))
		w.WriteTrailers(t)
		return
	}
}
