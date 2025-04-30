package server

import (
	"sync/atomic"
	"net"
	//"fmt"
	"log"
	"strconv"
)

type Server struct {
	IsOpen atomic.Bool
	Listener net.Listener
}

func Serve(port int) (*Server, error) {
	p := ":" + strconv.Itoa(port)
	l, err := net.Listen("tcp", p)
	if err != nil {
		log.Fatal(err)
	}
	s := Server{Listener: l}
	s.IsOpen.Store(true)

	go s.listen()

	return &s, nil
}

func (s *Server) Close() error {
	err := s.Listener.Close()
	if err != nil {
		log.Fatal(err)
	}
	s.IsOpen.Store(false)
	return nil
}

func (s *Server) listen() {
	for s.IsOpen.Load() == true {
		conn, err := s.Listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go func(c net.Conn) {
			s.handle(c)
		}(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	//fmt.Println("HTTP/1.1 200 OK")
	//fmt.Println("Content-Type: text/plain")
	//fmt.Println("\nHello World!")

	conn.Write([]byte("HTTP/1.1 200 OK\r\n"))
	conn.Write([]byte("Content-Length: 0\r\n"))
	conn.Write([]byte("Connection: close\r\n"))
	conn.Write([]byte("Content-Type: text/plain\r\n\r\n"))
	conn.Write([]byte("Hello World!"))

	conn.Close()
}