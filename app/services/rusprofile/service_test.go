package rusprofile

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestService_GetCompanyInfo(t *testing.T) {
	t.Run("it parses data", func(t *testing.T) {
		testCases := []struct {
			name string
			info CompanyInfo
			err  error
		}{
			{
				name: "ozon",
				info: CompanyInfo{
					Inn:      "7704217370",
					Kpp:      "770301001",
					Title:    "ООО \"Интернет Решения\"",
					Director: "Шульгин Александр Александрович",
				},
				err:  nil,
			},
			{
				name: "rtk",
				info: CompanyInfo{
					Inn:      "5030065734",
					Kpp:      "775101001",
					Title:    "ООО \"РТК ИТ\"",
					Director: "Ерохин Виталий Владимирович",
				},
				err:  nil,
			},
			{
				name: "not_found",
				info: CompanyInfo{},
				err:  ErrNotFound,
			},
		}

		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				srv, err := NewService(WithDataProvider(&fileClient{}))
				require.NoError(t, err)

				path := fmt.Sprintf("testdata/%s.html", testCase.name)
				info, err := srv.GetCompanyInfo(context.Background(), path)

				require.Equal(t, testCase.err, err)
				assert.Equal(t, testCase.info, info)
			})
		}
	})
}
