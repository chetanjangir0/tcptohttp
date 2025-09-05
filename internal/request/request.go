package request

import (
	"bytes"
	"chetanhttpserver/internal/headers"
	"fmt"
	"io"
	"strconv"
)

type ParserState string

const (
	StateInit    ParserState = "init"
	StateDone    ParserState = "done"
	StateBody    ParserState = "body"
	StateHeaders ParserState = "headers"
	StateError   ParserState = "error"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	ParserState ParserState
	Headers     *headers.Headers
	Body        string
}

var ERROR_MALFORMED_REQUEST_LINE = fmt.Errorf("malformed request-line")
var ERROR_UNSUPPORTED_HTTP_VERSION = fmt.Errorf("unsupported http version")
var ERROR_REQUEST_IN_ERROR_STATE = fmt.Errorf("request in error state")
var SEPERATOR = []byte("\r\n")

func getInt(headers *headers.Headers, name string, defaultValue int) int {
	valueStr, exists := headers.Get(name)
	if !exists {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func newRequest() *Request {
	return &Request{
		ParserState: StateInit,
		Headers:     headers.NewHeaders(),
		Body:        "",
	}
}

// returns the parsed request line and the number of bytes it consumed
func parseRequestLine(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, SEPERATOR)

	if idx == -1 {
		return nil, 0, nil
	}
	requestLine := b[:idx]
	read := idx + len(SEPERATOR)

	parts := bytes.Split(requestLine, []byte(" "))
	if len(parts) != 3 {
		return nil, 0, ERROR_MALFORMED_REQUEST_LINE
	}

	httpParts := bytes.Split(parts[2], []byte("/"))
	if len(httpParts) != 2 || string(httpParts[0]) != "HTTP" || string(httpParts[1]) != "1.1" {
		return nil, 0, ERROR_MALFORMED_REQUEST_LINE
	}

	rl := &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   string(httpParts[1]),
	}
	return rl, read, nil
}

func (r *Request) hasBody() bool {
	// TODO: when doing chunked encoding update this method
	length := getInt(r.Headers, "content-length", 0)
	return length > 0
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0

outer:
	for {
		currentData := data[read:]
		if len(currentData) == 0 {
			break outer
		}
		switch r.ParserState {
		case StateError:
			return 0, ERROR_REQUEST_IN_ERROR_STATE
		case StateInit:
			rl, n, err := parseRequestLine(currentData)
			if err != nil {
				r.ParserState = StateError
				return 0, err
			}

			// return what’s been read so far → wait for more input.
			if n == 0 {
				break outer
			}

			r.RequestLine = *rl
			read += n
			r.ParserState = StateHeaders

		case StateHeaders:
			n, done, err := r.Headers.Parse(currentData)
			if err != nil {
				r.ParserState = StateError
				return 0, err
			}

			if n == 0 {
				break outer
			}
			read += n

			// in real world we would not get EOF error after reading data
			// therefore we would nicely transition to body, which would allow use to 
			// transition to done, but i am doing the transtion here
			// so we would not need to check for hasBody
			if done {
				if r.hasBody() {
					r.ParserState = StateBody
				} else {
					r.ParserState = StateDone
				}
			}
		case StateBody:
			length := getInt(r.Headers, "content-length", 0)
			if length == 0 {
				panic("chunked not implemented")
			}

			remaining := min(length-len(r.Body), len(currentData))
			r.Body += string(currentData[:remaining])
			read += remaining

			if len(r.Body) == length {
				r.ParserState = StateDone
			}

		case StateDone:
			break outer
		default:
			panic("skill issue")
		}
	}

	return read, nil

}

func (r *Request) done() bool {
	return r.ParserState == StateDone || r.ParserState == StateError
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()
	// NOTE: buffer could overrun (a header or body that exceeds 1k bytes)
	buf := make([]byte, 1024)

	bufLen := 0
	for !request.done() {

		// Reading into the buffer
		n, err := reader.Read(buf[bufLen:])
		if err != nil { // TODO: what do to with error
			return nil, err
		}
		bufLen += n

		// parsing the buffer
		readN, err := request.parse(buf[:bufLen])
		// readN: how many bytes (from starting/idx 0) it successfully consumed
		if err != nil {
			return nil, err
		}

		// remove the data that is parsed to save memory
		copy(buf, buf[readN:bufLen])
		bufLen -= readN
	}

	return request, nil
}
