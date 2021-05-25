package rusprofile

import (
	"bytes"
	"context"
	"github.com/ivanovaleksey/rusprofile/app/services/rusprofile/mocks"
	"github.com/ivanovaleksey/rusprofile/pkg/models"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"math/rand"
	"testing"
)

func TestService_GetCompanyInfo(t *testing.T) {
	inn := randString(7)
	info := models.CompanyInfo{
		Inn:      randString(7),
		Kpp:      randString(7),
		Title:    randString(10),
		Director: randString(20),
	}

	t.Run("with data in cache", func(t *testing.T) {
		fx := newFixture(t)
		defer fx.Finish()

		fx.cacheProvider.On("Get", inn).Return(info, true)

		resp, err := fx.srv.GetCompanyInfo(fx.ctx, inn)
		require.NoError(t, err)

		assert.Equal(t, info, resp)
	})

	t.Run("without data in cache", func(t *testing.T) {
		t.Run("with client error", func(t *testing.T) {
			fx := newFixture(t)
			defer fx.Finish()

			clientErr := errors.New(randString(10))
			fx.cacheProvider.On("Get", inn).Return(models.CompanyInfo{}, false)
			fx.dataProvider.On("GetData", fx.ctx, inn).Return(nil, clientErr)

			resp, err := fx.srv.GetCompanyInfo(fx.ctx, inn)
			require.Equal(t, clientErr, errors.Cause(err))

			assert.Empty(t, resp)
		})

		t.Run("without client error", func(t *testing.T) {
			t.Run("with parser error", func(t *testing.T) {
				fx := newFixture(t)
				defer fx.Finish()

				fx.cacheProvider.On("Get", inn).Return(models.CompanyInfo{}, false)
				data := ioutil.NopCloser(bytes.NewBufferString(randString(10)))
				fx.dataProvider.On("GetData", fx.ctx, inn).Return(data, nil)
				parserErr := errors.New(randString(10))
				fx.parser.On("Parse", data).Return(models.CompanyInfo{}, parserErr)

				resp, err := fx.srv.GetCompanyInfo(fx.ctx, inn)
				require.Equal(t, parserErr, err)

				assert.Empty(t, resp)
			})

			t.Run("without parser error", func(t *testing.T) {
				fx := newFixture(t)
				defer fx.Finish()

				fx.cacheProvider.On("Get", inn).Return(models.CompanyInfo{}, false)
				data := ioutil.NopCloser(bytes.NewBufferString(randString(10)))
				fx.dataProvider.On("GetData", fx.ctx, inn).Return(data, nil)
				fx.parser.On("Parse", data).Return(info, nil)
				fx.cacheProvider.On("Add", inn, info).Return(false)

				resp, err := fx.srv.GetCompanyInfo(fx.ctx, inn)
				require.NoError(t, err)

				assert.Equal(t, info, resp)
			})
		})
	})
}

type fixture struct {
	t   *testing.T
	ctx context.Context

	srv           *Service
	cacheProvider *mocks.CacheProvider
	dataProvider  *mocks.DataProvider
	parser        *mocks.Parser
}

func newFixture(t *testing.T) *fixture {
	fx := &fixture{
		t:   t,
		ctx: context.Background(),

		cacheProvider: &mocks.CacheProvider{},
		dataProvider:  &mocks.DataProvider{},
		parser:        &mocks.Parser{},
	}
	srv, err := NewService(
		WithCacheProvider(fx.cacheProvider),
		WithDataProvider(fx.dataProvider),
		WithParser(fx.parser),
	)
	require.NoError(t, err)
	fx.srv = srv
	return fx
}

func (fx *fixture) Finish() {
	fx.cacheProvider.AssertExpectations(fx.t)
	fx.dataProvider.AssertExpectations(fx.t)
	fx.parser.AssertExpectations(fx.t)
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
