package main

import (
	"cs425-mp/internals/SDFS"
	"cs425-mp/internals/failureDetector"
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	// go startFailureDetector(&wg)
	startSDFS(&wg)
	wg.Wait()

}

func startFailureDetector(wg *sync.WaitGroup) {
	fmt.Println("Failure Detector Started")
	// Enable logging
	failureDetector.EnableLog()
	// Join the group and initialize local states
	err := failureDetector.JoinGroupAndInit()
	if err != nil {
		fmt.Printf("Cannot join the group: %v\n", err.Error())
		return
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		// Handle http requests for suspicion toggle and drop rate adjustment
		failureDetector.HandleExternalSignals()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		// Handle group messages
		failureDetector.HandleGroupMessages()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		// Periodic local state and memebrship refresh
		failureDetector.PeriodicUpdate()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		// Handle CLI input from user
		failureDetector.HandleUserInput()
	}()
}

func startSDFS(wg *sync.WaitGroup) {
	wg.Add(1)
	fmt.Println("SDFS Started")
	go func() {
		defer wg.Done()
		SDFS.HandleUserInput()
	}()
}