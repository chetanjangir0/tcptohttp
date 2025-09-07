package server

import (
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

type Handler func(w *response.Writer, r *request.Request) 


func runConnection(s *Server, conn io.ReadWriteCloser) {
	defer conn.Close() // only closing this conn not the server

	responseWriter := response.NewWriter(conn)
	r, err := request.RequestFromReader(conn)
	if err != nil {
		responseWriter.WriteStatusLine(response.StatusBadRequest) 
		return
	}

	s.handler(responseWriter, r)

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
