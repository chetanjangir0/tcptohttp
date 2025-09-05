package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers struct {
	headers map[string]string
}

var rn = []byte("\r\n")

func NewHeaders() *Headers {
	return &Headers{
		headers: map[string]string{},
	}
}

func parseHeader(fieldLine []byte) (string, string, error) {

	parts := bytes.SplitN(fieldLine, []byte(":"), 2) // splitN because the fieldvalue may have a `:` in it
	if len(parts) != 2 {
		return "", "", fmt.Errorf("malformed field line")
	}

	name := parts[0]
	value := bytes.TrimSpace(parts[1])

	if bytes.HasSuffix(name, []byte(" ")) {
		return "", "", fmt.Errorf("malformed field name")
	}

	return string(name), string(value), nil

}

func isToken(str []byte) bool {
	for _, ch := range str {
		found := false

		if ch >= 'A' && ch <= 'Z' ||
			ch >= 'a' && ch <= 'z' ||
			ch >= '0' && ch <= '9' {
			found = true
		}
		if strings.IndexByte("!#$%&'*+-.^_`|~", ch) >= 0 {
			found = true
		}

		if !found {
			return false
		}

	}
	return true
}

func (h *Headers) Get(name string) string {
	return h.headers[strings.ToLower(name)]
}

func (h *Headers) Set(name, val string) {
	name = strings.ToLower(name)

	if oldVal, ok := h.headers[name]; ok{
		val = oldVal + "," + val
	}
	h.headers[name] = val
}

func (h *Headers) Parse(data []byte) (int, bool, error) {
	read := 0
	done := false

	for {
		idx := bytes.Index(data[read:], rn)
		if idx == -1 {
			break
		}

		// headers end with empty  line with the \r\n at the start
		if idx == 0 {
			done = true
			read += len(rn)
			break
		}
		name, value, err := parseHeader(data[read : read+idx])
		if err != nil {
			return 0, false, err
		}

		if !isToken([]byte(name)) {
			return 0, false, fmt.Errorf("malformed header name")
		}
		read += idx + len(rn)
		h.Set(name, value)
	}

	return read, done, nil
}
