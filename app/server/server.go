package server

import (
	"context"
	"github.com/ivanovaleksey/rusprofile/app/services/rusprofile"
	"github.com/ivanovaleksey/rusprofile/pkg/models"
	pb "github.com/ivanovaleksey/rusprofile/pkg/pb/rusprofile"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
)

type Server struct {
	pb.UnimplementedRusProfileServiceServer

	logger        *zap.Logger
	rusprofileSrv RusprofileService
}

type RusprofileService interface {
	io.Closer
	GetCompanyInfo(ctx context.Context, inn string) (models.CompanyInfo, error)
}

func NewServer(opts ...Option) (*Server, error) {
	srv := &Server{}
	for _, opt := range opts {
		opt(srv)
	}
	return srv, nil
}

func (srv *Server) GetCompanyInfo(ctx context.Context, req *pb.GetCompanyInfoRequest) (*pb.GetCompanyInfoResponse, error) {
	info, err := srv.rusprofileSrv.GetCompanyInfo(ctx, req.Inn)
	switch {
	case err == rusprofile.ErrNotFound:
		return nil, status.Error(codes.NotFound, err.Error())
	case err != nil:
		return nil, err
	}
	resp := &pb.GetCompanyInfoResponse{
		Inn:      info.Inn,
		Kpp:      info.Kpp,
		Title:    info.Title,
		Director: info.Director,
	}
	return resp, nil
}

func (srv *Server) Close() error {
	return srv.rusprofileSrv.Close()
}
