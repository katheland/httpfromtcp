package main

import (
	"fmt"
	"log"
	"net"
	"github.com/katheland/httpfromtcp/internal/request"
)

func main() {
	listen, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
		}
		
		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s", req.RequestLine.Method, req.RequestLine.RequestTarget, req.RequestLine.HttpVersion)
		
		conn.Close()
	}

}