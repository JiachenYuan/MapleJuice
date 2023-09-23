# CS425-MP2

## Environment
* golang 1.21


## To Run
<<<<<<< HEAD
1. Start introducer VM first (fa23-cs425-1801.cs.illinois.edu), cd to the project root directory of mp2
2. enter the command `go run cmd/cs425-mp2/main.go` to start the introducer
3. Start other VMs, cd to the project root directory of mp2, and run `go run cmd/cs425-mp2/main.go` to start other membership processes
4. We have provided a bunch of CLI commands for you to interact with each VM, and they are listed below:
    a. ##leave##. This command will ask the current VM to gracefully leave the system and enter the "left" mode. That is, it does neither response to any external gossips nor send any gossip messages, but the process is still running. 
    b. ##rejoin## This command will enable a process in the "left" mode to rejoin the system with a new node ID by sending a introduction message to the introducer.
    c. ##list_mem## This command will list all the entries in the current VM's membership list. 
    d. ##list_self## This command will list the node ID of the current VM. 
    e. ##suspicion_toggl## This command will make the system toggle between Gossip mode and Gossip + S mode. 
    f. ##drop_rate {rate}## The command will read the {rate} variable in the user input and use it to adjust the message drop rate of the system. You can only input a floting point number between 0.0 and 1.0. 
=======
1. Log into (fa23-cs425-1801.cs.illinois.edu), cd to the project root directory of mp2. Run `go run cmd/cs425-mp2/main.go` to start the Introducer node
2. On other VMs, run `go run cmd/cs425-mp2/main.go` to start member nodes
3. We have provided a bunch of CLI commands to interact with each VM, and they are listed below:
*  **leave**. This command will ask the current VM to gracefully leave the system and enter the "LEFT" mode. That is, it neither respond to any external gossips nor send any gossip messages.
* **rejoin** This command will bring a "LEFT" node back to the group with a new node ID by sending a introduction message to the Introducer node.
* **list_mem** This command will list all the entries in the current VM's membership list. 
* **list_self** This command will list the node ID of the current VM. 
* **suspicion_toggle** This command will make the system toggle between Gossip mode and Gossip + S mode. 
* **drop_rate {rate}** The command will read the {rate} variable in the user input and use it to adjust the message drop rate of the system. You can only input a floting point number between 0.0 and 1.0. 

4. The logs are written to `log.log`
>>>>>>> main
