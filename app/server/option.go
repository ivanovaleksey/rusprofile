package server

type Option func(*Server)

func WithRusprofileService(srv RusprofileService) Option {
	return func(s *Server) {
		s.rusprofileSrv = srv
	}
}
