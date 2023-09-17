package main

import (
	"cs425-mp/internals/failureDetector"
	"fmt"
	"time"
)

func main() {
	fmt.Println("Failure Detector Started")
	failureDetector.InitializeNodeList()
	go failureDetector.HandleRequests()
	go failureDetector.StartGossipDetector()

	for {
		failureDetector.SendGossip(failureDetector.Gossip)
		time.Sleep(failureDetector.GOSSIP_RATE)
	}

}
