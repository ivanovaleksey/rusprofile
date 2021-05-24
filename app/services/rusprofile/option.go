package rusprofile

type Option func(*Service)

func WithDataProvider(p dataProvider) Option {
	return func(s *Service) {
		s.data = p
	}
}

func WithCacheProvider(p cacheProvider) Option {
	return func(s *Service) {
		s.cache = p
	}
}
