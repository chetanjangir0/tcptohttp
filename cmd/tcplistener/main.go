package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("error", err)
	}

	fmt.Println("Server is listening on :8080")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("error", err)
		}

		lines := getLinesChannel(conn)
		for line := range lines {
			fmt.Printf("read: %s\n", line)
		}
	}

}

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string, 1)

	go func() {
		defer close(out)
		defer f.Close()

		str := ""
		for {
			b := make([]byte, 8)
			n, err := f.Read(b)
			if err != nil {
				break
			}

			b = b[:n]

			if i := bytes.IndexByte(b, '\n'); i != -1 {
				str += string(b[:i])
				b = b[i+1:]
				out <- str
				str = ""
			}
			str += string(b)

		}

		if len(str) != 0 {
			out <- str
		}
	}()

	return out
}
