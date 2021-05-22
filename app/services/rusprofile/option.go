package rusprofile

type Option func(*Service)

func WithDataProvider(p dataProvider) Option {
	return func(s *Service) {
		s.provider = p
	}
}
