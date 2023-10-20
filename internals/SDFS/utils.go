package SDFS

import (
	"crypto/md5"
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

func getDefaultReplicaMachineIDs(id string) []string {
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

func getScpHostNameFromID(id string) string {
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
	return id
}
