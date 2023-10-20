# CS425-MP3

## Environment
* golang 1.21


## TODO
1. Implement Leader Election
2. Read/write lock (maybe a message queue in the leader server?)
3. Embed failure detector in file system. There are two options: 
    a. Make FD and SDFS into two classes and let them communicate with each other using observer mechanism, channel, or callback func
    b. Let FD send message to SDFS's UDP port after detecting failure. 
    

## To Run
