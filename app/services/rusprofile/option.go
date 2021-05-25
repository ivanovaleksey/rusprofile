package rusprofile

type Option func(*Service)

func WithDataProvider(p DataProvider) Option {
	return func(s *Service) {
		s.data = p
	}
}

func WithCacheProvider(p CacheProvider) Option {
	return func(s *Service) {
		s.cache = p
	}
}

func WithParser(p Parser) Option {
	return func(s *Service) {
		s.parser = p
	}
}
