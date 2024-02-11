package grpchandlers

import (
	"context"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type log struct {
	req any `json:"req"`
}

func marshalJSON(s *log) ([]byte, error) {
	m, err := json.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("не удалось сериализовать лог: %v", err)
	}
	return m, nil
}

func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler, logger zap.SugaredLogger) (interface{}, error) {
	logItem, err := marshalJSON(&log{
		req: req,
	})
	logger.Info(logItem)
	resp, err := handler(ctx, req)
	return resp, err
}