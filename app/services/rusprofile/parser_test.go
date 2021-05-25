package rusprofile

import (
	"fmt"
	"github.com/ivanovaleksey/rusprofile/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	t.Run("it parses data", func(t *testing.T) {
		testCases := []struct {
			name string
			info models.CompanyInfo
			err  error
		}{
			{
				name: "ozon",
				info: models.CompanyInfo{
					Inn:      "7704217370",
					Kpp:      "770301001",
					Title:    "ООО \"Интернет Решения\"",
					Director: "Шульгин Александр Александрович",
				},
				err: nil,
			},
			{
				name: "rtk",
				info: models.CompanyInfo{
					Inn:      "5030065734",
					Kpp:      "775101001",
					Title:    "ООО \"РТК ИТ\"",
					Director: "Ерохин Виталий Владимирович",
				},
				err: nil,
			},
			{
				name: "not_found",
				info: models.CompanyInfo{},
				err:  ErrNotFound,
			},
		}

		for _, testCase := range testCases {
			p := parser{}

			t.Run(testCase.name, func(t *testing.T) {
				data, err := os.Open(fmt.Sprintf("testdata/%s.html", testCase.name))
				require.NoError(t, err)

				info, err := p.Parse(data)

				require.Equal(t, testCase.err, err)
				assert.Equal(t, testCase.info, info)
			})
		}
	})
}
