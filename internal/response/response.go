package response

import (
	"chetanhttpserver/internal/headers"
	"fmt"
	"io"
	"strconv"
)

type Response struct {
	StatusCode StatusCode
}

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

func getDefaultHeaders(contentLen int) *headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", strconv.Itoa(contentLen))
	h.Set("Connection", "close") 
	h.Set("Content-Type", "text/plain") 
	return h
}

func WriteHeaders(w io.Writer, h *headers.Headers) error {
	b := []byte{}
	h.ForEach(func (n, v string) {
		b = fmt.Appendf(b, "%s: %s\r\n", n, v)
	})
	_, err := w.Write(b)
	return err
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	statusLine := []byte{}
	switch statusCode {
	case StatusOK:
		statusLine = []byte("HTTP/1.1 200 OK")
	case StatusBadRequest:
		statusLine = []byte("HTTP/1.1 400 Bad Request")
	case StatusInternalServerError:
		statusLine = []byte("HTTP/1.1 500 Internal Server Error")
	default:
		return fmt.Errorf("unrecognized error code")
	}

	_, err := w.Write(statusLine)
	return err 
}
