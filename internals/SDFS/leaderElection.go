package SDFS

import (
	"context"
	_ "cs425-mp/internals/failureDetector"
	"cs425-mp/internals/global"
	pb "cs425-mp/protobuf"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Server states enum
type ServerStatus int64

const (
	Leader ServerStatus = iota
	Candidate
	Follower
)

var s *MemberServer

type MemberServer struct {
	pb.UnimplementedLeaderElectionServer
	currentTerm     int64
	leaderID        int32
	votedFor        int32
	state           ServerStatus
	electionDeadline   time.Time
	serverStateLock sync.Mutex // protects the above local states
}

func init() {
	s = newMemberServer()
}

func GetLeaderID() int {
	return int(s.leaderID)
}

func newMemberServer() *MemberServer {
	return &MemberServer{
		currentTerm:     0,
		leaderID:        -1,
		votedFor:        -1,
		state:           Follower,
		serverStateLock: sync.Mutex{},
		electionDeadline:   time.Now().Add(randomDuration()),
	}
}



// HeartBeat implements protobuf.LeaderElectionServer.
func (s *MemberServer) HeartBeat(ctx context.Context, ping *pb.Ping) (*pb.Pong, error) {
	s.serverStateLock.Lock()
	defer s.serverStateLock.Unlock()

	if s.currentTerm <= ping.Term {
		// update current term and convert back to follower
		s.currentTerm = ping.Term
		s.leaderID = ping.LeaderID
		s.state = Follower
		s.electionDeadline = time.Now().Add(randomDuration()) 

		return &pb.Pong{
			Term:    s.currentTerm,
			Success: true,
		}, nil
	}
	// I have equal or higher term number, ignore or reject this heartbeat
	return &pb.Pong{
		Term:    s.currentTerm,
		Success: false,
	}, nil

}

// RequestVotes implements protobuf.LeaderElectionServer.
func (s *MemberServer) RequestVotes(ctx context.Context, request *pb.VoteRequest) (*pb.VoteResponse, error) {
	s.serverStateLock.Lock()
	defer s.serverStateLock.Unlock()

	if s.currentTerm < request.Term {
		// update current term and convert back to follower
		s.currentTerm = request.Term
		s.state = Follower
		s.votedFor = request.CandidateID
		s.electionDeadline = time.Now().Add(randomDuration()) 

		return &pb.VoteResponse{
			Term:        s.currentTerm,
			VoteGranted: true,
		}, nil

	}
	return &pb.VoteResponse{
		Term:        s.currentTerm,
		VoteGranted: false,
	}, nil
}

func startServer() {
	hostname, err := os.Hostname()
	if err != nil {
		return
	}

	lis, err := net.Listen("tcp", hostname + ":" + global.LEADER_ELECTION_PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// Register services
	pb.RegisterLeaderElectionServer(grpcServer, s)

	log.Printf("GRPC server listening on %v", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}


func leaderTask() {
	for {
		// s.serverStateLock.Lock()
		if s.state == Leader {
			// Perform heartbeats 
			// send request vote to every live peers
			s.serverStateLock.Lock()
			originalState := Leader
			originalTerm := s.currentTerm
			localID := int32(getLocalServerID())
			s.serverStateLock.Unlock()
			aliveServerAddrs := getAlivePeersAddrs()
			var wg sync.WaitGroup

			for _, hostname := range aliveServerAddrs {
				wg.Add(1)
				go func(_hostname string) {
					defer wg.Done()
					// Set up a connection to the server
					ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
					conn, err := grpc.DialContext(ctx, _hostname+":"+global.LEADER_ELECTION_PORT) 
					if err != nil {
						log.Fatalf("Failed to dial: %v", err)
					}
					defer conn.Close()
					// Request vote to peers and and respond if states haven't changed since start of the election
					client := pb.NewLeaderElectionClient(conn)
					timeout := 2 * time.Second
					ctx, cancel := context.WithTimeout(context.Background(), timeout)
					defer cancel()
					pong, err := client.HeartBeat(ctx, &pb.Ping{
						Term: originalTerm,
						LeaderID: localID,
					});
					if err != nil {
						// Check if the error is due to a timeout
						if ctx.Err() == context.DeadlineExceeded {
							log.Fatalf("gRPC call timed out after %s", timeout)
							return
						} else {
							log.Fatalf("Failed to call Heartbeat: %v", err)
							return
						}
					}

					s.serverStateLock.Lock()
					if s.state != originalState || s.currentTerm != originalTerm {
						return
					}
					// The other server has a higher term. Yield and convert to Follower
					if pong.Term > s.currentTerm {
						s.currentTerm = pong.Term
						s.state = Follower
						s.electionDeadline = time.Now().Add(randomDuration()) 
					}

					s.serverStateLock.Unlock()

				}(hostname)	
			}
			
		}
		// s.serverStateLock.Unlock()
		time.Sleep(500 * time.Millisecond)
	}
}

func followerTask() {
	lastTerm := int64(-1)
	for {
		s.serverStateLock.Lock()
		
		if (s.state == Follower || s.state == Candidate) {
			if time.Now().After(s.electionDeadline) {
				s.state = Candidate
				originalTerm := s.currentTerm
				s.serverStateLock.Unlock()
				if lastTerm == -1 || originalTerm != lastTerm {
					fmt.Println("Timer goes off, starting election")
					go startElection()
					lastTerm = originalTerm
				}
				continue
			}
		}
		s.serverStateLock.Unlock()
		time.Sleep(500 * time.Millisecond)
	}
}

func startElection() {
	if s.state != Candidate {
		return
	}
	localID := getLocalServerID()
	s.serverStateLock.Lock()
	s.electionDeadline = time.Now().Add(randomDuration())
	fmt.Printf("Next deadline is %v\n", s.electionDeadline)
	originalState := s.state
	s.currentTerm++
	originalTerm := s.currentTerm
	s.votedFor = int32(localID)
	s.serverStateLock.Unlock()
	
	// send request vote to every live peers
	aliveServerAddrs := getAlivePeersAddrs()

	fmt.Printf("Going to send request vote to %v\n", aliveServerAddrs)
	
	var wg sync.WaitGroup
	
	numVotes := 1 // 1 because it votes for itself
	convertedToLeader := false
	for _, hostname := range aliveServerAddrs {
		wg.Add(1)
		go func(_hostname string) {
			defer wg.Done()
			// Set up a connection to the server
			
			ctx, dialCancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer dialCancel()
			conn, err := grpc.DialContext(ctx, _hostname+":"+global.LEADER_ELECTION_PORT, grpc.WithTransportCredentials(insecure.NewCredentials())) 
			
			if err != nil {
				fmt.Printf("Failed to dial: %v\n", err)
				return
			}
			fmt.Println("dialing")
			defer conn.Close()
			// Request vote to peers and and respond if states haven't changed since start of the election
			client := pb.NewLeaderElectionClient(conn)
			timeout := 2 * time.Second
			ctx, callCancel := context.WithTimeout(context.Background(), timeout)
			defer callCancel()
			fmt.Println("sending RequestVotes")
			voteResponse, err := client.RequestVotes(ctx, &pb.VoteRequest{
				Term: originalTerm,
				CandidateID: int32(localID),
			}, );
			if err != nil {
				// Check if the error is due to a timeout
				if ctx.Err() == context.DeadlineExceeded {
					fmt.Printf("gRPC call timed out after %s", timeout)
					return
				} else {
					fmt.Printf("Failed to call RequestVotes: %v", err)
					return
				}
			}
			// Only update local server if state not changed
			s.serverStateLock.Lock()
			defer s.serverStateLock.Unlock()
			if (s.state != originalState || s.currentTerm != originalTerm) {
				return
			}

			if  (!convertedToLeader) && voteResponse.VoteGranted{
				numVotes++
				// if numVotes >= 10/2+1 { // Quorum reached
				if numVotes >= 2 { // ! Testing
					s.state = Leader
					convertedToLeader = true
					return
				}
			}
			if voteResponse.Term > s.currentTerm {
				// update current term and convert back to follower
				s.currentTerm = voteResponse.Term
				s.state = Follower
				s.electionDeadline = time.Now().Add(randomDuration()) 
			}
		}(hostname)
	}
	wg.Wait()

	fmt.Println("Finished election")

}

func startClient() {
	go leaderTask()
	go followerTask()
}

func StartLeaderElection() {
	go startServer()
	startClient()
}

