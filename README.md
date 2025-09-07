# tcptohttp — HTTP from scratch in Go

**Tagline:**  
> Building HTTP/1.1 from raw TCP, no `net/http` attached.

---

## What’s This?

A simple implementation of HTTP/1.1 built directly on top of TCP in Go.  
No external HTTP libraries, just raw sockets and custom parsing for requests, headers, and body.

---

## Features

- Listen and accept connections over TCP.
- Parse HTTP/1.1 request line, headers, and body.
- Basic response handling with headers and body.
- No dependency on Go's `net/http`.

---

## How to Run

    git clone https://github.com/chetanjangir0/tcptohttp.git
    cd tcptohttp
    go run cmd/httpserver/main.go

Visit http://localhost:42069 or test with curl.

---

## Why Build This?

- To understand how HTTP/1.1 works under the hood.
- To practice low-level networking with Go.
- To see how much work libraries actually save you.

---

## How It Works

1. Start a TCP listener.  
2. Read raw bytes from the connection.  
3. Parse the request line (`GET / HTTP/1.1`).  
4. Collect headers until `\r\n\r\n`.  
5. Parse body if `Content-Length` is set.  
6. Build and send a response.

---

## License

MIT (or your choice).
