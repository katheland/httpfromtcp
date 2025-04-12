package main

import (
	"fmt"
	"os"
	"log"
	"strings"
	"io"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	ch := getLinesChannel(file)
	for true {
		n, ok := <-ch
		if !ok {
			return
		}
		fmt.Printf("read: %s\n", n)
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