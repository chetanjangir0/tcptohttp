package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("./message.txt")
	if err != nil {
		panic(err)
	}

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
			fmt.Printf("read: %s\n", str)
			str = ""
		}
		str += string(b)

	}

	if len(str) != 0 {
		fmt.Printf("read: %s\n", str)
	}

}
