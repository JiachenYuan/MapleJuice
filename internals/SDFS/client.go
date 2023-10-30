package SDFS

import (
	"context"
	"cs425-mp/internals/global"
	pb "cs425-mp/protobuf"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// SDFS file operations
func handleGetFile(sdfsFileName string, localFileName string) {
	var conn *grpc.ClientConn
	var c pb.SDFSClient
	var err error

	for {
		if conn == nil {
			conn, err = grpc.Dial(global.GetLeaderAddress()+":"+global.SDFS_PORT, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				fmt.Printf("did not connect: %v\n", err)
				continue
			}
			c = pb.NewSDFSClient(conn)
		}

		shouldWaitForLock := true
		resp := &pb.GetResponse{}
		for shouldWaitForLock {
			r, err := c.GetFile(context.Background(), &pb.GetRequest{
				FileName: sdfsFileName,
			})

			if err != nil {
				fmt.Printf("Failed to call get: %v\n", err)
				conn.Close()
				conn = nil
				time.Sleep(global.RETRY_CONN_SLEEP_TIME)
				fmt.Printf("Retrying to get file %s\n", sdfsFileName)
				break
			}

			if r != nil {
				if !r.Success {
					fmt.Println("failed to acquire the list of vms to read the file from")
					conn.Close()
					return
				}

				if r.ShouldWait {
					fmt.Printf("Waiting for read lock on file %s\n", sdfsFileName)
					time.Sleep(global.RETRY_LOCK_SLEEP_TIME)
				} else {
					shouldWaitForLock = false
					resp = r
				}
			} else {
				time.Sleep(global.RETRY_CONN_SLEEP_TIME)
				fmt.Printf("Retrying to get file %s\n", sdfsFileName)
			}
		}

		if shouldWaitForLock {
			// Failed to acquire lock or reach leader, retry
			continue
		}

		replicas := resp.VMAddresses
		if len(replicas) == 0 {
			fmt.Printf("No target read replicas provided\n")
			conn.Close()
			return
		}
		for _, r := range replicas {
			fmt.Printf("Trying to get file %s from replica: %s\n", sdfsFileName, r)
			remotePath := getScpHostNameFromHostName(r) + ":" + filepath.Join(SDFS_PATH, sdfsFileName)
			cmd := exec.Command("scp", remotePath, localFileName)
			err := cmd.Start()
			if err != nil {
				fmt.Printf("Failed to start command: %v\n", err)
				return
			}

			err = cmd.Wait()
			if err != nil {
				fmt.Printf("Command finished with error: %v\n", err)
				continue
			}
			break
		}

		sendGetACKToLeader(sdfsFileName)
		fmt.Printf("Successfully get file %s\n", sdfsFileName)
		conn.Close()
		break
	}
}

func sendGetACKToLeader(sdfsFileName string) {
	var conn *grpc.ClientConn
	var c pb.SDFSClient
	var err error

	for {
		if conn == nil {
			conn, err = grpc.Dial(global.GetLeaderAddress()+":"+global.SDFS_PORT, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				fmt.Printf("did not connect: %v\n", err)
				continue
			}
			c = pb.NewSDFSClient(conn)
		}

		ackResponse, err := c.GetACK(context.Background(), &pb.GetACKRequest{
			FileName: sdfsFileName,
		})
		if err != nil {
			fmt.Printf("Leader failed to process get ACK: %v\n", err)
			conn.Close()
			conn = nil
			time.Sleep(global.RETRY_CONN_SLEEP_TIME)
			fmt.Printf("Retrying to send get ACK to leader %s\n", sdfsFileName)
			continue
		}
		if ackResponse == nil || !ackResponse.Success {
			fmt.Printf("Leader process get ACK unsuccessfully: %v\n", err)
			conn.Close()
			return
		}
		fmt.Printf("Leader processed get ACK successfully \n")
		conn.Close()
		break
	}
}

func handlePutFile(localFileName string, sdfsFileName string) {
	if _, err := os.Stat(localFileName); os.IsNotExist(err) {
		fmt.Printf("Local file not exist: %s\n", localFileName)
		return
	}

	var conn *grpc.ClientConn
	var c pb.SDFSClient
	var err error

	for {
		// Establish a new connection if it doesn't exist or previous leader failed
		if conn == nil {
			conn, err = grpc.Dial(global.GetLeaderAddress()+":"+global.SDFS_PORT, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				fmt.Printf("did not connect: %v\n", err)
				continue
			}
			c = pb.NewSDFSClient(conn)
		}

		shouldWaitForLock := true
		resp := &pb.PutResponse{}
		for shouldWaitForLock {
			r, err := c.PutFile(context.Background(), &pb.PutRequest{
				FileName: sdfsFileName,
			})

			if err != nil {
				fmt.Printf("Failed to call put: %v\n", err)
				// Close the connection and break to outer loop to retry
				conn.Close()
				conn = nil
				time.Sleep(global.RETRY_CONN_SLEEP_TIME)
				fmt.Printf("Retrying to get file %s\n", sdfsFileName)
				break
			}

			if r != nil {
				if !r.Success {
					fmt.Printf("Failed to put file %s to sdfs %s \n", localFileName, sdfsFileName)
					conn.Close()
					return
				}

				if r.ShouldWait {
					fmt.Printf("Waiting for write lock on file %s\n", sdfsFileName)
					time.Sleep(global.RETRY_LOCK_SLEEP_TIME)
				} else {
					shouldWaitForLock = false
					resp = r
				}
			} else {
				time.Sleep(global.RETRY_CONN_SLEEP_TIME)
				fmt.Printf("Retrying to get file %s\n", sdfsFileName)
			}
		}

		if shouldWaitForLock {
			// Failed to acquire lock or reach leader, retry
			continue
		}

		targetReplicas := resp.VMAddresses
		if len(targetReplicas) == 0 {
			fmt.Printf("No target write replicas provided\n")
			conn.Close()
			return
		}

		fmt.Printf("Starting to put file: %s to SDFS file: %s \n", localFileName, sdfsFileName)
		err = transferFilesConcurrent(localFileName, sdfsFileName, targetReplicas)
		if err != nil {
			fmt.Printf("Failed to transfer file: %v\n", err)
		} else {
			sendPutACKToLeader(sdfsFileName, targetReplicas)
		}
		conn.Close()
		break
	}
}

func transferFilesConcurrent(localFileName string, sdfsFileName string, targetReplicas []string) error {
	var wg sync.WaitGroup
	var transferErrors []error

	for _, r := range targetReplicas {
		wg.Add(1)
		go func(replica string) {
			defer wg.Done()
			err := transferFileToReplica(localFileName, sdfsFileName, replica)
			if err != nil {
				transferErrors = append(transferErrors, err)
			}
		}(r)
	}

	wg.Wait()

	if len(transferErrors) > 0 {
		return fmt.Errorf("Some transfers failed: %v", transferErrors)
	}

	fmt.Printf("Successfully put file to replicas: %+q\n", targetReplicas)
	return nil
}

func transferFileToReplica(localFileName string, sdfsFileName string, replica string) error {
	targetHostName := getScpHostNameFromHostName(replica)
	remotePath := targetHostName + ":" + filepath.Join(SDFS_PATH, sdfsFileName)
	cmd := exec.Command("scp", localFileName, remotePath)
	err := cmd.Start()
	if err != nil {
		fmt.Printf("Failed to start command: %v\n", err)
		return err
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Printf("Command finished with error: %v\n", err)
	}

	return err
}

func sendPutACKToLeader(sdfsFileName string, targetReplicas []string) {
	var conn *grpc.ClientConn
	var c pb.SDFSClient
	var err error

	for {
		if conn == nil {
			conn, err = grpc.Dial(global.GetLeaderAddress()+":"+global.SDFS_PORT, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				fmt.Printf("did not connect: %v\n", err)
				continue
			}
			c = pb.NewSDFSClient(conn)
		}
		r, err := c.PutACK(context.Background(), &pb.PutACKRequest{
			FileName:         sdfsFileName,
			ReplicaAddresses: targetReplicas,
		})
		if err != nil {
			fmt.Printf("Leader failed to process put ACK: %v\n", err)
			conn.Close()
			conn = nil
			time.Sleep(global.RETRY_CONN_SLEEP_TIME)
			fmt.Printf("Retrying to send put ACK to leader %s\n", sdfsFileName)
			continue
		}
		if r == nil || !r.Success {
			fmt.Printf("Leader process put ACK unsuccessfully: %v\n", err)
			conn.Close()
			return
		}
		fmt.Printf("Leader processed put ACK successfully \n")
		conn.Close()
		break
	}
}

func handleDeleteFile(sdfsFileName string) {
	conn, err := grpc.Dial(global.GetLeaderAddress()+":"+global.SDFS_PORT, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("did not connect: %v\n", err)
	}
	defer conn.Close()

	c := pb.NewSDFSClient(conn)
	r, err := c.DeleteFileLeader(context.Background(), &pb.DeleteRequestLeader{
		FileName: sdfsFileName,
	})

	if err != nil {
		fmt.Printf("Failed to call delete: %v\n", err)
	}
	if r.Success {
		fmt.Printf("Successfully deleted file %s\n", sdfsFileName)
	} else {
		fmt.Printf("Failed to delete file %s\n", sdfsFileName)
	}
}

func handleListFileHolders(sdfsFileName string) {
	ctx, dialCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer dialCancel()
	conn, err := grpc.DialContext(ctx, global.GetLeaderAddress()+":"+global.SDFS_PORT, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("did not connect: %v\n", err)
	}
	defer conn.Close()

	c := pb.NewSDFSClient(conn)
	timeout := 2 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	r, err := c.ListFileHolder(ctx, &pb.ListFileHolderRequest{
		FileName: sdfsFileName,
	})

	if err != nil {
		// Check if the error is due to a timeout
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Printf("gRPC call timed out after %s\n", timeout)
		} else {
			fmt.Printf("Failed to call ls: %v\n", err)
		}
	}
	if r.Success {
		fmt.Printf("%+q\n", r.VMAddresses)
	} else {
		fmt.Printf("Failed to list VMs for file: %s\n", sdfsFileName)
	}
}

func handleListLocalFiles() {
	ctx, dialCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer dialCancel()
	conn, err := grpc.DialContext(ctx, global.GetLeaderAddress()+":"+global.SDFS_PORT, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("did not connect: %v\n", err)
	}
	defer conn.Close()

	c := pb.NewSDFSClient(conn)
	timeout := 2 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	r, err := c.ListLocalFiles(ctx, &pb.ListLocalFilesRequest{
		SenderAddress: HOSTNAME,
	})

	if err != nil {
		// Check if the error is due to a timeout
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Printf("gRPC call timed out after %s\n", timeout)
		} else {
			fmt.Printf("Failed to call store: %v\n", err)
		}
	}
	if r.Success {
		fmt.Printf("%+q\n", r.FileNames)
	} else {
		fmt.Printf("Failed to list local files for machine %s\n", HOSTNAME)
	}
}
