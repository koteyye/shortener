package grpchandlers

import (
	"context"
	"strings"

	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	IPHeader   = "X-Real-IP" // IPHeader заголовок запроса, содержащий IP адрес
	methodStat = "/shortener.Shortener/GetStats"
)

var subnetMethods = []string{methodStat}

// SubnetInterceptor проверка IP адреса клиента на вхождение в доверенную подсеть
func (g *GRPCHandlers) SubnetInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var count int
	for i := range subnetMethods {
		if strings.Contains(info.FullMethod, subnetMethods[i]) {
			count += 1
		}
	}
	if count == 0 {
		return handler(ctx, req)
	}
	var ip string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		val := md.Get(IPHeader)
		if len(val) > 0 {
			ip = val[0]
		}
	}
	if ip == "" || !g.subnet.Contains(net.ParseIP(ip)) {
		return nil, status.Error(codes.Unavailable, "метод недоступен из данной подсети")
	}
	return handler(ctx, req)
}
