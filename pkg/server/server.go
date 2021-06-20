package server

import (
	"context"
	"time"

	"github.com/mikberg/irrigation/pkg/sensing"
	"github.com/mikberg/irrigation/pkg/water"
	pb "github.com/mikberg/irrigation/protobuf"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type ServerConfig struct {
	MoistureSensors  map[uint32]*sensing.MoistureSensor
	Waterer          water.Waterer
	WaterLevelSensor *sensing.WaterLevelSensor
}

type grpcServer struct {
	config *ServerConfig
	pb.UnimplementedIrrigationServer
}

func NewServer(config *ServerConfig) *grpc.Server {
	s := grpc.NewServer()

	server := &grpcServer{
		config: config,
	}

	pb.RegisterIrrigationServer(s, server)
	return s
}

func (s *grpcServer) Water(ctx context.Context, in *pb.WaterRequest) (*pb.WaterResponse, error) {
	log.Info().Msgf("applying water on channel %d for %d seconds", in.GetChannel(), in.GetDuration())

	channel := water.Channel(in.GetChannel())

	// Asynchronously water, so that we can return immediately
	go func() {
		if err := s.config.Waterer.Water(channel, time.Duration(in.GetDuration())*time.Second); err != nil {
			log.Error().Err(err).Msg("failed to water")
		}
	}()

	return &pb.WaterResponse{}, nil
}

func (s *grpcServer) GetRelativeMoisture(ctx context.Context, in *pb.GetRelativeMoistureRequest) (*pb.GetRelativeMoistureResponse, error) {
	channel := in.GetChannel()

	sensor, ok := s.config.MoistureSensors[channel]
	if !ok {
		return nil, grpc.Errorf(codes.InvalidArgument, "no such channel: %d", channel)
	}

	moisture, err := sensor.Read()
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "sensor malfunction: %w", err)
	}

	return &pb.GetRelativeMoistureResponse{
		Moisture: float32(moisture),
	}, nil
}

func (s *grpcServer) GetWaterLevel(ctx context.Context, in *pb.GetWaterLevelRequest) (*pb.GetWaterLevelResponse, error) {
	distance, err := s.config.WaterLevelSensor.Read()
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "sensor malfunction: %w", err)
	}

	return &pb.GetWaterLevelResponse{
		Distance: float32(distance),
		Liters:   float32(40 - distance),
	}, nil
}
