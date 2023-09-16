# CS425-MP2

## Environment
* golang 1.21

## Update
1. Migrated the code to VM
2. Designate 1801 as introducer to handle JOIN messages 
3. Update membership table (heartbeat count, local timer, etc.)
4. Mark node as failed after detection

## Todo
1. Add suspicion stage
2. Integrate MP2 to MP1

## To Run
1. Start VM, cd to the project root directory of mp2
2. enter the command `go run cmd/cs425-mp2/main.go`  You will see three things happening
    a. The VM will send a JOIN message to 1801
    b. The VM will start listening gossips from other VMs on port 55556
    b. The current VM will periodically send its membership list to 3 other randomly chosen machines