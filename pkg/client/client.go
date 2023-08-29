package client

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

func SendRequest(server string, query string) {
	conn, err := net.Dial("tcp", server)
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	conn.Write([]byte(query))
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading response:", err)
			return
		}
		fmt.Print("Result from", server, ":", line)
	}
}
