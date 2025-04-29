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
		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n", req.RequestLine.Method, req.RequestLine.RequestTarget, req.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for key, val := range req.Headers {
			fmt.Println(fmt.Sprintf("- %s: %s", key, val))
		}
		fmt.Println("Body:")
		fmt.Println(string(req.Body))
		
		conn.Close()
	}

}