package SDFS

import (
	"crypto/md5"
	"cs425-mp/internals/global"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
	fd "cs425-mp/internals/failureDetector"
)

type Empty struct{}

func hashFileName(fileName string) string {
	hash := md5.Sum([]byte(fileName))
	hashString := fmt.Sprintf("%v", hash)
	total := 0
	for i := 0; i < 10; i++ {
		total += int(hashString[i])
	}
	// add one because our VM ids start from 1
	return fmt.Sprintf("%v", total%10+1)
}

func getDefaultReplicaVMAddresses(id string) []string {
	//TODO: check for two conditions:
	//1. membership list size smaller than 4
	//2. node not exist in membership list in the for loop
	replicas := make([]string, 4)
	val, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Input id cannot be parsed to int")
	}
	for i := 0; i < 4; i++ {
		replicas[i] = getFullHostNameFromID(fmt.Sprintf("%v", ((val+i)%10 + 1)))
	}
	return replicas
}

func getFullHostNameFromID(id string) string {
	numID, err := strconv.Atoi(id) // Convert string to integer
	if err != nil {
		fmt.Println("getFullHostNameFromID: Invalid input ID")
		return ""
	}
	return fmt.Sprintf("fa23-cs425-18%02d.cs.illinois.edu", numID)
}

func getScpHostNameFromHostName(hostName string) string {
	id := getIDFromFullHostName(hostName)
	return fmt.Sprintf("cs425-%s", id)
}

func getIDFromFullHostName(hostName string) string { // might use it later if we decide to do the decode/encode locally
	components := strings.Split(hostName, "-")
	if len(components) < 3 {
		fmt.Printf("Get ID from host name with invalid name: %v \n", hostName)
		return ""
	}
	idWithSuffix := components[2]
	idOnly := strings.Split(idWithSuffix, ".")[0]

	id := idOnly[2:]
	if len(id) > 1 && strings.HasPrefix(id, "0") {
		return id[1:]
	}
	return id
}

func getAllLocalSDFSFilesForVM(vmAddress string) []string {
	var fileNames []string
	files, exists := memTable.VMToFileMap[vmAddress]
	if !exists {
		return fileNames
	}
	for fileName := range files {
		fileNames = append(fileNames, fileName)
	}
	return fileNames
}

func listSDFSFileVMs(sdfsFileName string) []string {
	var VMList []string
	val, exists := memTable.fileToVMMap[sdfsFileName]
	if !exists {
		fmt.Println("Error: file not exist")
	} else {
		for k := range val {
			VMList = append(VMList, k)
		}
	}
	return VMList
}

// Get local server ID based on hostname. Each machines should have a unique ID defined in global.go
// Returns -1 if not defined
func getLocalServerID() int {

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println("Error getting host name: ", err)
		return -1;
	}
	for i, addr := range global.SERVER_ADDRS {
		if hostname == addr {
			return i+1;
		}
	}
	return -1;
}

// Returns server's hostname given its ID
// Returns "" empty string if ID is not defined
func getServerName(id int) string {
	if (id < 1 || id > len(global.SERVER_ADDRS)) {
		return ""
	}
	return global.SERVER_ADDRS[id-1];
}


// Generate random duration bewteen 3-8 seconds
func randomDuration() time.Duration {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(6) + 3
	return time.Duration(n) * time.Second
}

func getAlivePeersAddrs() []string {
	localServerAddr := getServerName(getLocalServerID())
	fd.NodeListLock.Lock()
	addrList := []string{}
	for _, node := range fd.NodeInfoList {
		// Strip away the trailing ":port" part for every address in the list
		delimiterIndex := strings.Index(node.NodeAddr, ":")
		if delimiterIndex == -1 {
			panic("Something wrong in failure detector's NodeInfoList")
		}
		nodeName := node.NodeAddr[:delimiterIndex]
		if (nodeName != localServerAddr) && (node.Status == fd.Alive || node.Status == fd.Suspected){
			addrList = append(addrList, nodeName)
		}
	}
	fd.NodeListLock.Unlock()

	return addrList
}


