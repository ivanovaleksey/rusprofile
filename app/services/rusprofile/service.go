package rusprofile

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	lru "github.com/hashicorp/golang-lru"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"strings"
	"time"
)

const cacheSize = 1000

var ErrNotFound = errors.New("rusprofile: not found")

type Service struct {
	data  dataProvider
	cache cacheProvider
}

type dataProvider interface {
	GetData(ctx context.Context, inn string) (io.ReadCloser, error)
}

type cacheProvider interface {
	Add(key, value interface{}) bool
	Get(key interface{}) (interface{}, bool)
}

func NewService(opts ...Option) (*Service, error) {
	srv := &Service{
		data: &webClient{
			httpClient: &http.Client{
				Timeout: 3 * time.Second,
			},
		},
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

type CompanyInfo struct {
	Inn      string
	Kpp      string
	Title    string
	Director string
}

func (srv *Service) GetCompanyInfo(ctx context.Context, inn string) (CompanyInfo, error) {
	entry, ok := srv.cache.Get(inn)
	if ok {
		return entry.(CompanyInfo), nil
	}

	body, err := srv.data.GetData(ctx, inn)
	if err != nil {
		return CompanyInfo{}, errors.Wrap(err, "can't get data")
	}
	defer body.Close()

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return CompanyInfo{}, errors.Wrap(err, "can't build document")
	}

	foundInn := doc.Find("#clip_inn").Text()
	if foundInn == "" {
		return CompanyInfo{}, ErrNotFound
	}

	info := CompanyInfo{
		Inn:      foundInn,
		Kpp:      doc.Find("#clip_kpp").Text(),
		Title:    strings.TrimSpace(doc.Find("h1[itemprop='name']").Text()),
		Director: getDirector(doc),
	}
	srv.cache.Add(inn, info)

	return info, nil
}

func getDirector(doc *goquery.Document) string {
	return doc.Find("span.company-info__title").
		FilterFunction(func(_i int, selection *goquery.Selection) bool {
			return selection.Text() == "Руководитель"
		}).
		SiblingsFiltered("span.company-info__text").
		Text()
}
