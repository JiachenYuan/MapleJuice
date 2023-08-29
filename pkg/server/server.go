package server

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	HOST = "localhost"
)

var portString = ""

func Start(port int) {
	fmt.Println("Server started")
	portString = strconv.Itoa(port)
	server, err := net.Listen("tcp", HOST+":"+portString)
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

	readLength, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("Error reading request: %v\n", err.Error())
		return
	}

	response := processRequest(string(buffer[:readLength]), conn)
	_, err = conn.Write(response)
	if err != nil {
		fmt.Printf("Error sending response: %v\n", err.Error())
		return
	}
}

func processRequest(request string, conn net.Conn) []byte {
	fmt.Println("Message Received:", request)
	parts := strings.SplitN(request, " ", 2)
	if len(parts) != 2 {
		return []byte("Invalid command.")
	}
	command, argument := parts[0], parts[1]
	out, err := exec.Command(command, argument, portString+".txt").Output() // to be replaced with actual log file name
	if err != nil {
		return []byte("Error executing command.")
	}
	return out
}
