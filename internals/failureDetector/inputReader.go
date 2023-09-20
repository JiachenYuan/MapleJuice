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
		switch strings.TrimRight(query, "\n") {
		case "leave":
			handleLeave()
		case "rejoin":
			handleRejoin()
		default:
			fmt.Println("Error: input command not supported.")
		}

	}
}

func handleLeave() {
	NodeInfoList[LOCAL_NODE_KEY].Status = Failed
	NodeInfoList[LOCAL_NODE_KEY].SeqNo++
	leaveMessage := newMessageOfType(pb.GroupMessage_LEAVE)
	SendGossip(leaveMessage)
	delete(NodeInfoList, LOCAL_NODE_KEY)
}

func handleRejoin() {
	selfNode, ok := NodeInfoList[LOCAL_NODE_KEY]
	if !ok || selfNode.Status != Failed {
		fmt.Println("Error: cannot rejoin when the current node is still in the network")
		return
	}
	updateLocalNodeKey()
	NodeInfoList[LOCAL_NODE_KEY].Status = Alive
	NodeInfoList[LOCAL_NODE_KEY].SeqNo++
	rejoinMesasge := newMessageOfType(pb.GroupMessage_JOIN)
	SendGossip(rejoinMesasge)
}
