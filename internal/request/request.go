package request

import (
	"io"
	"strings"
	"unicode"
	"fmt"
)

const bufferSize = 8

type Status int

const (
	Initialized Status = iota
	Done
)

type Request struct {
	RequestLine RequestLine
	Status Status
}

type RequestLine struct {
	HttpVersion string
	RequestTarget string
	Method string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	req := Request{Status: Initialized,}
	buf := make([]byte, bufferSize, bufferSize)
	readToIndex := 0

	for req.Status != Done {
		if len(buf) >= cap(buf) {
			neo := make([]byte, cap(buf)*2, cap(buf)*2)
			copy(neo, buf)
			buf = neo
		} 

		l, err := reader.Read(buf[readToIndex:])
		if err == io.EOF {
			req.Status = Done
			break
		}
		if err != nil {
			return nil, err
		}

		readToIndex = readToIndex + l
		n, err := req.parse(buf)
		if err != nil {
			return nil, err
		}
		remSuccessful := make([]byte, cap(buf), cap(buf))
		copy(remSuccessful, buf[n:])
		buf = remSuccessful
		readToIndex = readToIndex - n
	}
	return &req, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	// first we split by \r\n
	splitslines := strings.Split(string(data), "\r\n")
	if len(splitslines) < 2 {
		return nil, 0, nil
	}
	
	// next we parse the request line
	splitReq := strings.Split(splitslines[0], " ")
	if len(splitReq) != 3 {
		return nil, 0, fmt.Errorf("Incorrect length of request line")
	}
	method := splitReq[0]
	requestTarget := splitReq[1]
	httpVersion := splitReq[2]
	if !isUpper(method) {
		return nil, 0, fmt.Errorf("method must be all capital letters")
	}
	if httpVersion != "HTTP/1.1" {
		return nil, 0, fmt.Errorf("only supports HTTP/1.1")
	} 

	return &RequestLine{
		HttpVersion: strings.Split(httpVersion, "/")[1],
		RequestTarget: requestTarget,
		Method: method,
	}, len(data), nil
}

func isUpper(s string) bool {
    for _, r := range s {
        if !unicode.IsUpper(r) && unicode.IsLetter(r) {
            return false
        }
    }
    return true
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.Status {
	case Initialized:
		line, read, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if read == 0 {
			return 0, nil
		}
		r.RequestLine = *line
		r.Status = Done
		return len(data), nil
	case Done:
		return 0, fmt.Errorf("Trying to read data in a done state")
	default:
		return 0, fmt.Errorf("Unknown state")
	}
}