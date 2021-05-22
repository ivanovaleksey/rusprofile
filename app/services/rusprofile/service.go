package rusprofile

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"strings"
	"time"
)

var ErrNotFound = errors.New("rusprofile: not found")

type Service struct {
	provider dataProvider
}

type dataProvider interface {
	GetData(ctx context.Context, inn string) (io.ReadCloser, error)
}

func NewService(opts ...Option) *Service {
	srv := &Service{
		provider: &webClient{
			httpClient: &http.Client{
				Timeout: 3 * time.Second,
			},
		},
	}
	for _, opt := range opts {
		opt(srv)
	}
	return srv
}

type CompanyInfo struct {
	Inn      string
	Kpp      string
	Title    string
	Director string
}

func (srv *Service) GetCompanyInfo(ctx context.Context, inn string) (CompanyInfo, error) {
	body, err := srv.provider.GetData(ctx, inn)
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
