package main

import (
	"github.com/ivanovaleksey/rusprofile/pkg/closer"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"io"
	"time"
)

const closeTimeout = 3 * time.Second

type grpcCloser struct {
	srv *grpc.Server
}

func (srv grpcCloser) Close() error {
	srv.srv.GracefulStop()
	return nil
}

type closerWrapper struct {
	io.Closer
	logger *zap.Logger
}

func newCloserWrapper(srv io.Closer, logger *zap.Logger) closerWrapper {
	return closerWrapper{
		Closer: srv,
		logger:   logger,
	}
}

func (c closerWrapper) Close() error {
	c.logger.Info("closing")
	tc := closer.NewTimeoutCloser(c.Closer, closeTimeout)
	if err := tc.Close(); err != nil {
		c.logger.Error("close error", zap.Error(err))
		return err
	}
	c.logger.Info("closed")
	return nil
}
