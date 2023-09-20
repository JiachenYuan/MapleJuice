package failureDetector

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
)

// helper function to randomly select B nodes to gossip to
func randomlySelectNodes(num int) []*Node {
	num = min(num, len(NodeInfoList)-1)
	keys := make([]string, 0, len(NodeInfoList))
	for k := range NodeInfoList {
		// do not add self to target list
		if k == LOCAL_NODE_KEY {
			continue
		}
		keys = append(keys, k)
	}

	rand.Shuffle(len(keys), func(i, j int) { keys[i], keys[j] = keys[j], keys[i] })
	selectedNodes := make([]*Node, 0, len(keys))
	for _, nodeKey := range keys {
		nodeInfo := NodeInfoList[nodeKey]
		if nodeInfo.Status == Failed || nodeInfo.Status == Left || nodeInfo.Status == Suspected {
			continue
		}
		selectedNodes = append(selectedNodes, nodeInfo)
		num--
		if num == 0 {
			break
		}
	}
	return selectedNodes
}


func getLocalNodeAddress() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println("Error getting host name: ", err)
		return "", err
	}
	key := hostname + ":" + PORT
	return key, nil
}

// Given a nodeKey in format of [hostname]:[port]:[timestamp], extract the [hostname]:[port] part as a string
func GetAddrFromNodeKey(nodeKey string) string {
	idSplitted := strings.Split(nodeKey, ":")
	peer_name := idSplitted[0]
	peer_port := idSplitted[1]
	return peer_name + ":" + peer_port
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
