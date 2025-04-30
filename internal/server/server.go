package server

import (
	"sync/atomic"
	"net"
	"log"
	"strconv"
	"github.com/katheland/httpfromtcp/internal/response"
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
	response.WriteStatusLine(conn, 200)
	response.WriteHeaders(conn, response.GetDefaultHeaders(0))

	conn.Close()
}