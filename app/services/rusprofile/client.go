package rusprofile

import (
	"context"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type fileClient struct{}

func (c *fileClient) GetData(_ context.Context, inn string) (io.ReadCloser, error) {
	file, err := os.Open(inn)
	if err != nil {
		return nil, errors.Wrap(err, "can't read file")
	}
	return file, err
}

type WebClient struct {
	cfg        Config
	httpClient httpClient
}

func NewWebClient(cfg Config) *WebClient {
	return &WebClient{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

type Config struct {
	URL string `default:"https://www.rusprofile.ru/"`
}

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

func (c *WebClient) GetData(ctx context.Context, inn string) (io.ReadCloser, error) {
	url := strings.TrimSuffix(c.cfg.URL, "/") + "/search?query=" + inn
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
