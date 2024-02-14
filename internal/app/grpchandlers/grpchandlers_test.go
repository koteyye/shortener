package grpchandlers

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"

	"github.com/koteyye/shortener/config"
	"github.com/koteyye/shortener/internal/app/deleter"
	"github.com/koteyye/shortener/internal/app/models"
	"github.com/koteyye/shortener/internal/app/service"

	mockservice "github.com/koteyye/shortener/internal/app/storage/mocks"
	pb "github.com/koteyye/shortener/proto"
)

const (
	testSecretKey = "super_secret_key"
	tokenHeader   = "token"
	testURL       = "http://yandex.ru"
	testIP        = "127.0.0.1"
	testSubnet    = "127.0.0.1/24"
)

func TestGRPCHandlers_InitGRPCHandlers(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		testService := service.Service{}
		testDelCh := make(chan deleter.DeleteURL)
		grpchandler := InitGRPCHandlers(&testService, &zap.SugaredLogger{}, testDelCh, testSecretKey, &net.IPNet{})

		assert.Equal(t, &GRPCHandlers{
			services:  &testService,
			logger:    &zap.SugaredLogger{},
			secretKey: testSecretKey,
			delURLch:  testDelCh,
			subnet:    &net.IPNet{},
		}, grpchandler)
	})
}

func initTestGRPCHandler(t *testing.T) (*GRPCHandlers, *mockservice.MockURLStorage) {
	c := gomock.NewController(t)
	defer c.Finish()

	repo := mockservice.NewMockURLStorage(c)
	service := service.NewService(repo, &config.Shortener{Listen: "http://localhost:8081"}, &zap.SugaredLogger{})

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	log := *logger.Sugar()

	grpchandler := InitGRPCHandlers(service, &log, make(chan deleter.DeleteURL), testSecretKey, &net.IPNet{})
	return grpchandler, repo
}

func dialer(g *GRPCHandlers) func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			g.AuthInterceptor,
			g.LogInterceptor,
			g.SubnetInterceptor,
		),
	}

	s := grpc.NewServer(opts...)

	pb.RegisterShortenerServer(s, g)

	go func() {
		if err := s.Serve(listener); err != nil {
			panic(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func TestGRPCHandlers_Batch(t *testing.T) {
	testRequest := &pb.BatchRequest{Batch: []*pb.BatchRequestItem{
		{
			CorrelationId: "1",
			OriginalUrl:   "http://yandex.ru",
		},
		{
			CorrelationId: "2",
			OriginalUrl:   "http://rambler.ru",
		},
	}}
	t.Run("batch", func(t *testing.T) {
		g, s := initTestGRPCHandler(t)
		ctx := context.Background()
		conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialer(g)))
		assert.NoError(t, err)
		defer conn.Close()

		client := pb.NewShortenerClient(conn)
		t.Run("success", func(t *testing.T) {
			testToken, err := g.buildJWTString()
			assert.NoError(t, err)
			testUserID, err := g.getUserID(testToken)
			assert.NoError(t, err)

			s.EXPECT().BatchAddURL(gomock.Any(), gomock.Any(), testUserID).Return(error(nil))

			md := metadata.New(map[string]string{tokenHeader: testToken})
			ctx = metadata.NewOutgoingContext(ctx, md)

			res, err := client.Batch(ctx, testRequest)
			assert.NoError(t, err)
			for i := range res.Batch {
				assert.Contains(t, res.Batch[i].ShortUrl, "")
			}
		})
		t.Run("err", func(t *testing.T) {
			testToken, err := g.buildJWTString()
			assert.NoError(t, err)
			testUserID, err := g.getUserID(testToken)
			assert.NoError(t, err)

			s.EXPECT().BatchAddURL(gomock.Any(), gomock.Any(), testUserID).Return(errors.New("some err"))
			md := metadata.New(map[string]string{tokenHeader: testToken})
			ctx = metadata.NewOutgoingContext(ctx, md)

			_, err = client.Batch(ctx, testRequest)
			assert.Error(t, err)
		})
	})
}

func TestGRPCHandlers_AddURL(t *testing.T) {
	testRequest := &pb.AddURLRequest{
		Url: "http://mail.ru",
	}
	t.Run("AddURL", func(t *testing.T) {
		g, s := initTestGRPCHandler(t)
		ctx := context.Background()
		conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialer(g)))
		assert.NoError(t, err)
		defer conn.Close()

		client := pb.NewShortenerClient(conn)
		t.Run("success", func(t *testing.T) {
			testToken, err := g.buildJWTString()
			assert.NoError(t, err)
			testUserID, err := g.getUserID(testToken)
			assert.NoError(t, err)

			s.EXPECT().AddURL(gomock.Any(), gomock.Any(), testRequest.Url, testUserID).Return(error(nil))
			md := metadata.New(map[string]string{tokenHeader: testToken})
			ctx = metadata.NewOutgoingContext(ctx, md)

			res, err := client.AddURL(ctx, testRequest)
			assert.NoError(t, err)
			assert.Contains(t, res.Result, "")
		})
		t.Run("err", func(t *testing.T) {
			testToken, err := g.buildJWTString()
			assert.NoError(t, err)
			testUserID, err := g.getUserID(testToken)
			assert.NoError(t, err)

			s.EXPECT().AddURL(gomock.Any(), gomock.Any(), testRequest.Url, testUserID).Return(errors.New("some err"))
			md := metadata.New(map[string]string{tokenHeader: testToken})
			ctx = metadata.NewOutgoingContext(ctx, md)

			_, err = client.AddURL(ctx, testRequest)
			assert.Error(t, err)
		})
	})
}

func TestGRPCHandlers_GetOriginalURL(t *testing.T) {
	testRequest := &pb.GetOriginalURLRequest{
		ShortUrl: "someShortURL",
	}
	t.Run("get original URL", func(t *testing.T) {
		g, s := initTestGRPCHandler(t)
		ctx := context.Background()
		conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialer(g)))
		assert.NoError(t, err)
		defer conn.Close()

		client := pb.NewShortenerClient(conn)
		t.Run("success", func(t *testing.T) {
			testToken, err := g.buildJWTString()
			assert.NoError(t, err)

			s.EXPECT().GetURL(gomock.Any(), testRequest.ShortUrl).Return(&models.SingleURL{URL: testURL, ShortURL: testRequest.ShortUrl, IsDeleted: false}, error(nil))
			md := metadata.New(map[string]string{tokenHeader: testToken})
			ctx = metadata.NewOutgoingContext(ctx, md)

			res, err := client.GetOriginalURL(ctx, testRequest)
			assert.NoError(t, err)
			assert.Equal(t, testURL, res.OriginalUrl)
		})
		t.Run("err", func(t *testing.T) {
			testToken, err := g.buildJWTString()
			assert.NoError(t, err)

			s.EXPECT().GetURL(gomock.Any(), testRequest.ShortUrl).Return(nil, errors.New("some err"))
			md := metadata.New(map[string]string{tokenHeader: testToken})
			ctx = metadata.NewOutgoingContext(ctx, md)

			_, err = client.GetOriginalURL(ctx, testRequest)
			assert.Error(t, err)
		})
	})
}

func TestGRPCHandlers_GetURLByUserID(t *testing.T) {
	testRequest := &pb.NullRequest{}
	t.Run("get urls by userID", func(t *testing.T) {
		g, s := initTestGRPCHandler(t)
		ctx := context.Background()
		conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialer(g)))
		assert.NoError(t, err)
		defer conn.Close()

		client := pb.NewShortenerClient(conn)
		t.Run("success", func(t *testing.T) {
			testToken, err := g.buildJWTString()
			assert.NoError(t, err)
			testUserID, err := g.getUserID(testToken)
			assert.NoError(t, err)

			wantURLsByUser := []*models.URLList{
				{
					URL:      "http://yandex.ru",
					ShortURL: "http://localhost:8080/someshorturl1",
				},
				{
					URL:      "http://mail.ru",
					ShortURL: "http://localhost:8080/someshorturl2",
				},
			}

			s.EXPECT().GetURLByUser(gomock.Any(), testUserID).Return(wantURLsByUser, error(nil))
			md := metadata.New(map[string]string{tokenHeader: testToken})
			ctx = metadata.NewOutgoingContext(ctx, md)

			res, err := client.GetURLByUserID(ctx, testRequest)
			wantResItem := make([]*pb.GetURLByUserItem, len(wantURLsByUser))
			for i := range wantURLsByUser {
				wantResItem[i] = &pb.GetURLByUserItem{
					OriginalUrl: wantURLsByUser[i].URL,
					ShortUrl:    wantURLsByUser[i].ShortURL,
				}
			}

			assert.NoError(t, err)
			assert.Equal(t, wantResItem, res.Urls)
		})
		t.Run("err", func(t *testing.T) {
			testToken, err := g.buildJWTString()
			assert.NoError(t, err)
			testUserID, err := g.getUserID(testToken)
			assert.NoError(t, err)

			s.EXPECT().GetURLByUser(gomock.Any(), testUserID).Return(nil, errors.New("some err"))
			md := metadata.New(map[string]string{tokenHeader: testToken})
			ctx = metadata.NewOutgoingContext(ctx, md)

			_, err = client.GetURLByUserID(ctx, testRequest)
			assert.Error(t, err)
		})
	})
}

func TestGRPCHandlers_Ping(t *testing.T) {
	testRequest := &pb.NullRequest{}
	t.Run("ping", func(t *testing.T) {
		g, s := initTestGRPCHandler(t)
		ctx := context.Background()
		conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialer(g)))
		assert.NoError(t, err)
		defer conn.Close()

		client := pb.NewShortenerClient(conn)

		testToken, err := g.buildJWTString()
		assert.NoError(t, err)

		s.EXPECT().GetDBPing(gomock.Any()).Return(error(nil))
		md := metadata.New(map[string]string{tokenHeader: testToken})
		ctx = metadata.NewOutgoingContext(ctx, md)

		res, err := client.Ping(ctx, testRequest)
		assert.NoError(t, err)
		assert.Equal(t, "ping success", res.Msg)
	})
}

func TestGRPCHandlers_GetStats(t *testing.T) {
	testRequest := &pb.NullRequest{}
	t.Run("get stats", func(t *testing.T) {
		c := gomock.NewController(t)
		defer c.Finish()

		s := mockservice.NewMockURLStorage(c)
		service := service.NewService(s, &config.Shortener{Listen: "http://localhost:8081"}, &zap.SugaredLogger{})

		logger, err := zap.NewDevelopment()
		if err != nil {
			panic(err)
		}
		defer logger.Sync()
		log := *logger.Sugar()

		cfg := config.Config{TrustSubnet: testSubnet}
		subnet, err := cfg.CIDR()
		assert.NoError(t, err)

		g := InitGRPCHandlers(service, &log, make(chan deleter.DeleteURL), testSecretKey, subnet)
		ctx := context.Background()
		conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialer(g)))
		assert.NoError(t, err)
		defer conn.Close()

		client := pb.NewShortenerClient(conn)
		t.Run("success", func(t *testing.T) {
			s.EXPECT().GetCount(gomock.Any()).Return(1, 1, error(nil))
			ip := net.ParseIP(testIP)

			md := metadata.New(map[string]string{IPHeader: ip.String()})
			ctx = metadata.NewOutgoingContext(ctx, md)

			res, err := client.GetStats(ctx, testRequest)
			wantURL := int32(1)
			wantUsers := int32(1)

			assert.NoError(t, err)
			assert.Equal(t, wantURL, res.Url)
			assert.Equal(t, wantUsers, res.Users)
		})
		t.Run("error", func(t *testing.T) {
			ip := net.ParseIP("127.1.1.0")

			md := metadata.New(map[string]string{IPHeader: ip.String()})
			ctx = metadata.NewOutgoingContext(ctx, md)

			_, err := client.GetStats(ctx, testRequest)
			assert.Error(t, err)
		})
	})
}
