package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8888")

	if err != nil {
		fmt.Println("An error occurred:")
		fmt.Println(err)
		return
	}

	inputReader := bufio.NewReader(os.Stdin)
	reader := bufio.NewReader(conn)

	for {
		fmt.Printf("Enter a message: ")
		input, _ := inputReader.ReadString('\n')
		input = strings.TrimSpace(input)

		fmt.Fprintf(conn, input+"\n")
		line, err := reader.ReadString('\n')

        if err != nil {
            fmt.Println("Error reading: ")
            fmt.Println(err)
            break
        }

        fmt.Println(line)
	}

	conn.Close()
}
