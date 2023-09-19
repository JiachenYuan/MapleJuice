package failureDetector

import (
	"bufio"
	pb "cs425-mp/protobuf"
	"fmt"
	"os"
	"strings"
)

func HandleUserInput() {
	inputChan := make(chan string)
	//start another goroutine to keep reading user input in a loop
	go GetUserInputInLoop(inputChan)
	//process and read user input. The two goroutines will communicate via a channel
	ProcessUserInputInLoop(inputChan)
}

func GetUserInputInLoop(inputChan chan<- string) {
	for {
		reader := bufio.NewReader(os.Stdin)

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("An error occurred while reading input. Please try again", err)
			return
		}
		//send the user input to inputChan
		inputChan <- input
	}
}

func ProcessUserInputInLoop(inputChan <-chan string) {
	for {
		//read the user input from inputChan
		query := <-inputChan
		if strings.TrimRight(query, "\n") == "leave" {
			NodeInfoList[LOCAL_NODE_KEY].Status = Failed
			leaveMessage := newMessageOfType(pb.GroupMessage_LEAVE)
			SendGossip(leaveMessage)
		}
	}
}
