package rusprofile

import (
	"context"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"os"
)

const host = "https://www.rusprofile.ru/"

type fileClient struct{}

func (c *fileClient) GetData(_ context.Context, inn string) (io.ReadCloser, error) {
	file, err := os.Open(inn)
	if err != nil {
		return nil, errors.Wrap(err, "can't read file")
	}
	return file, err
}

type webClient struct {
	httpClient httpClient
}

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

func (c *webClient) GetData(ctx context.Context, inn string) (io.ReadCloser, error) {
	url := host + "search?query=" + inn
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "can't create request")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "can't do request")
	}
	if resp.StatusCode != 200 {
		return nil, errors.Errorf("response error, code: %d", resp.StatusCode)
	}

	return resp.Body, nil
}
