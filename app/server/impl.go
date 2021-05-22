package server

import (
	"context"
	"github.com/ivanovaleksey/rusprofile/app/services/rusprofile"
	pb "github.com/ivanovaleksey/rusprofile/pkg/pb/rusprofile"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedRusProfileServiceServer

	rusprofileSrv RusprofileService
}

type RusprofileService interface {
	GetCompanyInfo(ctx context.Context, inn string) (rusprofile.CompanyInfo, error)
}

func NewServer(opts ...Option) *Server {
	srv := &Server{
		rusprofileSrv: rusprofile.NewService(),
	}
	for _, opt := range opts {
		opt(srv)
	}
	return srv
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
