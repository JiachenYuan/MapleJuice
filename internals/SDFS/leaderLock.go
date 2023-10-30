package SDFS

import (
	"fmt"
<<<<<<< HEAD
	"sync"
)

type RequestType int

const (
	READ RequestType = iota
	WRITE
)

type FileLock struct {
	readQueue         []string
	writeQueue        []string
	readCount         int
	writeCount        int
	consecutiveReads  int
	consecutiveWrites int
	fileLocksMutex    sync.Mutex
}

var (
	fileLocks      = make(map[string]*FileLock)
	globalFileLock = sync.Mutex{}
)

func requestLock(requestorAddress string, fileName string, requestType RequestType) {
	globalFileLock.Lock()
	lock, exists := fileLocks[fileName]
	if !exists {
		lock = &FileLock{}
		fileLocks[fileName] = lock
	}
	globalFileLock.Unlock()

	lock.fileLocksMutex.Lock()
	defer lock.fileLocksMutex.Unlock()
	if requestType == READ {
		lock.readQueue = append(lock.readQueue, requestorAddress)
	} else {
		lock.writeQueue = append(lock.writeQueue, requestorAddress)
=======
	"cs425-mp/internals/global"
)



func requestLock(requestorAddress string, fileName string, requestType global.RequestType) {
	global.GlobalFileLock.Lock()
	lock, exists := global.FileLocks[fileName]
	if !exists {
		lock = &global.FileLock{}
		global.FileLocks[fileName] = lock
	}
	global.GlobalFileLock.Unlock()

	lock.FileLocksMutex.Lock()
	defer lock.FileLocksMutex.Unlock()
	if requestType == global.READ {
		lock.ReadQueue = append(lock.ReadQueue, requestorAddress)
	} else {
		lock.WriteQueue = append(lock.WriteQueue, requestorAddress)
>>>>>>> main
	}
	canProceed := false
	hasPrintedLog := false
	for !canProceed {
		switch requestType {
<<<<<<< HEAD
		case READ:
			canProceed = lock.writeCount == 0 && lock.readCount < 2 && (len(lock.writeQueue) == 0 || lock.consecutiveReads < 4) && lock.readQueue[0] == requestorAddress
		case WRITE:
			canProceed = lock.writeCount == 0 && lock.readCount == 0 && (len(lock.readQueue) == 0 || lock.consecutiveWrites < 4) && lock.writeQueue[0] == requestorAddress
=======
		case global.READ:
			canProceed = lock.WriteCount == 0 && lock.ReadCount < 2 && (len(lock.WriteQueue) == 0 || lock.ConsecutiveReads < 4) && lock.ReadQueue[0] == requestorAddress
		case global.WRITE:
			canProceed = lock.WriteCount == 0 && lock.ReadCount == 0 && (len(lock.ReadQueue) == 0 || lock.ConsecutiveWrites < 4) && lock.WriteQueue[0] == requestorAddress
>>>>>>> main
		}
		if !canProceed {
			if !hasPrintedLog {
				fmt.Printf("Waiting for lock for file %s\n", fileName)
				hasPrintedLog = true
			}
<<<<<<< HEAD
			lock.fileLocksMutex.Unlock()
			// Optionally, sleep for a short time before trying again
			// time.Sleep(time.Millisecond * 50)
			lock.fileLocksMutex.Lock()
		}
	}
	// Grant lock
	if requestType == READ {
		fmt.Printf("Granted read lock for file %s\n", fileName)
		lock.readCount++
		lock.consecutiveReads++
		lock.consecutiveWrites = 0
	} else {
		fmt.Printf("Granted write lock for file %s\n", fileName)
		lock.writeCount++
		lock.consecutiveWrites++
		lock.consecutiveReads = 0
	}
}

func releaseLock(fileName string, requestType RequestType) {
	globalFileLock.Lock()
	lock, exists := fileLocks[fileName]
	globalFileLock.Unlock()
	if !exists {
		return
	}
	lock.fileLocksMutex.Lock()
	defer lock.fileLocksMutex.Unlock()
	if requestType == READ {
		fmt.Printf("Released read lock for file %s\n", fileName)
		lock.readCount--
		lock.readQueue = lock.readQueue[1:]
	} else {
		fmt.Printf("Released write lock for file %s\n", fileName)
		lock.writeCount--
		lock.writeQueue = lock.writeQueue[1:]
=======
			lock.FileLocksMutex.Unlock()
			// Optionally, sleep for a short time before trying again
			// time.Sleep(time.Millisecond * 50)
			lock.FileLocksMutex.Lock()
		}
	}
	// Grant lock
	if requestType == global.READ {
		fmt.Printf("Granted read lock for file %s\n", fileName)
		lock.ReadCount++
		lock.ConsecutiveReads++
		lock.ConsecutiveWrites = 0
	} else {
		fmt.Printf("Granted write lock for file %s\n", fileName)
		lock.WriteCount++
		lock.ConsecutiveWrites++
		lock.ConsecutiveReads = 0
	}
}

func releaseLock(fileName string, requestType global.RequestType) {
	global.GlobalFileLock.Lock()
	lock, exists := global.FileLocks[fileName]
	global.GlobalFileLock.Unlock()
	if !exists {
		return
	}
	lock.FileLocksMutex.Lock()
	defer lock.FileLocksMutex.Unlock()
	if requestType == global.READ {
		fmt.Printf("Released read lock for file %s\n", fileName)
		lock.ReadCount--
		lock.ReadQueue = lock.ReadQueue[1:]
	} else {
		fmt.Printf("Released write lock for file %s\n", fileName)
		lock.WriteCount--
		lock.WriteQueue = lock.WriteQueue[1:]
>>>>>>> main
	}
}
