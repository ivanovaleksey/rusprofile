package main

import (
	"context"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/ivanovaleksey/rusprofile/app/server"
	"github.com/ivanovaleksey/rusprofile/pkg/pb/rusprofile"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"net/http"
)

func main() {
	logger, _ := zap.NewDevelopment()
	grpc_zap.ReplaceGrpcLoggerV2(logger)

	handler := func(p interface{}) (err error) {
		logger.Error("panic occurred", zap.Any("err", p))
		return status.Errorf(codes.Internal, "%v", p)
	}

	srv := grpc.NewServer(grpc.UnaryInterceptor(
		grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(handler)),
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(logger),
		),
	))
	impl, err := server.NewServer()
	if err != nil {
		logger.Fatal("can't create server", zap.Error(err))
	}

	rusprofile.RegisterRusProfileServiceServer(srv, impl)
	reflection.Register(srv)

	go func() {
		mux := runtime.NewServeMux()
		opts := []grpc.DialOption{
			grpc.WithInsecure(),
		}
		regErr := rusprofile.RegisterRusProfileServiceHandlerFromEndpoint(context.Background(), mux, "127.0.0.1:7001", opts)
		if regErr != nil {
			log.Fatal(regErr)
		}
		err := http.ListenAndServe(":7002", mux)
		fmt.Println(err)
	}()

	lis, _ := net.Listen("tcp", ":7001")
	serveErr := srv.Serve(lis)
	fmt.Println(serveErr)
}
