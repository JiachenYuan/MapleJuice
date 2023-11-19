package maplejuice

import (
	"context"
	sdfs "cs425-mp/internals/SDFS"
	"cs425-mp/internals/global"
	pb "cs425-mp/protobuf"
	"fmt"
	"net"

	"google.golang.org/grpc"
)

type MapleJuiceServer struct {
	pb.UnimplementedMapleJuiceServer
}

func (s *MapleJuiceServer) GetMapleWorkerList(ctx context.Context, in *pb.MapleWorkerListRequest) (*pb.MapleWorkerListeResponse, error) {
	if sdfs.IsCurrentNodeLeader() {
		workersCount := in.NumMaples
		dir := in.SdfsSrcDirectory
		workerAssignments := assignWorkload(dir, int(workersCount))
		resp := &pb.MapleWorkerListeResponse{
			Success:     true,
			Assignments: workerAssignments,
		}
		return resp, nil
	} else {
		fmt.Println("Not the leader, cannot get maple worker list")
		return nil, fmt.Errorf("not the leader, cannot get maple worker list")
	}
}

func (s *MapleJuiceServer) Maple(ctx context.Context, in *pb.MapleRequest) (*pb.MapleResponse, error) {
	return nil, nil
}

func StartMapleJuiceServer() {
	// start listening for incoming connections
	lis, err := net.Listen("tcp", ":"+global.MAPLE_JUICE_PORT)
	if err != nil {
		fmt.Printf("failed to listen: %v\n", err)
	}
	s := grpc.NewServer()
	pb.RegisterMapleJuiceServer(s, &MapleJuiceServer{})
	if err := s.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v\n", err)
	}
}
