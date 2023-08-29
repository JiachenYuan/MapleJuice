# CS425-MP1

## What does this commit do?
1. Basic client-server architecture using socket
2. User input reading
3. Overall project structure
4. A quick python script to generate some log entries
## TODOs
1. Sending request to other VMs (not set up yet)
2. grep command execution

## To run
1. Using one terminal, at the root level of the project, run `go run cmd/cs425-mp1/main.go`. The process represents both a client and server
2. To test it being a client: enter some queries, or just some random strings you like. Now the actual query is not yet being sent to the other endpoints. That is a **TODO**!
3. To test it being a server: on another terminal as client by running `netcat localhost 55555` to connect to the server. You can then type in some words and the server should display those texts.
