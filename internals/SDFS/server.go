package SDFS

import (
	"cs425-mp/internals/failureDetector"
	pb "cs425-mp/protobuf"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"
)

const (
	INTRODUCER_ADDRESS = "fa23-cs425-1801.cs.illinois.edu:55557" // Introducer node's receiving address
	PORT               = "55557"
	CONN_TIMEOUT       = 500 * time.Millisecond
)

var (
	SDFS_PATH   string
	fileToVMMap = make(map[string]map[string]Empty)
	VMToFileMap = make(map[string]map[string]Empty)
)

func init() {
	usr, err := user.Current()
	if err != nil {
		fmt.Printf("Error getting user home directory: %v \n", err)
	}
	SDFS_PATH = filepath.Join(usr.HomeDir, "SDFS_Files")
}

func HandleSDFSMessages() {
	conn, err := net.ListenPacket("udp", ":"+PORT)
	if err != nil {
		fmt.Println("Error listening to UDP packets: ", err)
		os.Exit(1)
	}
	defer conn.Close()
	buffer := make([]byte, 4096)
	for {
		n, _, err := conn.ReadFrom(buffer)
		if err != nil {
			fmt.Printf("Error reading: %v\n", err.Error())
			continue
		}
		SDFSMessage := &pb.SDFSMessage{}
		err = proto.Unmarshal(buffer[:n], SDFSMessage)
		if err != nil {
			fmt.Printf("Error unmarshalling SDFS message: %v\n", err.Error())
		}
		switch SDFSMessage.MessageType {
		case pb.SDFSMessageType_UPDATE:
			processUpdateMessage(SDFSMessage)
		case pb.SDFSMessageType_DELETE:
			processDeleteMessage(SDFSMessage)
		}

	}
}

func processDeleteMessage(message *pb.SDFSMessage) {
	fmt.Println("Received Delete Message")
	fileName := message.SdfsFileName
	delete(fileToVMMap, fileName)
	deleteLocalSDFSFile(fileName)
}

func processUpdateMessage(message *pb.SDFSMessage) {
	fmt.Println("Received Update Message")
	fileName := message.SdfsFileName
	replicas := message.Replicas
	if _, exists := fileToVMMap[fileName]; !exists {
		fileToVMMap[fileName] = make(map[string]Empty)
	}
	for _, r := range replicas {
		fileToVMMap[fileName][r] = Empty{}
	}
}

func putFile(localFileName string, sdfsFileName string) {
	if _, err := os.Stat(localFileName); os.IsNotExist(err) {
		fmt.Printf("Local file not exist: %s\n", localFileName)
		return
	}
	var targetReplicas []string
	val, exists := fileToVMMap[sdfsFileName]
	if !exists {
		fileToVMMap[sdfsFileName] = make(map[string]Empty)
		targetReplicas = getDefaultReplicaMachineIDs(hashFileName(localFileName))
		for _, r := range targetReplicas {
			fileToVMMap[sdfsFileName][r] = Empty{}
		}
	} else {
		for k := range val {
			targetReplicas = append(targetReplicas, k)
		}
	}
	fmt.Printf("Put file %s to sdfs %s \n", localFileName, sdfsFileName)
	for _, id := range targetReplicas {
		targetHostName := getScpHostNameFromID(id)
		remotePath := targetHostName + ":" + filepath.Join(SDFS_PATH, sdfsFileName)
		cmd := exec.Command("scp", localFileName, remotePath)
		err := cmd.Start()
		if err != nil {
			fmt.Printf("Failed to start command: %v\n", err)
			return
		}

		err = cmd.Wait()
		if err != nil {
			fmt.Printf("Command finished with error: %v\n", err)
		}
	}
	sendPutFileMessage(sdfsFileName, targetReplicas)
}

func sendPutFileMessage(fileName string, replicas []string) {
	putMessage := &pb.SDFSMessage{
		MessageType:  pb.SDFSMessageType_UPDATE,
		Replicas:     replicas,
		SdfsFileName: fileName,
	}
	messageBytes, err := proto.Marshal(putMessage)
	if err != nil {
		fmt.Printf("Failed to marshal PutMessage: %v\n", err.Error())
	}
	sendMesageToAllHosts(messageBytes)
}

func sendMesageToAllHosts(messageBytes []byte) {
	allHosts := failureDetector.GetAllNodeAddresses()
	var wg sync.WaitGroup
	for _, hostAddr := range allHosts {
		wg.Add(1)
		go func(address string) {
			defer wg.Done()

			conn, err := net.DialTimeout("udp", address, CONN_TIMEOUT)
			if err != nil {
				// fmt.Println("Error dialing UDP: ", err)
				return
			}
			conn.SetWriteDeadline(time.Now().Add(CONN_TIMEOUT))
			defer conn.Close()
			_, err = conn.Write(messageBytes)
			if err != nil {
				fmt.Println("Error sending UDP: ", err)
				return
			}
		}(hostAddr)
	}
	wg.Wait()
}

func getFile(sdfsFileName string, localFileName string) {
	replicas, exists := fileToVMMap[sdfsFileName]
	if !exists {
		fmt.Printf("SDFS file %v does not exist", sdfsFileName)
	}
	var firstReplicaID string
	for key := range replicas {
		firstReplicaID = key
		break
	}
	targetHostName := getScpHostNameFromID(firstReplicaID)
	fmt.Printf("Get file %s from VM %s", sdfsFileName, targetHostName)
	remotePath := targetHostName + ":" + filepath.Join(SDFS_PATH, sdfsFileName)
	cmd := exec.Command("scp", remotePath, localFileName)
	err := cmd.Start()
	if err != nil {
		fmt.Printf("Failed to start command: %v\n", err)
		return
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Printf("Command finished with error: %v\n", err)
	}
}

func deleteFile(sdfsFileName string) {
	delete(fileToVMMap, sdfsFileName)
	deleteLocalSDFSFile(sdfsFileName)
	sendDeleteFileMessage(sdfsFileName)
}

func deleteLocalSDFSFile(sdfsFileName string) {
	files, err := os.ReadDir(SDFS_PATH)
	if err != nil {
		fmt.Printf("Failed to list files in directory %s: %v\n", SDFS_PATH, err)
		return
	}
	for _, file := range files {
		filePath := filepath.Join(SDFS_PATH, file.Name())
		if !file.IsDir() && strings.HasPrefix(file.Name(), sdfsFileName) {
			fmt.Printf("Try to delete file %s.\n", file.Name())
			err := os.Remove(filePath)
			if err != nil {
				fmt.Printf("Failed to delete file %s: %v\n", filePath, err)
			}
		}
	}
}

func sendDeleteFileMessage(fileName string) {
	deleteMessage := &pb.SDFSMessage{
		MessageType:  pb.SDFSMessageType_DELETE,
		SdfsFileName: fileName,
	}
	messageBytes, err := proto.Marshal(deleteMessage)
	if err != nil {
		fmt.Printf("Failed to marshal DeleteMessage: %v\n", err.Error())
	}
	sendMesageToAllHosts(messageBytes)
}

func getAllLocalSDFSFiles() []string {
	files, err := os.ReadDir(SDFS_PATH)
	if err != nil {
		fmt.Printf("Failed to list files in directory %s: %v\n", SDFS_PATH, err)
		return nil
	}
	var fileNames []string
	for _, f := range files {
		fileNames = append(fileNames, f.Name())
	}
	return fileNames
}

func listSDFSFileVMs(sdfsFileName string) []string {
	var VMList []string
	val, exists := fileToVMMap[sdfsFileName]
	if !exists {
		fmt.Println("Error: file not exist")
	} else {
		for k := range val {
			VMList = append(VMList, getFullHostNameFromID(k))
		}
	}
	return VMList
}
