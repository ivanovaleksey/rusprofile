package main

import (
	"context"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/ivanovaleksey/rusprofile/app/config"
	"github.com/ivanovaleksey/rusprofile/app/server"
	"github.com/ivanovaleksey/rusprofile/pkg/closer"
	"github.com/ivanovaleksey/rusprofile/pkg/pb/rusprofile"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, _ := zap.NewDevelopment()
	grpc_zap.ReplaceGrpcLoggerV2(logger)

	cfg, err := config.New()
	if err != nil {
		logger.Fatal("can't create config", zap.Error(err))
	}

	appCloser := closer.New(syscall.SIGTERM, syscall.SIGINT)
	appCloser.Add(func() error {
		cancel()
		return nil
	})

	runApps(ctx, cfg, logger, appCloser)

	appCloser.Wait()
}

func runApps(ctx context.Context, cfg config.Config, logger *zap.Logger, appCloser *closer.Closer) {
	srv, err := server.NewServer(server.WithLogger(logger))
	if err != nil {
		logger.Fatal("can't create server", zap.Error(err))
	}
	{
		wrap := newCloserWrapper(srv, logger.With(zap.String("component", "app")))
		appCloser.Add(wrap.Close)
	}

	go func() {
		l := logger.With(zap.String("component", "grpc"))

		grpcSrv := grpc.NewServer(grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_recovery.UnaryServerInterceptor(),
				grpc_ctxtags.UnaryServerInterceptor(),
				grpc_zap.UnaryServerInterceptor(logger),
			),
		))

		rusprofile.RegisterRusProfileServiceServer(grpcSrv, srv)
		reflection.Register(grpcSrv)

		wrap := newCloserWrapper(grpcCloser{srv: grpcSrv}, l)
		appCloser.Add(wrap.Close)

		lis, _ := net.Listen("tcp", cfg.GRPCAddr)
		l.Info("starting")
		if err := grpcSrv.Serve(lis); err != nil {
			l.Error("serve error", zap.Error(err))
			return
		}
	}()

	go func() {
		const (
			httpReadTimeout     = 3 * time.Second
			httpWriteTimeout    = 3 * time.Second
			httpShutdownTimeout = 3 * time.Second
		)

		l := logger.With(zap.String("component", "http"))

		mux := runtime.NewServeMux()
		opts := []grpc.DialOption{
			grpc.WithInsecure(),
		}
		err := rusprofile.RegisterRusProfileServiceHandlerFromEndpoint(ctx, mux, cfg.GRPCAddr, opts)
		if err != nil {
			l.Fatal("can't register http gateway", zap.Error(err))
		}

		srv := &http.Server{
			Addr:         cfg.HTTPAddr,
			Handler:      mux,
			ReadTimeout:  httpReadTimeout,
			WriteTimeout: httpWriteTimeout,
		}

		appCloser.Add(func() error {
			ctx, cancel := context.WithTimeout(context.Background(), httpShutdownTimeout)
			defer cancel()

			l.Info("closing")
			if err := srv.Shutdown(ctx); err != nil {
				l.Error("close error", zap.Error(err))
				return err
			}
			l.Info("closed")
			return nil
		})

		l.Info("starting")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Error("serve error", zap.Error(err))
			return
		}
	}()
}
