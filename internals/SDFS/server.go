package SDFS

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

const (
	INTRODUCER_ADDRESS = "fa23-cs425-1801.cs.illinois.edu:55557" // Introducer node's receiving address
	PORT               = "55557"
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
