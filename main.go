package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"cs425-mp1/pkg/client"
	"cs425-mp1/pkg/server"
)

var servers = []string{"localhost:8080", "localhost:8081"} //needs to be deleted after we get the actual VM addresses
var currentServerName string = ""

func getUserInputInLoop(portString string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println(">>> Enter Query: ")
		os.Stdout.Sync()
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("An error occured while reading input. Please try again", err)
			return
		}

		input = strings.TrimSuffix(input, "\n")
		fmt.Printf("The query pattern you entered is %s\n", input)
		for _, server := range servers {
			if strings.Compare("localhost:"+portString, server) != 0 {
				client.SendRequest(server, input)
			}
		}
	}
}

func main() {
	portString := os.Args[1]
	currentServerName = portString
	port, err := strconv.Atoi(portString)
	if err != nil {
		fmt.Println("Error parsing port:", err)
		return
	}

	go server.Start(port)

	// 2. Start waiting for user input
	getUserInputInLoop(portString)

}
