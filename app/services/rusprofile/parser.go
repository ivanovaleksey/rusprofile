package rusprofile

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/ivanovaleksey/rusprofile/pkg/models"
	"github.com/pkg/errors"
	"io"
	"strings"
)

type parser struct{}

func (p parser) Parse(data io.Reader) (models.CompanyInfo, error) {
	doc, err := goquery.NewDocumentFromReader(data)
	if err != nil {
		return models.CompanyInfo{}, errors.Wrap(err, "can't build document")
	}

	inn := doc.Find("#clip_inn").Text()
	if inn == "" {
		return models.CompanyInfo{}, ErrNotFound
	}

	info := models.CompanyInfo{
		Inn:      inn,
		Kpp:      doc.Find("#clip_kpp").Text(),
		Title:    strings.TrimSpace(doc.Find("h1[itemprop='name']").Text()),
		Director: p.getDirector(doc),
	}
	return info, nil
}

func (p parser) getDirector(doc *goquery.Document) string {
	return doc.Find("span.company-info__title").
		FilterFunction(func(_i int, selection *goquery.Selection) bool {
			return selection.Text() == "Руководитель"
		}).
		SiblingsFiltered("span.company-info__text").
		Text()
}
