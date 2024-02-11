package grpchandlers

import (
	"context"
	"net"

	"github.com/koteyye/shortener/internal/app/deleter"
	"github.com/koteyye/shortener/internal/app/service"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/koteyye/shortener/proto"
)

// GRPCHandlers обработчик grpc запросов
type GRPCHandlers struct {
	services  *service.Service
	logger    *zap.SugaredLogger
	delURLch  chan deleter.DeleteURL
	secretKey string
	subnet    *net.IPNet
	pb.ShortenerServer
}

func InitGRPCHandlers(service *service.Service, logger *zap.SugaredLogger, delURLch chan deleter.DeleteURL, secretKey string, subnet *net.IPNet) *GRPCHandlers {
	return &GRPCHandlers{
		services:  service,
		logger:    logger,
		delURLch:  delURLch,
		secretKey: secretKey,
		subnet:    subnet,
	}
}

func (g *GRPCHandlers) AddURL(ctx context.Context, in *pb.AddURLRequest) (*pb.AddURLResponse, error) {
	var userID string

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		val := md.Get("user")
		if len(val) > 0 {
			userID = val[0]
		} else if len(val) == 0 {
			return nil, status.Errorf(codes.InvalidArgument, "user id is empty")
		}
	}

	result, err := g.services.AddShortURL(ctx, in.Url, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "can't add url: %s", err.Error())
	}
	return &pb.AddURLResponse{Result: result}, nil
}
