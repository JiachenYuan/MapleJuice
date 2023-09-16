package failureDetector

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
)

// helper function to calculate the max of two nodes
func max(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// helper function to parse node list into byte array
func parseNodeList() []byte {
	byteSlice, err := json.Marshal(nodeList)
	if err != nil {
		fmt.Println("Error encoding node list to json: ", err)
		return nil
	}

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
