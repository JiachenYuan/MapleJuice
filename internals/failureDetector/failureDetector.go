package failureDetector

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"
)

const (
	GOSSIPRATE = 500 * time.Millisecond // 1000ms
	B          = 3                      //number of nodes to gossip to
	PORT       = "55556"
	HOST       = "localhost"
)

var (
	localTimer = 0
	nodeList   = make(map[string]*Node)
)

type Node struct {
	Id               int    `json:"id"`
	IpAddress        string `json:"ipAddress"`
	HeartbeatCounter int    `json:"heartbeatCounter"`
	IsAlive          bool   `json:"isAlive"`
	TimeStamp        int    `json:"timeStamp"`
}

func SetNodeList(nodes map[string]*Node) {
	nodeList = nodes
}

// listen to gossip from other nodes
func StartGossipDetector(port string) {
	server, err := net.ListenPacket("udp", ":"+port)
	fmt.Println("Listening on address: ", ":"+port)
	if err != nil {
		fmt.Println("Error listening to UDP packets: ", err)
		os.Exit(1)
	}
	defer server.Close()
	buffer := make([]byte, 1024)
	for {
		//accept connection
		n, _, err := server.ReadFrom(buffer)
		if err != nil {
			fmt.Printf("Error reading: %v\n", err.Error())
			os.Exit(1)
		}

		var receivedMembershipList map[string]*Node
		if err := json.Unmarshal(buffer[:n], &receivedMembershipList); err != nil {
			fmt.Println("Error decoding membership list: ", err)
		} else {
			updateMembershipList(receivedMembershipList)
		}
	}

}

// update current membership list with incoming list
func updateMembershipList(receivedMembershipList map[string]*Node) {
	fmt.Println("Updating membership list")
	for key, receivedNode := range receivedMembershipList {
		fmt.Println("received node is : ", *receivedNode)
		if val, ok := nodeList[key]; ok {
			if val.HeartbeatCounter < receivedNode.HeartbeatCounter {
				val.HeartbeatCounter = receivedNode.HeartbeatCounter
				val.TimeStamp = localTimer
			}
		} else {
			nodeList[key] = receivedNode
		}
	}

}

// send gossip to other nodes
func SendGossip() {
	selectedNodes := randomlySelectNodes(2)
	var wg sync.WaitGroup
	for _, node := range selectedNodes {
		wg.Add(1)
		go func(address string) {
			defer wg.Done()
			parsedNodes := parseNodeList()
			conn, err := net.Dial("udp", address)
			if err != nil {
				fmt.Println("Error dialing UDP: ", err)
				return
			}
			fmt.Println("Sending gossip to: ", address)
			fmt.Println("Gossip: ", string(parsedNodes))
			defer conn.Close()
			_, err = conn.Write(parsedNodes)
			if err != nil {
				fmt.Println("Error sending UDP: ", err)
				return
			}
		}(node.IpAddress)
	}
	wg.Wait()
	fmt.Println("Finished sending gossip")
}

// helper function to parse node list into byte array
func parseNodeList() []byte {
	byteSlice, err := json.Marshal(nodeList)
	if err != nil {
		fmt.Println("Error encoding node list to json: ", err)
		return nil
	}

	// fmt.Println("Encoded node list to json: ", string(byteSlice))
	return byteSlice
}

// helper function to randomly select B nodes to gossip to
func randomlySelectNodes(num int) []*Node {
	num = max(num, len(nodeList))
	keys := make([]string, 0, len(nodeList))
	for k := range nodeList {
		keys = append(keys, k)
	}

	rand.Shuffle(len(keys), func(i, j int) { keys[i], keys[j] = keys[j], keys[i] })
	selectedNodes := make([]*Node, 0, num)
	for i := 0; i < num; i++ {
		selectedNodes = append(selectedNodes, nodeList[keys[i]])
	}
	fmt.Println("Selected nodes: ", selectedNodes)
	return selectedNodes
}
