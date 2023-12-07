# MapleJuice

MapleJuice is a distributed systems project built from the ground up using the Go programming language. This system incorporates MapReduce-like functionalities and is backed by what is called Simple Distributed File System, drawing inspiration from Hadoop Distributed File System (HDFS).

## Overview

This project originated as a course assignment (refactored in Golang) in a Distributed Systems course, where the goal was to design and implement a Hadoop-like system. Through the development of MapleJuice, I gained valuable insights into the intricate world of distributed systems.

## Key Features

- **MapReduce Functionality:** MapleJuice leverages a distributed processing paradigm similar to MapReduce, providing a scalable and fault-tolerant framework for processing large datasets.

- **Simple Distributed File System (SDFS):** Our project is supported by a Simple Distributed File System inspired by HDFS, ensuring reliable storage and efficient data retrieval across the distributed environment.

## Environment
* golang 1.21


## How to Start 
1. In directory `~/MapleJuice` and **main** branch, run 
```
go run cmd/MapleJuice/main.go
```

## Supported Commands
* maple <maple_exe> <num_maples> <sdfs_intermediate_filename_prefix> <sdfs_src_directory>
* juice <juice_exe> <num_juices> <sdfs_intermediate_filename_prefix> <sdfs_dest_filename> delete_input={0,1} use_range_partition?={1,0}
  * **Special Note:** for the last command line argument, 1 means use range partitioning and 0 means use hash partitioning
* SQL Query: SELECT ALL FROM Dataset WHERE <regex condition>
* SQL Query:  Given two Datasets D1 and D2: SELECT ALL FROM D1, D2 WHERE <one
specific field’s value in a line of D1 = one specific field’s value in a line of D2>, e.g., “WHERE D1.name = D2.ID”
* get [remote filename] [local filename]
* put [local filename] [remote filename]
* delete [remote filename]
* ls [remote filename]
* store
* **multiread** [remote filename] [local filename] [machine id1] [machine id2] [machine id3] ...
* list_mem
* list_self
* leave
* enable_suspicion
* disable_suspicion
* join (handled implicitly when the node starts)
