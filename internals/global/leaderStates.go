package global

import (
	pb "cs425-mp/protobuf"
	"fmt"
	"strconv"
	"sync"
)

/************************************************* Leader's States ***************************************************************/
type Empty struct{}
type _MemTable struct {
	FileToVMMap map[string]map[string]Empty
	VMToFileMap map[string]map[string]Empty
}

type RequestType int

const (
	READ RequestType = iota
	WRITE
)

type FileLock struct {
	ReadQueue         []string
	WriteQueue        []string
	ReadCount         int
	WriteCount        int
	ConsecutiveReads  int
	ConsecutiveWrites int
	FileLocksMutex    sync.Mutex
}

var (
	FileLocks      = make(map[string]*FileLock)
	GlobalFileLock = sync.Mutex{}
)

var (
	MemTable = &_MemTable{
		FileToVMMap: make(map[string]map[string]Empty), // go does not have sets, so we used a map with empty value to repersent set
		VMToFileMap: make(map[string]map[string]Empty),
	}
	Version  = 0  // Monotonically incremented when leader send Gossip in Failure Detector Module
)

// update mem tables
func (mt *_MemTable) Delete(sdfsFileName string) {
	for _, files := range mt.VMToFileMap {
		delete(files, sdfsFileName)
	}
	delete(mt.FileToVMMap, sdfsFileName)
}

func (mt *_MemTable) Put(sdfsFileName string, replicas []string) {
	if _, exists := mt.FileToVMMap[sdfsFileName]; !exists {
		mt.FileToVMMap[sdfsFileName] = make(map[string]Empty)
	}
	for _, r := range replicas {
		if _, exists := mt.VMToFileMap[r]; !exists {
			mt.VMToFileMap[r] = make(map[string]Empty)
		}
		mt.VMToFileMap[r][sdfsFileName] = Empty{}
		mt.FileToVMMap[sdfsFileName][r] = Empty{}
	}
}

// Serialize leader state to protobuf format. 
// myAddr is the hostname of the calling node of this function
func LeaderStatesToPB(myAddr string) *pb.LeaderState {
	res := &pb.LeaderState{}
	// Include memtable
	// !TODO: Memtable not lock???
	for key, value := range MemTable.FileToVMMap {
		vm_addr_list, ok := res.FileToVMMap[key]
		if !ok {
			res.FileToVMMap[key] = &pb.LeaderState_AddrList{}
			vm_addr_list = res.FileToVMMap[key]
		}
		vm_list := vm_addr_list.VMAddr
		for vm_addr := range value {
			vm_list = append(vm_list, vm_addr)
		}
	}

	for key, value := range MemTable.VMToFileMap {
		file_list, ok := res.VMToFileMap[key]
		if !ok {
			res.VMToFileMap[key] = &pb.LeaderState_FileList{}
			file_list = res.VMToFileMap[key]
		}
		file_name_list := file_list.FileNames
		for file_name := range value {
			file_name_list = append(file_name_list, file_name)
		}
	}

	// Include file lock
	GlobalFileLock.Lock()
	res.FileLocks = make(map[string]*pb.LeaderState_FileLock)
	for filename, filelock := range FileLocks {
		res.FileLocks[filename] = &pb.LeaderState_FileLock{}
		res.FileLocks[filename].ReadQueue = make([]string, 0)
		res.FileLocks[filename].WriteQueue = make([]string, 0)
		res.FileLocks[filename].ReadQueue = append(res.FileLocks[filename].ReadQueue, filelock.ReadQueue...)
		res.FileLocks[filename].WriteQueue = append(res.FileLocks[filename].WriteQueue, filelock.WriteQueue...)
		res.FileLocks[filename].ReadCount = int32(filelock.ReadCount)
		res.FileLocks[filename].WriteCount = int32(filelock.WriteCount)
		res.FileLocks[filename].ConsecutiveReads = int32(filelock.ConsecutiveReads)
		res.FileLocks[filename].ConsecutiveWrites = int32(filelock.ConsecutiveWrites)
	}
	GlobalFileLock.Unlock()


	// If Leader node, increase Memtable Version number by one
	if myAddr == GetLeaderAddress() {
		Version ++;
	}
	res.Version = int64(Version)

	return res
}

func UpdateLeaderStateIfNecessary(leaderStates *pb.LeaderState) {
	if (leaderStates.Version <= int64(Version) ) {
		return
	}

	// !TODO: Memtable not lock???
	new_file_to_VM_map := make(map[string]map[string]Empty)
	for k, v := range leaderStates.FileToVMMap {
		addr_list := v.VMAddr
		new_file_to_VM_map[k] = make(map[string]Empty)
		for _, addr := range addr_list {
			new_file_to_VM_map[k][addr] = Empty{}
		}
	}
	MemTable.FileToVMMap = new_file_to_VM_map

	new_VM_to_file_map := make(map[string]map[string]Empty)
	for k, v := range leaderStates.VMToFileMap {
		file_list := v.FileNames
		new_VM_to_file_map[k] = make(map[string]Empty)
		for _, filename := range file_list {
			new_VM_to_file_map[k][filename] = Empty{}
		}
	}
	MemTable.VMToFileMap = new_VM_to_file_map

	FileLocks = make(map[string]*FileLock)
	for filename, v := range leaderStates.FileLocks {
		fmt.Printf("%v\n", v)
		FileLocks[filename].ReadQueue = v.ReadQueue
		FileLocks[filename].WriteQueue = v.WriteQueue
		FileLocks[filename].ReadCount = int(v.ReadCount)
		FileLocks[filename].WriteCount = int(v.WriteCount)
		FileLocks[filename].ConsecutiveReads = int(v.ConsecutiveReads)
		FileLocks[filename].ConsecutiveWrites = int(v.ConsecutiveWrites)
	}

	fmt.Printf("My leader state replica's version = %v\n", Version)
	fmt.Printf("My leader state's fileToVMMap = %v\n", MemTable.FileToVMMap)
}




/************************************************* Leader Information ***************************************************************/
var (
	LeaderID = -1 // Record which machine is the leader for now. -1 is the default value and mean there is no leader, at least this machine thinks
)

func GetLeaderID() int {
	return int(LeaderID)
}

func GetLeaderAddress() string {
	id := GetLeaderID()
	var format_parameter string
	if id == 10 {
		format_parameter = strconv.Itoa(id)
	} else {
		format_parameter = "0" + strconv.Itoa(id)
	}
	return fmt.Sprintf("fa23-cs425-18%v.cs.illinois.edu", format_parameter)
}


