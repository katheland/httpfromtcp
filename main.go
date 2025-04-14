package main

import (
	"fmt"
	"log"
	"strings"
	"io"
	"net"
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
		fmt.Println("Connection established")
		
		ch := getLinesChannel(conn)
		for true {
			n, ok := <-ch
			if !ok {
				break
			}
			fmt.Printf("%s\n", n)
		}
		fmt.Println("Connection closed")
		conn.Close()
	}

}

// reads a file 8 bytes at a time, retuning each line via a channel
func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	var b []byte
	currentLine := ""

	go func() {
		for true {
			b = make([]byte, 8)
			_, err := f.Read(b)
			if err != nil {
				break
			}
			s := strings.Split(string(b), "\n")
			currentLine = currentLine + s[0]
			if len(s) == 2 {
				ch <- currentLine
				currentLine = s[1]
			}
		}
		ch <- currentLine
		close(ch)
	}()

	return ch
}