package grpcserver

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"shortener/internal/interceptors"
	"shortener/internal/logger"
	"shortener/internal/models"
	"shortener/internal/service"
	"shortener/internal/storage"
	pb "shortener/pkg/service/proto"
)

type GRPCServer struct {
	pb.UnimplementedURLShortenerServiceServer
	svc *service.Service
}

// GRPCServe runs the gRPC server.
func GRPCServe(svc *service.Service, log *logger.Log) error {
	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		return fmt.Errorf("failed to listen a port: %w", err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(interceptors.UserIDUnaryInterceptor(svc)))
	pb.RegisterURLShortenerServiceServer(s, &GRPCServer{svc: svc})
	reflection.Register(s)
	log.Debug("Starting gRPC server..")

	return s.Serve(listen)
}

// Save method saves long url and replies short one.
func (g *GRPCServer) Save(ctx context.Context, long *wrapperspb.StringValue) (*wrapperspb.StringValue, error) {
	if long == nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid URL passed")
	}

	short, err := g.svc.SaveURL(ctx, long.String())
	if err != nil {
		var duplicateError *storage.DuplicateRecordError
		if errors.As(err, &duplicateError) {
			return nil, status.Error(codes.AlreadyExists, duplicateError.Message)
		}
		g.svc.Log.Err("failed to save URL", err)
		return nil, status.Error(codes.Internal, "failed to save URL")
	}

	return &wrapperspb.StringValue{Value: g.svc.BaseURL + "/" + short}, nil
}

// Ping checks if connection to database can be established.
func (g *GRPCServer) Ping(ctx context.Context, _ *pb.PingRequest) (*pb.PingResponse, error) {
	if err := g.svc.Storage.Ping(ctx); err != nil {
		g.svc.Log.Err("failed to ping database", err)
		return nil, status.Error(codes.Unavailable, "")
	}

	return &pb.PingResponse{}, nil
}

// SavedByUser method gets saved urls by the user from the ctx.
func (g *GRPCServer) SavedByUser(ctx context.Context, _ *pb.SavedByUserRequest) (*pb.SavedByUserResponse, error) {
	urls, err := g.svc.GetUserURLs(ctx)
	if err != nil {
		g.svc.Log.Err("failed to get user urls", err)
		return nil, status.Error(codes.Internal, "")
	}
	result := &pb.SavedByUserResponse{}
	for _, url := range urls {
		tmp := &pb.URL{OriginalUrl: url.OriginalURL, ShortUrl: url.ShortURL}
		result.Urls = append(result.Urls, tmp)
	}
	return result, nil
}

// Get long URL by short value.
func (g *GRPCServer) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	long, err := g.svc.GetURL(ctx, in.GetShort())
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrURLDeleted):
			g.svc.Log.Info("requested deleted url", "short", in.GetShort())
			return nil, status.Error(codes.Unavailable, "Requested deleted URL")
		case errors.Is(err, service.ErrURLNotFound):
			return nil, status.Error(codes.NotFound, "Requested URL not found")
		default:
			g.svc.Log.Err("failed to get URL", err)
			return nil, status.Error(codes.Internal, "")
		}
	}

	return &pb.GetResponse{Long: g.svc.BaseURL + "/" + long}, nil
}

// Batch saves many urls for the one call.
func (g *GRPCServer) Batch(ctx context.Context, in *pb.BatchRequest) (*pb.BatchResponse, error) {
	req := make([]models.BatchRequest, 0)
	for _, u := range in.GetUrls() {
		req = append(req, models.BatchRequest{
			OriginalURL: u.GetOriginalUrl(), CorrelationID: u.GetCorrelationId()},
		)
	}
	saved, err := g.svc.SaveURLs(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, "")
	}
	res := &pb.BatchResponse{}
	for _, u := range saved {
		res.Urls = append(
			res.Urls,
			&pb.BatchResponseEntity{
				CorrelationId: u.CorrelationID,
				ShortUrl:      u.ShortURL,
			})
	}

	return &pb.BatchResponse{Urls: res.GetUrls()}, nil
}

// DeleteMany deletes many urls for the one call.
func (g *GRPCServer) DeleteMany(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	if err := g.svc.DeleteURLs(ctx, in.GetUrls()); err != nil {
		g.svc.Log.Err("failed to delete URLs", err)
		return nil, status.Error(codes.Internal, "")
	}

	return &pb.DeleteResponse{Urls: in.GetUrls()}, nil
}

// Shorten method saves long and returns short url.
func (g *GRPCServer) Shorten(ctx context.Context, in *pb.ShortenRequest) (*pb.ShortenResponse, error) {
	short, err := g.svc.SaveURL(ctx, in.GetUrl())
	if err != nil {
		var duplicateErr *storage.DuplicateRecordError
		if errors.As(err, &duplicateErr) {
			g.svc.Log.Warn("failed to save url", err)
			duplicate := g.svc.BaseURL + "/" + duplicateErr.Message
			return nil, status.Error(codes.AlreadyExists, duplicate)
		} else {
			g.svc.Log.Err("failed to save url", err)
			return nil, status.Error(codes.Internal, "")
		}
	}

	return &pb.ShortenResponse{Result: g.svc.BaseURL + "/" + short}, nil
}

// Stats method shows internal info about saved users and urls.
func (g *GRPCServer) Stats(ctx context.Context, _ *pb.StatsRequest) (*pb.StatsResponse, error) {
	const permissionDeniedMsg = "Untrusted subnet"
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, permissionDeniedMsg)
	}
	tcpAddr, ok := p.Addr.(*net.TCPAddr)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, permissionDeniedMsg)
	}

	if !g.svc.IsSubnetTrusted(tcpAddr.IP.String()) {
		return nil, status.Error(codes.PermissionDenied, permissionDeniedMsg)
	}

	stats, err := g.svc.GetStats(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &pb.StatsResponse{Urls: strconv.Itoa(stats.URLs), Users: strconv.Itoa(stats.Users)}, nil
}
