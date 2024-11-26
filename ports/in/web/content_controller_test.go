package web

import (
	"contentgit/testdata/testserver"
	"contentgit/testdata/testsuite"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ContentControllerTestSuite struct {
	testsuite.BaseDatabaseTestSuite
}

func TestContentControllerTestSuite(t *testing.T) {
	suite.Run(t, new(ContentControllerTestSuite))
}

func (suite *ContentControllerTestSuite) TestCreateBulkContent() {
	// given
	sut := testserver.NewTestAppServerBuilder(Router{}, suite.TestDbContainer).Build()

	requestBody := `[
		{
		"name": "2024 최신형 공기 살균기",
		"mainImage": "https://gdimg.gmarket.co.kr/2367233519/still/280?ver=1645526559",
				"price": "250000",
				"taxRate": "10.2"
		},
			{
		"name": "이지듀 멜라토닝 원데이 앰플 대용량 28ml",
		"mainImage": "https://cdn.011st.com/11dims/resize/600x600/quality/75/11src/product/5966707693/B.jpg?920000000",
				"price": "3000",
				"taxRate": "10.2",
				"liveShowInventoryQuantity": 2000
		}
	]`

	req := httptest.NewRequest(http.MethodPost, "/api/tenants/bettercode/contents/bulk", strings.NewReader(requestBody))
	rec := httptest.NewRecorder()

	// when
	sut.ServeHTTP(rec, req)
	fmt.Println(rec.Body.String())

	// then
	suite.Equal(http.StatusCreated, rec.Code)
}

func (suite *ContentControllerTestSuite) TestCreateContent() {
	// given
	sut := testserver.NewTestAppServerBuilder(Router{}, suite.TestDbContainer).Build()

	requestBody := `{
		"name": "2024 최신형 공기 살균기",
		"mainImage": "https://gdimg.gmarket.co.kr/2367233519/still/280?ver=1645526559",
		"price": "250000",
		"taxRate": "10.2"
	}`
	req := httptest.NewRequest(http.MethodPost, "/api/tenants/bettercode/contents", strings.NewReader(requestBody))
	rec := httptest.NewRecorder()

	// when
	sut.ServeHTTP(rec, req)

	// then
	suite.Equal(http.StatusCreated, rec.Code)
}

func (suite *ContentControllerTestSuite) TestGetContent() {
	// given
	sut := testserver.NewTestAppServerBuilder(Router{}, suite.TestDbContainer).DatabaseFixture().Build()

	req := httptest.NewRequest(http.MethodGet, "/api/tenants/bettercode/contents/6f3bbc99-55aa-4340-89f6-1ddd4dfdb8cd", nil)
	rec := httptest.NewRecorder()

	// when
	sut.ServeHTTP(rec, req)
	fmt.Println(rec.Body.String())

	// then
	suite.Equal(http.StatusOK, rec.Code)

	var actual any
	json.Unmarshal(rec.Body.Bytes(), &actual)

	expected := map[string]any{
		"id": "6f3bbc99-55aa-4340-89f6-1ddd4dfdb8cd",
		"content": map[string]any{
			"name":                      "2024 최신형 공기 살균기",
			"mainImage":                 "https://cdn.011st.com/11dims/resize/600x600/quality/75/11src/product/5966707693/B.jpg?920000000",
			"price":                     "3000",
			"taxRate":                   "10.2",
			"liveShowInventoryQuantity": "2000",
		},
		"fieldComments": []any{
			map[string]any{
				"field": "price",
				"comments": []any{
					map[string]any{
						"id":            float64(1),
						"comment":       "금액 결정되었나요?",
						"createdAt":     "1982-01-05T00:00:00+09:00",
						"createdById":   "1",
						"createdByName": "사이트 관리자",
					},
					map[string]any{
						"id":            float64(2),
						"comment":       "네 250000으로 결정되었네요.",
						"createdAt":     "1982-01-06T00:00:00+09:00",
						"createdById":   "2",
						"createdByName": "김영희",
					},
				},
			},
			map[string]any{
				"field": "mainImage",
				"comments": []any{
					map[string]any{
						"id":            float64(3),
						"comment":       "이미지는 https://gdimg.gmarket.co.kr/2367233519/still/280?ver=1645526559 이것으로 해주세요",
						"createdAt":     "1982-01-07T00:00:00+09:00",
						"createdById":   "3",
						"createdByName": "이수민",
					},
				},
			},
		},
		"fieldChanges": []any{
			map[string]any{
				"field": "name",
				"changes": []any{
					map[string]any{
						"id":            float64(2),
						"beforeValue":   "최신형 공기 살균기",
						"afterValue":    "2024 최신형 공기 살균기",
						"createdAt":     "1982-01-06T00:00:00+09:00",
						"createdById":   "2",
						"createdByName": "김영희",
					},
					map[string]any{
						"id":            float64(1),
						"beforeValue":   "살균기",
						"afterValue":    "최신형 공기 살균기",
						"createdAt":     "1982-01-04T00:00:00+09:00",
						"createdById":   "1",
						"createdByName": "사이트 관리자",
					},
				},
			},
		},
		"createdAt": "1982-01-07T00:00:00+09:00",
		"updatedAt": "1982-01-07T00:00:00+09:00",
	}

	suite.Equal(expected, actual)
}

func (suite *ContentControllerTestSuite) TestGetContent_아이디에_해당하는_데이터가_없으면_NotFound를_반환한다() {
	// given
	sut := testserver.NewTestAppServerBuilder(Router{}, suite.TestDbContainer).DatabaseFixture().Build()

	req := httptest.NewRequest(http.MethodGet, "/api/tenants/bettercode/contents/unknown-id", nil)
	rec := httptest.NewRecorder()

	// when
	sut.ServeHTTP(rec, req)

	// then
	suite.Equal(http.StatusNotFound, rec.Code)
}

func (suite *ContentControllerTestSuite) TestGetContents_페이징() {
	// given
	sut := testserver.NewTestAppServerBuilder(Router{}, suite.TestDbContainer).DatabaseFixture().Build()

	req := httptest.NewRequest(http.MethodGet, "/api/tenants/bettercode/contents?page=1&pageSize=1", nil)
	rec := httptest.NewRecorder()

	// when
	sut.ServeHTTP(rec, req)
	fmt.Println(rec.Body.String())

	// then
	suite.Equal(http.StatusOK, rec.Code)

	var actual any
	json.Unmarshal(rec.Body.Bytes(), &actual)

	expected := map[string]any{
		"result": []any{
			map[string]any{
				"id": "074c7322-e7fa-4d5c-8938-8dbe0ce67465",
				"content": map[string]any{
					"name":      "불스원샷",
					"mainImage": "https://gdimg.gmarket.co.kr/2367233519/still/280?ver=1645526559",
					"price":     "250000",
					"taxRate":   "10.2",
				},
				"createdAt": "1982-01-04T00:00:00+09:00",
				"updatedAt": "1982-01-04T00:00:00+09:00",
			},
		},
		"totalCount": float64(2),
	}
	suite.Equal(expected, actual)
}

func (suite *ContentControllerTestSuite) TestGetContents_Sort_By_Desc_CreatedAt() {
	// given
	sut := testserver.NewTestAppServerBuilder(Router{}, suite.TestDbContainer).DatabaseFixture().Build()

	req := httptest.NewRequest(http.MethodGet, "/api/tenants/bettercode/contents?page=1&pageSize=2&sortBy=desc(created_at)", nil)
	rec := httptest.NewRecorder()

	// when
	sut.ServeHTTP(rec, req)
	fmt.Println(rec.Body.String())

	// then
	suite.Equal(http.StatusOK, rec.Code)

	var actual any
	json.Unmarshal(rec.Body.Bytes(), &actual)

	expected := map[string]any{
		"result": []any{
			map[string]any{
				"id": "6f3bbc99-55aa-4340-89f6-1ddd4dfdb8cd",
				"content": map[string]any{
					"name":                      "2024 최신형 공기 살균기",
					"mainImage":                 "https://cdn.011st.com/11dims/resize/600x600/quality/75/11src/product/5966707693/B.jpg?920000000",
					"price":                     "3000",
					"taxRate":                   "10.2",
					"liveShowInventoryQuantity": "2000",
				},
				"createdAt": "1982-01-07T00:00:00+09:00",
				"updatedAt": "1982-01-07T00:00:00+09:00",
			},
			map[string]any{
				"id": "074c7322-e7fa-4d5c-8938-8dbe0ce67465",
				"content": map[string]any{
					"name":      "불스원샷",
					"mainImage": "https://gdimg.gmarket.co.kr/2367233519/still/280?ver=1645526559",
					"price":     "250000",
					"taxRate":   "10.2",
				},
				"createdAt": "1982-01-04T00:00:00+09:00",
				"updatedAt": "1982-01-04T00:00:00+09:00",
			},
		},
		"totalCount": float64(2),
	}
	suite.Equal(expected, actual)
}

func (suite *ContentControllerTestSuite) TestUpdateContentField() {
	// given
	sut := testserver.NewTestAppServerBuilder(Router{}, suite.TestDbContainer).DatabaseFixture().Build()

	requestBody := `{
			"beforeValue": "불스원샷",
			"afterValue": "불스원샷 플러스",
			"createdById": "1",
			"createdByName": "사이트 관리자"
		}`

	req := httptest.NewRequest(http.MethodPut, "/api/tenants/bettercode/contents/074c7322-e7fa-4d5c-8938-8dbe0ce67465/name", strings.NewReader(requestBody))
	rec := httptest.NewRecorder()

	// when
	sut.ServeHTTP(rec, req)

	// then
	suite.Equal(http.StatusNoContent, rec.Code)
}

func (suite *ContentControllerTestSuite) TestUpdateContentField_Conflict_Field_Value() {
	// given
	sut := testserver.NewTestAppServerBuilder(Router{}, suite.TestDbContainer).DatabaseFixture().Build()

	requestBody := `{
			"beforeValue": "불스원샷 플러스",
			"afterValue": "불스원샷+",
			"createdById": "1",
			"createdByName": "사이트 관리자"
		}`

	req := httptest.NewRequest(http.MethodPut, "/api/tenants/bettercode/contents/074c7322-e7fa-4d5c-8938-8dbe0ce67465/name", strings.NewReader(requestBody))
	rec := httptest.NewRecorder()

	// when
	sut.ServeHTTP(rec, req)

	// then
	suite.Equal(http.StatusConflict, rec.Code)
}

func (suite *ContentControllerTestSuite) TestUpdateContentField_NotFound_Field() {
	// given
	sut := testserver.NewTestAppServerBuilder(Router{}, suite.TestDbContainer).DatabaseFixture().Build()

	requestBody := `{
			"beforeValue": "불스원샷 플러스",
			"afterValue": "불스원샷+",
			"createdById": "1",
			"createdByName": "사이트 관리자"
		}`

	req := httptest.NewRequest(http.MethodPut, "/api/tenants/bettercode/contents/074c7322-e7fa-4d5c-8938-8dbe0ce67465/unknownField", strings.NewReader(requestBody))
	rec := httptest.NewRecorder()

	// when
	sut.ServeHTTP(rec, req)

	// then
	suite.Equal(http.StatusNotFound, rec.Code)
}

func (suite *ContentControllerTestSuite) TestAddFieldComment() {
	// given
	sut := testserver.NewTestAppServerBuilder(Router{}, suite.TestDbContainer).DatabaseFixture().Build()

	requestBody := `{
			"comment": "공기 살균기에 대한 설명",
			"createdById": "1",
			"createdByName": "사이트 관리자"
		}`

	req := httptest.NewRequest(http.MethodPost, "/api/tenants/bettercode/contents/074c7322-e7fa-4d5c-8938-8dbe0ce67465/name/comments", strings.NewReader(requestBody))
	rec := httptest.NewRecorder()

	// when
	sut.ServeHTTP(rec, req)

	// then
	suite.Equal(http.StatusNoContent, rec.Code)
}

func (suite *ContentControllerTestSuite) TestAddFieldComment_NotFound_Field() {
	// given
	sut := testserver.NewTestAppServerBuilder(Router{}, suite.TestDbContainer).DatabaseFixture().Build()
	requestBody := `{
		"comment": "공기 살균기에 대한 설명",
		"createdById": "1",
		"createdByName": "사이트 관리자"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/tenants/bettercode/contents/074c7322-e7fa-4d5c-8938-8dbe0ce67465/unknownField/comments", strings.NewReader(requestBody))
	rec := httptest.NewRecorder()

	// when
	sut.ServeHTTP(rec, req)

	// then
	suite.Equal(http.StatusNotFound, rec.Code)
}
