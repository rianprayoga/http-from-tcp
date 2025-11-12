package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"net/http"
	"strings"
)

func Urproblem(body string, w *response.Writer, req *request.Request) error {

	if req.RequestLine.RequestTarget == "/yourproblem" {
		w.WriteStatusLine(response.BadRequest)
		h := response.GetDefaultHeaders(len(body))
		h.Set("Content-Type", "text/html")
		w.WriteHeaders(h)
		w.WriteBody([]byte(body))
		return nil
	}
	return nil
}

func MyProblem(body string, w *response.Writer, req *request.Request) error {

	if req.RequestLine.RequestTarget == "/myproblem" {
		w.WriteStatusLine(response.InternalServerError)
		h := response.GetDefaultHeaders(len(body))
		h.Set("Content-Type", "text/html")
		w.WriteHeaders(h)
		w.WriteBody([]byte(body))
		return nil
	}
	return nil
}

func HttpBin(w *response.Writer, req *request.Request) error {
	target := req.RequestLine.RequestTarget
	if after, ok := strings.CutPrefix(target, "/httpbin"); ok {

		res, err := http.Get(fmt.Sprintf("httpbin.org%s", after))
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			return err
		}
		// buffer := make([]byte, 1025)

		// for {
		// 	n, err := res.Body.Read(buffer)
		// 	if err != nil {
		// 		if errors.Is(err, io.EOF) {
		// 			// TODO
		// 			break
		// 		}
		// 		return err
		// 	}

		// }

		// log.Println(after)
		// httpbin/stream/100
		// after
		return nil
	}
	return nil
}
