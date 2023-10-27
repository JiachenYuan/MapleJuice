package SDFS

import (
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
	}
	canProceed := false
	for !canProceed {
		switch requestType {
		case READ:
			canProceed = lock.writeCount == 0 && lock.readCount < 2 && lock.consecutiveReads < 4 && lock.readQueue[0] == requestorAddress
		case WRITE:
			canProceed = lock.writeCount == 0 && lock.readCount == 0 && lock.consecutiveWrites < 4 && lock.writeQueue[0] == requestorAddress
		}
		if !canProceed {
			lock.fileLocksMutex.Unlock()
			// Optionally, sleep for a short time before trying again
			// time.Sleep(time.Millisecond * 50)
			lock.fileLocksMutex.Lock()
		}
	}
	// Grant lock
	if requestType == READ {
		lock.readCount++
		lock.consecutiveReads++
		lock.consecutiveWrites = 0
	} else {
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
		lock.readCount--
		lock.readQueue = lock.readQueue[1:]
	} else {
		lock.writeCount--
		lock.writeQueue = lock.writeQueue[1:]
	}
}
