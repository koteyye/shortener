package grpchandlers

import (
	"context"
	"encoding/json"
	"time"

	"google.golang.org/grpc"
)

type log struct {
	Method   string `json:"method"`
	Body     any    `json:"body"`
	Duration int64  `json:"duration"`
}

func (l *log) marshalJSON() []byte {
	m, _ := json.Marshal(l)
	return m
}

// LogInterceptor логирование запросов
func (g *GRPCHandlers) LogInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	res, err := handler(ctx, req)

	duration := time.Since(start).Milliseconds()
	logItem := &log{
		Method:   info.FullMethod,
		Duration: duration,
		Body:     req,
	}

	g.logger.Infow("GRPC Request", "event", string(logItem.marshalJSON()))
	return res, err
}
