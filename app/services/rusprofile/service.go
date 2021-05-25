package rusprofile

import (
	"context"
	lru "github.com/hashicorp/golang-lru"
	"github.com/ivanovaleksey/rusprofile/pkg/models"
	"github.com/pkg/errors"
	"io"
)

const cacheSize = 1000

var ErrNotFound = errors.New("rusprofile: not found")

type Service struct {
	data   DataProvider
	cache  CacheProvider
	parser Parser
}

type DataProvider interface {
	GetData(ctx context.Context, inn string) (io.ReadCloser, error)
}

type CacheProvider interface {
	Add(key, value interface{}) bool
	Get(key interface{}) (interface{}, bool)
}

type Parser interface {
	Parse(io.Reader) (models.CompanyInfo, error)
}

func NewService(opts ...Option) (*Service, error) {
	srv := &Service{
		parser: parser{},
	}
	cache, err := lru.New(cacheSize)
	if err != nil {
		return nil, err
	}
	srv.cache = cache

	for _, opt := range opts {
		opt(srv)
	}

	return srv, nil
}

func (srv *Service) Close() error {
	return nil
}

func (srv *Service) GetCompanyInfo(ctx context.Context, inn string) (models.CompanyInfo, error) {
	entry, ok := srv.cache.Get(inn)
	if ok {
		return entry.(models.CompanyInfo), nil
	}

	body, err := srv.data.GetData(ctx, inn)
	if err != nil {
		return models.CompanyInfo{}, errors.Wrap(err, "can't get data")
	}
	defer body.Close()

	info, err := srv.parser.Parse(body)
	if err != nil {
		return models.CompanyInfo{}, err
	}

	srv.cache.Add(inn, info)
	return info, nil
}
