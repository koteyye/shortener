package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"github.com/koteyye/shortener/config"
	"github.com/koteyye/shortener/internal/app/deleter"
	"github.com/koteyye/shortener/internal/app/grpchandlers"
	"github.com/koteyye/shortener/internal/app/handlers"
	"github.com/koteyye/shortener/internal/app/service"
	"github.com/koteyye/shortener/internal/app/storage"
	pb "github.com/koteyye/shortener/proto"
	"github.com/koteyye/shortener/server"

	"net/http"
	_ "net/http/pprof"

	_ "github.com/lib/pq"
)

// @Title Shortener
// @Description Сервис для сокращения URL.
// @Version 1.0

// @Contact.email koteyye@yandex.ru

// @BasePath /
// @Host localhost:8081

// @Tag.name Info
// @Tag.description "Группа запросов состояния сервиса"

// @Tag.name Shortener
// @Tag.desctiption "Группа запросов для сокращения URL"

const (
	shutdownTimeout = 5 * time.Second
)

var (
	buildVersion = "N/A" // Версия сборки.
	buildDate    = "N/A" // Дата сборки.
	buildCommit  = "N/A" // Последний коммит.
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	g, gCtx := errgroup.WithContext(ctx)

	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	log := *logger.Sugar()

	cfg, err := config.GetConfig()
	if err != nil {
		log.Errorw(err.Error(), "event", "get config")
		return
	}
	log.Info("Server address: ", cfg.Server.Listen)
	log.Info("BaseURL: ", cfg.Shortener.Listen)
	log.Info("File storage path: ", cfg.FileStoragePath)
	log.Info("DataBase DN: ", cfg.DataBaseDSN)

	var db *sqlx.DB
	//postgres Client
	if cfg.DataBaseDSN != "" {
		db, err = newPostgres(ctx, cfg)
		if err != nil {
			log.Fatalw(err.Error(), "event", "connect db")
		}
	}
	var subnet *net.IPNet
	if cfg.TrustSubnet != "" {
		subnet, err = cfg.CIDR()
		if err != nil {
			log.Fatal(err.Error(), "event", "cidr")
		}
	}

	//init internal
	storages, err := storage.NewURLHandle(&log, db, cfg.FileStoragePath)
	if err != nil {
		log.Fatalw(err.Error(), "event", "init storage")
	}
	delURLch := make(chan deleter.DeleteURL)
	worker := deleter.InitDeleter(delURLch, storages, &log)
	services := service.NewService(storages, cfg.Shortener, &log)
	handler := handlers.NewHandlers(services, &log, cfg.JWTSecretKey, delURLch, subnet)
	grpchandler := grpchandlers.InitGRPCHandlers(services, &log, delURLch, cfg.JWTSecretKey, subnet)

	if cfg.Pprof != "" {
		g.Go(func() error {
			if startErr := http.ListenAndServe(cfg.Pprof, nil); startErr != nil && !errors.Is(startErr, http.ErrServerClosed) {
				log.Fatalf("cant start server: %s", startErr)
			}
			return nil
		})
	}

	g.Go(func() error {
		runServer(gCtx, cfg, handler, log)
		return nil
	})

	g.Go(func() error {
		worker.StartWorker(gCtx)
		return nil
	})

	if cfg.GRPCServer != "" {
		g.Go(func() error {
			runGRPCServer(gCtx, cfg, grpchandler, &log)
			return nil
		})
	}

	if err = g.Wait(); err != nil {
		log.Fatal(err)
	}
}

func newPostgres(ctx context.Context, cfg *config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", cfg.DataBaseDSN)
	if err != nil {
		return nil, fmt.Errorf("can't create db: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err = db.PingContext(dbCtx)
	if err != nil {
		return nil, fmt.Errorf("can't ping db: %w", err)
	}

	return db, nil
}

func runServer(ctx context.Context, cfg *config.Config, handler *handlers.Handlers, log zap.SugaredLogger) error {
	restServer := new(server.Server)
	go func() {
		if err := restServer.Run(cfg.EnbaleHTTPS, cfg.Server.Listen, handler.InitRoutes(cfg.Server.BaseURL)); err != nil && err != http.ErrServerClosed {
			log.Fatalw(err.Error(), "event", "start server")
		}
	}()

	<-ctx.Done()

	log.Info("shutting down server")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := restServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}

	return nil
}

func runGRPCServer(ctx context.Context, cfg *config.Config, handler *grpchandlers.GRPCHandlers, log *zap.SugaredLogger) error {
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			handler.AuthInterceptor,
			handler.LogInterceptor,
			handler.SubnetInterceptor,
		),
	}

	s := grpc.NewServer(opts...)
	go func() {
		listen, err := net.Listen("tcp", cfg.GRPCServer)
		if err != nil {
			log.Fatalw(err.Error(), "event", "search port for server")
		}
		pb.RegisterShortenerServer(s, handler)

		log.Infof("start grpc server on %v port", cfg.GRPCServer)

		if err := s.Serve(listen); err != nil {
			log.Fatalw(err.Error(), "event", "listen serve")
		}
	}()

	<-ctx.Done()

	log.Info("shutting down grpc server")
	s.GracefulStop()

	return nil
}
