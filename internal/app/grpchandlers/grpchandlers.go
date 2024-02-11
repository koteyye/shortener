package grpchandlers

import (
	"context"
	"net"

	"github.com/koteyye/shortener/internal/app/deleter"
	"github.com/koteyye/shortener/internal/app/service"
	"go.uber.org/zap"

	pb "github.com/koteyye/shortener/proto"
)

// GRPCHandlers обработчик grpc запросов
type GRPCHandlers struct {
	services  *service.Service
	logger    *zap.SugaredLogger
	delURLch chan deleter.DeleteURL
	secretKey string
	subnet    *net.IPNet
}

func (g *GRPCHandlers) AddURL(ctx context.Context, in *pb.AddURLRequest) (*pb.AddURLResponse, error) {
	
}