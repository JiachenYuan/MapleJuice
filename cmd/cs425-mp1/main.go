package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"sync"

	// "cs425-mp1/internals/client"
	"cs425-mp1/internals/server"
)

func getUserInputInLoop(inputChan chan<- string) {
	for {
		fmt.Println(">>> Enter Query: ")
		reader := bufio.NewReader(os.Stdin)
		
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("An error occured while reading input. Please try again", err)
			return
		}

		input = strings.TrimSuffix(input, "\n")
		inputChan <- input
	}
}

func main() {
	fmt.Println("App started")

	var wg sync.WaitGroup

	// 1. Start log server 
	wg.Add(1)
	go func() {
		defer wg.Done()
		server.Start()
	}()

	// 2. Start waiting for user input
	inputChan := make(chan string)
	wg.Add(2)
	go func() {
		defer wg.Done()
		getUserInputInLoop(inputChan)
	}()
	go func() {
		defer wg.Done()
		for {
			queryPattern := <- inputChan
			// todo: hand over the logic of query processing to client package
			fmt.Printf("The query pattern you entered is %s\n", queryPattern)
		}
	}()
	
	// Wait indefinitely until Ctrl-C
	wg.Wait()
}