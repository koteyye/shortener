package grpchandlers

import (
	"context"
	"database/sql"
	"errors"
	"net"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/koteyye/shortener/internal/app/deleter"
	"github.com/koteyye/shortener/internal/app/models"
	"github.com/koteyye/shortener/internal/app/service"

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

// InitGRPCHandlers возвращает новый экземпляр GRPCHandlers
func InitGRPCHandlers(service *service.Service, logger *zap.SugaredLogger, delURLch chan deleter.DeleteURL, secretKey string, subnet *net.IPNet) *GRPCHandlers {
	return &GRPCHandlers{
		services:  service,
		logger:    logger,
		delURLch:  delURLch,
		secretKey: secretKey,
		subnet:    subnet,
	}
}

// AddURL добавляет сокращенный URL
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

// GetOriginalURL получить оригинальный URL по сокращенному
func (g *GRPCHandlers) GetOriginalURL(ctx context.Context, in *pb.GetOriginalURLRequest) (*pb.GetOriginalURLResponse, error) {
	result, err := g.services.GetOriginURL(ctx, in.ShortUrl)
	if err != nil {
		if errors.Is(err, models.ErrDeleted) {
			return nil, status.Errorf(codes.NotFound, "url is deleted: %s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.GetOriginalURLResponse{OriginalUrl: result}, nil
}

// Batch добавляет множество сокращенных URL
func (g *GRPCHandlers) Batch(ctx context.Context, in *pb.BatchRequest) (*pb.BatchResponse, error) {
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

	batch := make([]*models.URLList, len(in.Batch))
	for i, item := range in.Batch {
		batch[i] = &models.URLList{ID: item.CorrelationId, URL: item.OriginalUrl}
		err := batch[i].Validate()
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		batch[i].ID = strings.TrimSpace(batch[i].ID)
		batch[i].URL = strings.TrimSpace(batch[i].URL)
	}
	list, err := g.services.Batch(ctx, batch, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	res := make([]*pb.BatchResponseItem, len(list))
	for i := range list {
		res[i] = &pb.BatchResponseItem{CorrelationId: list[i].ID, ShortUrl: list[i].ShortURL}
	}
	return &pb.BatchResponse{Batch: res}, nil
}

// GetURLByUserID возвращает все URL пользователя
func (g *GRPCHandlers) GetURLByUserID(ctx context.Context, in *pb.NullRequest) (*pb.GetURLByUserResponse, error) {
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

	list, err := g.services.GetURLByUser(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "the user hasn't urls")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	res := make([]*pb.GetURLByUserItem, len(list))
	for i := range list {
		res[i] = &pb.GetURLByUserItem{OriginalUrl: list[i].URL, ShortUrl: list[i].ShortURL}
	}
	return &pb.GetURLByUserResponse{Urls: res}, nil
}

// Ping пингует базы данных
func (g *GRPCHandlers) Ping(ctx context.Context, in *pb.NullRequest) (*pb.PingResponse, error) {
	err := g.services.GetDBPing(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.PingResponse{Msg: "ping success"}, nil
}

// GetStats получить статистику по сервису, доступен только из доверенной подсети
func (g *GRPCHandlers) GetStats(ctx context.Context, in *pb.NullRequest) (*pb.GetStatsResponse, error) {
	stats, err := g.services.GetStats(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.GetStatsResponse{Url: int32(stats.URLs), Users: int32(stats.Users)}, nil
}

// DeleteURLs запрос на удаление сокращенных URL
func (g *GRPCHandlers) DeleteURLs(ctx context.Context, in *pb.DeleteURLsRequest) (*pb.DeleteURLsResponse, error) {
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

	go func() {
		g.delURLch <- deleter.DeleteURL{URL: in.Urls, UserID: userID}
	}()

	return &pb.DeleteURLsResponse{Msg: "accept"}, nil
}
