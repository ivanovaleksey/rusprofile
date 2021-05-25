package rusprofile

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebClient_GetData(t *testing.T) {
	inn := randString(7)

	t.Run("with error", func(t *testing.T) {
		fx := newClientFixture(t)
		defer fx.Finish()

		mock := fx.newServerMock(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "GET", r.Method)
			require.Equal(t, "/search?query="+inn, r.URL.String())

			resp := `{}`
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte(resp))
			require.NoError(t, err)
		})
		defer mock.Close()

		data, err := fx.client.GetData(fx.ctx, inn)

		require.Error(t, err)
		assert.Nil(t, data)
	})

	t.Run("without error", func(t *testing.T) {
		fx := newClientFixture(t)
		defer fx.Finish()

		resp := []byte(randString(20))
		mock := fx.newServerMock(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "GET", r.Method)
			require.Equal(t, "/search?query="+inn, r.URL.String())

			w.WriteHeader(http.StatusOK)
			_, err := w.Write(resp)
			require.NoError(t, err)
		})
		defer mock.Close()

		data, err := fx.client.GetData(fx.ctx, inn)

		require.NoError(t, err)
		readData, err := ioutil.ReadAll(data)
		require.NoError(t, err)
		assert.Equal(t, resp, readData)
	})
}

type clientFixture struct {
	t   *testing.T
	ctx context.Context

	client *WebClient
}

func newClientFixture(t *testing.T) *clientFixture {
	fx := &clientFixture{
		t:      t,
		ctx:    context.Background(),
		client: NewWebClient(Config{}),
	}
	return fx
}

func (fx *clientFixture) Finish() {}

func (fx *clientFixture) newServerMock(fn http.HandlerFunc) *httptest.Server {
	srv := httptest.NewServer(fn)
	fx.client.cfg.URL = srv.URL
	return srv
}
