# CS425-MP2

## Environment
* golang 1.21

## Update
1. Finished the basic data structure design

## Todo
1. Migrate the code to VMs && real IP address rather than port
2. Update membership table (heartbeat count, local timer, etc.)
3. Mark node as failed after detection
4. Think about how to handle the situation of new node joining

## To Run
1. Type in `go run cmd/cs425-mp2/main.go 8080` in the terminal. You will see two things happening
    a. The terminal periodically sends its own gossip table to 8080 and 8001
    b. The current terminal will receive the membership list