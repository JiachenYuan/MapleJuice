package main

import (
	"cs425-mp/internals/failureDetector"
	"fmt"
	"os"
	"time"
)

func main() {
	fmt.Println("Failure Detector Started")
	portNum := os.Args[1]

	initialNodeList := map[string]*failureDetector.Node{
		"localhost:8080": &failureDetector.Node{
			Id:               1,
			IpAddress:        "localhost:8080",
			HeartbeatCounter: 0,
			IsAlive:          true,
			TimeStamp:        0,
		},
		"localhost:8081": &failureDetector.Node{
			Id:               2,
			IpAddress:        "localhost:8081",
			HeartbeatCounter: 0,
			IsAlive:          true,
			TimeStamp:        0,
		},
	}

	failureDetector.SetNodeList(initialNodeList)

	go failureDetector.StartGossipDetector(portNum)

	for {
		failureDetector.SendGossip()
		time.Sleep(failureDetector.GOSSIPRATE)
	}

}
