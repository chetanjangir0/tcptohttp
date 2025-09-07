package server

import (
	"bytes"
	"chetanhttpserver/internal/request"
	"chetanhttpserver/internal/response"
	"fmt"
	"io"
	"net"
)

type Server struct {
	closed bool
	handler Handler
}


type HandlerError struct {
	StatusCode response.StatusCode
	Message string
}

type Handler func(w io.Writer, r *request.Request) *HandlerError


func runConnection(s *Server, conn io.ReadWriteCloser) {
	defer conn.Close() // only closing this conn not the server

	headers := response.GetDefaultHeaders(0)
	r, err := request.RequestFromReader(conn)
	if err != nil {
		response.WriteStatusLine(conn, response.StatusBadRequest) 
		response.WriteHeaders(conn, headers) 
		return
	}

	// we let the handler write into this write buffer
	// (instead of directly into conn) so that 
	// we can calculate more info about the body (eg content-length) for headers
	writer := bytes.NewBuffer([]byte{})
	handleError := s.handler(writer, r)

	var body []byte = nil
	var status response.StatusCode= response.StatusOK

	if handleError != nil {
		body = []byte(handleError.Message)
		status = handleError.StatusCode
	} else {
		body = writer.Bytes()
	}

	headers.Replace("Content-length", fmt.Sprintf("%d", len(body)))
	response.WriteStatusLine(conn, status) 
	response.WriteHeaders(conn, headers) 
	conn.Write(body)
}

func runServer(s *Server, listener net.Listener) {
	for {
		conn, err := listener.Accept()

		if s.closed {
			return
		}

		if err != nil {
			return
		}

		// to handle multiple requests
		go func() {
			runConnection(s, conn)
		}()
	}
}

func Serve(port uint16, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	server := &Server{
		closed: false,
		handler: handler,
	}

	go runServer(server, listener)

	return server, nil
}

func (s *Server) Close() error {
	s.closed = true
	return nil
}
