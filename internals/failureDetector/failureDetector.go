package failureDetector

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

const (
	GOSSIP_RATE         = 1000 * time.Millisecond // 1000ms
	T_FAIL              = 1                       // 1 second
	T_CLEANUP           = 1                       // 1 second
	NUM_NODES_TO_GOSSIP = 3                       //number of nodes to gossip to
	PORT                = "55556"
	HOST                = "0.0.0.0"
)

type MessageType int
type StatusType int

const (
	Alive StatusType = iota
	Failed
	Suspected
)

const (
	Join MessageType = iota
	Leave
	Gossip
)

var (
	nodeList           = make(map[string]*Node)
	INTRODUCER_ADDRESS = "fa23-cs425-1801.cs.illinois.edu"
	SERVER_ADDRS       = []string{
		"fa23-cs425-1801.cs.illinois.edu", "fa23-cs425-1802.cs.illinois.edu",
		"fa23-cs425-1803.cs.illinois.edu", "fa23-cs425-1804.cs.illinois.edu",
		"fa23-cs425-1805.cs.illinois.edu", "fa23-cs425-1806.cs.illinois.edu",
		"fa23-cs425-1807.cs.illinois.edu", "fa23-cs425-1808.cs.illinois.edu",
		"fa23-cs425-1809.cs.illinois.edu", "fa23-cs425-1810.cs.illinois.edu"}
)

type Node struct {
	Id               string     `json:"id"`
	Address          string     `json:"Address"`
	HeartbeatCounter int        `json:"heartbeatCounter"`
	Status           StatusType `json:"status"`
	TimeStamp        int        `json:"timeStamp"`
}

func InitializeNodeList() {
	key := getLocalNodeName() + fmt.Sprint(time.Now().Unix())
	initialNodeList := map[string]*Node{
		key: {
			Id:               key,
			Address:          getLocalNodeName(),
			HeartbeatCounter: 1,
			Status:           Alive,
			TimeStamp:        int(time.Now().Unix()),
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
	go startPerodicFailureCheck()
	startListeningToGossips()
}

// start periodic failure check
func startPerodicFailureCheck() {
	for key, node := range nodeList {
		switch node.Status {
		case Alive:
			if time.Now().Unix()-int64(node.TimeStamp) > T_FAIL {
				node.Status = Failed
			}
		case Failed:
			if time.Now().Unix()-int64(node.TimeStamp) > T_CLEANUP {
				delete(nodeList, key)
			}
		case Suspected:
			fmt.Println("Not implemented")
		}

	}
}

func startListeningToGossips() {
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
		if val, ok := nodeList[key]; ok {
			if val.HeartbeatCounter < receivedNode.HeartbeatCounter {
				val.HeartbeatCounter = receivedNode.HeartbeatCounter
				val.TimeStamp = int(time.Now().Unix())
			}
			val.Status = receivedNode.Status
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
	conn, err := net.Dial("udp", INTRODUCER_ADDRESS+":"+PORT)
	if err != nil {
		fmt.Println("Error dialing UDP to introducer: ", err)
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
			defer conn.Close()
			_, err = conn.Write(parsedNodes)
			if err != nil {
				fmt.Println("Error sending UDP: ", err)
				return
			}
		}(node.Address)
	}
	wg.Wait()
}
func SendLeaveMessage() {
	fmt.Println("TODO: SendLeaveMessage not implemented")
}
