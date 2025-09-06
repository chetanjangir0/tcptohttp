package response

import (
	"fmt"
	"io"
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
