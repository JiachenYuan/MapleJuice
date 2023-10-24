package SDFS

import (
	"crypto/md5"
	"cs425-mp/internals/failureDetector"
	"fmt"
	"strconv"
	"strings"
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

func isCurrentNodeLeader() bool {
	return HOSTNAME == LEADER_ADDRESS
}

func getDefaultReplicaVMAddresses(id string) []string {
	membershipList := failureDetector.GetAllNodeAddresses()
	replicaSize := max(4, len(membershipList))

	replicas := make([]string, replicaSize)

	if replicaSize < NUM_WRITE {
		replicas = append(replicas, membershipList...)
	} else {
		val, err := strconv.Atoi(id)
		if err != nil {
			fmt.Println("Input id cannot be parsed to int")
		}
		i := 0
		for i < 4 {
			hostName := getFullHostNameFromID(fmt.Sprintf("%v", ((val+i)%10 + 1)))
			if failureDetector.IsNodeAlive(hostName) {
				replicas[i] = hostName
				i++
			}

		}
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

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
