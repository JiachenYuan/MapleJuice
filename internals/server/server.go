package server

import (
	"fmt"
	"net"
	"os"
)

const (
	HOST = "localhost"
	PORT = "55555"
)

func Start() {
	fmt.Println("Server started")
	server, err := net.Listen("tcp", HOST+":"+PORT)
	if err != nil {
		fmt.Printf("Error listening: %v\n", err.Error())
		os.Exit(1)
	}
	defer server.Close()
	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Printf("Error accepting: %v\n", err.Error())
			os.Exit(1)
		}
		go serveConn(conn)
	}
}

func serveConn(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 2048)

	for {
		readLength, err := conn.Read(buffer)
		if err != nil {
			fmt.Printf("Error reading request: %v\n", err.Error())
			return
		}
		
		response := processRequest(string(buffer[:readLength]))
		_, err = conn.Write(response)
		if err != nil {
			fmt.Printf("Error sending response: %v\n", err.Error())
			return
		}
	}
}

func processRequest(request string) []byte {
	// TODO: actual log grep logic and 
	fmt.Println("Received: ", request)
	return []byte("Hello world!\n")
}