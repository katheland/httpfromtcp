package main

import (
	"net"
	"log"
	"bufio"
	"os"
	"fmt"
)

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	read := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		input, err := read.ReadString([]byte("\n")[0])
		if err != nil {
			log.Fatal(err)
		}
		_, err = conn.Write([]byte(input))
		if err != nil {
			log.Fatal(err)
		}
	}
}