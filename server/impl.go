package server

import (
	"github.com/ivanovaleksey/rusprofile/pkg/pb/rusprofile"
)

type Server struct {
	rusprofile.UnimplementedRusProfileServiceServer
}

func NewServer() *Server {
	return new(Server)
}
