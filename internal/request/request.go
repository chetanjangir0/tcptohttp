package request

import (
	"chetanhttpserver/internal/request"
	"errors"
	"fmt"
	"io"
	"strings"
)

type ParserState string

const (
	StateInit ParserState = "init"
	StateDone ParserState = "done"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	ParserState ParserState
}

var ERROR_MALFORMED_REQUEST_LINE = fmt.Errorf("malformed request-line")
var ERROR_UNSUPPORTED_HTTP_VERSION = fmt.Errorf("unsupported http version")
var SEPERATOR = "\r\n"

func newRequest() *Request {
	return &Request{
		ParserState: StateInit,
	}
}

// returns the parsed request line and the number of bytes it consumed
func parseRequestLine(b string) (*RequestLine, int, error) {
	idx := strings.Index(b, SEPERATOR)

	if idx == -1 {
		return nil, b, nil
	}
	requestLine := b[:idx]
	restOfMsg := b[idx+len(SEPERATOR):]

	parts := strings.Split(requestLine, " ")
	if len(parts) != 3 {
		return nil, restOfMsg, ERROR_MALFORMED_REQUEST_LINE
	}

	httpParts := strings.Split(parts[2], "/")
	if len(httpParts) != 2 || httpParts[0] != "HTTP" || httpParts[1] != "1.1" {
		return nil, restOfMsg, ERROR_MALFORMED_REQUEST_LINE
	}

	rl := &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   httpParts[1],
	}
	return rl, restOfMsg, nil
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
	
outer:
	for {
		switch r.ParserState{
		case StateInit:
		case StateDone:
			break outer
		}
	}

	return read, nil

}

func (r *Request) done() bool {
	return r.ParserState == StateDone
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()
	// NOTE: buffer could overrun (a header or body that exceeds 1k bytes)
	buf := make([]byte, 1024)

	bufLen := 0
	for !request.done() {
		n, err := reader.Read(buf[bufLen:])
		// TODO: what do to with error
		if err != nil {
			return nil, err
		}

		bufLen += n

		// readN: how many bytes (from starting/idx 0) it successfully consumed
		readN, err := request.parse(buf[:bufLen])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[readN:bufLen])
		bufLen -= readN
	}

	return request, nil
}
