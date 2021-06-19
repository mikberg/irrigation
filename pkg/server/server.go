package server

import (
	"context"

	pb "github.com/mikberg/irrigation/api/v1/irrigation"
	"google.golang.org/grpc"
)

type grpcServer struct {
	pb.UnimplementedIrrigationServer
}

func NewServer() *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterIrrigationServer(s, &grpcServer{})
	return s
}

func (s *grpcServer) Water(ctx context.Context, in *pb.WaterRequest) (*pb.WaterResponse, error) {
	return nil, nil
}
