package server

import "go.uber.org/zap"

type Option func(*Server)

func WithRusprofileService(srv RusprofileService) Option {
	return func(s *Server) {
		s.rusprofileSrv = srv
	}
}

func WithLogger(l *zap.Logger) Option {
	return func(s *Server) {
		s.logger = l
	}
}
