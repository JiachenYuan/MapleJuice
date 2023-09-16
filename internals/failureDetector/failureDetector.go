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
	GOSSIP_RATE         = 5000 * time.Millisecond // 1000ms
	NUM_NODES_TO_GOSSIP = 3                       //number of nodes to gossip to
	PORT                = "55556"
	HOST                = "0.0.0.0"
)

type MessageType int

const (
	Join MessageType = iota
	Leave
	Gossip
)

var (
	localTimer         = 0
	nodeList           = make(map[string]*Node)
	COORDINATOR_ADRESS = "fa23-cs425-1801.cs.illinois.edu"
	SERVER_ADDRS       = []string{
		"fa23-cs425-1801.cs.illinois.edu", "fa23-cs425-1802.cs.illinois.edu",
		"fa23-cs425-1803.cs.illinois.edu", "fa23-cs425-1804.cs.illinois.edu",
		"fa23-cs425-1805.cs.illinois.edu", "fa23-cs425-1806.cs.illinois.edu",
		"fa23-cs425-1807.cs.illinois.edu", "fa23-cs425-1808.cs.illinois.edu",
		"fa23-cs425-1809.cs.illinois.edu", "fa23-cs425-1810.cs.illinois.edu"}
)

type Node struct {
	Id               string `json:"id"`
	Address          string `json:"Address"`
	HeartbeatCounter int    `json:"heartbeatCounter"`
	IsAlive          bool   `json:"isAlive"`
	TimeStamp        int    `json:"timeStamp"`
}

func InitializeNodeList() {
	key := getLocalNodeName()
	initialNodeList := map[string]*Node{
		key: {
			Id:               key,
			Address:          getLocalNodeName(),
			HeartbeatCounter: 1,
			IsAlive:          true,
			TimeStamp:        0,
		},
	}

	SetNodeList(initialNodeList)
}

func SetNodeList(nodes map[string]*Node) {
	nodeList = nodes
}

// listen to gossip from other nodes
func StartGossipDetector() {
	go SendGossip(Join)
	server, err := net.ListenPacket("udp", ":"+PORT)
	fmt.Println("Listening on address: ", ":"+PORT)
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
func SendGossip(msgType MessageType) {
	switch msgType {
	case Join:
		SendJoinMessage()
	case Gossip:
		SendGossipMessage()
	case Leave:
		SendLeaveMessage()
	default:
		fmt.Println("Error: unsupported message type")
		os.Exit(1)
	}
}

func SendJoinMessage() {
	conn, err := net.Dial("udp", COORDINATOR_ADRESS+":"+PORT)
	if err != nil {
		fmt.Println("Error dialing UDP to Coordinator: ", err)
		return
	}
	defer conn.Close()
	fmt.Println(getLocalNodeName(), " is joining the cluster")
	parsedNodes := parseNodeList()
	conn.Write(parsedNodes)
}

func SendGossipMessage() {
	selectedNodes := randomlySelectNodes(NUM_NODES_TO_GOSSIP)
	fmt.Println("Number of selected nodes is:", len(selectedNodes))
	getLocalNodeFromNodeList().HeartbeatCounter++
	parsedNodes := parseNodeList()
	var wg sync.WaitGroup
	for _, node := range selectedNodes {
		wg.Add(1)
		go func(address string) {
			defer wg.Done()
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
		}(node.Address)
	}
	wg.Wait()
	fmt.Println("Finished sending gossip")
}
func SendLeaveMessage() {
	fmt.Println("TODO: SendLeaveMessage not implemented")
}

func getLocalNodeFromNodeList() *Node {
	key := getLocalNodeName()
	return nodeList[key]
}

func getLocalNodeName() string {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println("Error getting host name: ", err)
		os.Exit(1)
	}
	key := hostname + ":" + PORT
	return key
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

// helper function to calculate the max of two nodes
func max(a, b int) int {
	if a < b {
		return a
	}
	return b
}
