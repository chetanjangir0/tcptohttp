package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	file, err := os.Open("./message.txt")
	if err != nil {
		log.Fatal("error", err)
	}

	lines := getLinesChannel(file)
	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}

}

func getLinesChannel(file io.ReadCloser) <-chan string {
	out := make(chan string, 1)

	go func() {
		defer close(out)
		defer file.Close()

		str := ""
		for {
			b := make([]byte, 8)
			n, err := file.Read(b)
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
