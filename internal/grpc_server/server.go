package grpc_server

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"shortener/internal/logger"
	"shortener/internal/service"
	pb "shortener/pkg/service/proto"
)

type GRPCServer struct {
	pb.UnimplementedURLShortenerServiceServer
	svc  *service.Service
	Addr string
}

// Stats method shows internal info about saved users and urls.
func (g *GRPCServer) Stats(ctx context.Context, _ *pb.StatsRequest) (*pb.StatsResponse, error) {
	const permissionDeniedMsg = "Untrusted subnet"
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, permissionDeniedMsg)
	}
	realIP := p.Addr.(*net.TCPAddr).IP.String()

	if !g.svc.IsSubnetTrusted(realIP) {
		return nil, status.Error(codes.PermissionDenied, permissionDeniedMsg)
	}

	stats, err := g.svc.GetStats(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &pb.StatsResponse{Urls: strconv.Itoa(stats.URLs), Users: strconv.Itoa(stats.Users)}, nil
}

// GRPCServe runs the gRPC server.
func GRPCServe(svc *service.Service, log *logger.Log) error {
	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		return fmt.Errorf("failed to listen a port: %w", err)
	}
	s := grpc.NewServer()
	pb.RegisterURLShortenerServiceServer(s, &GRPCServer{svc: svc})
	reflection.Register(s)
	log.Debug("Starting gRPC server..")

	return s.Serve(listen)
}
