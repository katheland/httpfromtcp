package response

import (
	"io"
	"github.com/katheland/httpfromtcp/internal/headers"
	"fmt"
	"strconv"
	"log"
)

type StatusCode int

const (
	OK StatusCode = 200
	BadRequest StatusCode = 400
	InternalServer StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	reason := ""
	switch statusCode {
	case OK:
		reason = "OK"
	case BadRequest:
		reason = "Bad Request"
	case InternalServer:
		reason = "Internal Server Error"
	}
	_, err := w.Write([]byte(fmt.Sprintf("HTTP/1.1 %v %s\r\n", statusCode, reason)))
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	head := headers.NewHeaders()
	head["content-length"] = strconv.Itoa(contentLen)
	head["connection"] = "close"
	head["content-type"] = "text/plain"
	return head
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, value := range headers {
		_, err := w.Write([]byte(fmt.Sprintf("%s: %s\r\n", key, value)))
		if err != nil {
			log.Fatal(err)
		}
	}
	w.Write([]byte("\r\n"))
	return nil
}