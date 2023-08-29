# CS425-MP1

## What does this commit do?
1. Using different ports to locally simulate different machines (since we do not have the VMs yet)
2. Finish the client code of sending requests to other servers
3. Finished the server code to handle grep commands
4. The code can now take in consecutive requests, separated by `\n`
5. Changed the project structure by bringing main.go to the root level and changing `internal` to `pkg`
## TODOs
1. Modify grep command execution to allow more query arguments
2. Replace the test .txt files and port numbers with real IP address once we have access to VMs

## To run
1. Using one terminal, at the root level of the project, run `go run main.go 8080`. The process represents both a client and server
2. Open another terminal, at the root level of the project, run `go run main.go 8081`. This command starts another machine that, again, acts both as a client and as a server
3. Using a random terminal, type in `grep 8080` or `grep 8080` to see the results
