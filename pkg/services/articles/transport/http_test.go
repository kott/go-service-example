package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kott/go-service-example/pkg/errors"
	"github.com/kott/go-service-example/pkg/services/articles"
)

type mockService struct {
	GetResult    articles.Article
	GetErr       error
	GetAllResult []articles.Article
	GetAllErr    error
	CreateResult articles.Article
	CreateErr    error
	UpdateResult articles.Article
	UpdateErr    error
}

func (s *mockService) Get(ctx context.Context, id string) (articles.Article, error) {
	return s.GetResult, s.GetErr
}

func (s *mockService) GetAll(ctx context.Context, limit, offset int) ([]articles.Article, error) {
	return s.GetAllResult, s.GetAllErr
}

func (s *mockService) Create(ctx context.Context, ar articles.ArticleCreateUpdate) (articles.Article, error) {
	return s.CreateResult, s.CreateErr
}

func (s *mockService) Update(ctx context.Context, ar articles.ArticleCreateUpdate, id string) (articles.Article, error) {
	return s.UpdateResult, s.UpdateErr
}

func TestHandlerGet(t *testing.T) {
	id := uuid.New().String()
	tests := map[string]struct {
		mockService articles.Service
		uri         string
		response    interface{}
		status      int
	}{
		"Happy path": {
			mockService: &mockService{
				GetResult: articles.Article{ID: id},
				GetErr:    nil,
			},
			uri: fmt.Sprintf("/articles/%s", id),
			response: articles.Article{
				ID:         id,
				Title:      "",
				Body:       "",
				CreatedAt:  time.Time{},
				UpdatedAt:  time.Time{},
				DisabledAt: nil,
			},
			status: http.StatusOK,
		},
		"Not found": {
			mockService: &mockService{
				GetResult: articles.Article{},
				GetErr:    articles.ErrArticleNotFound,
			},
			uri: fmt.Sprintf("/articles/%s", id),
			response: errors.AppError{
				Code:        errors.NotFound,
				Description: articles.ErrArticleNotFound.Error(),
				Field:       "id",
			},
			status: http.StatusNotFound,
		},
		"Server error": {
			mockService: &mockService{
				GetResult: articles.Article{},
				GetErr:    fmt.Errorf("internal"),
			},
			uri: fmt.Sprintf("/articles/%s", id),
			response: errors.AppError{
				Code:        errors.InternalServerError,
				Description: "internal",
				Field:       "unknown",
			},
			status: http.StatusInternalServerError,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			response := httptest.NewRecorder()
			router := gin.New()
			newHandler(router, test.mockService)

			req, err := http.NewRequest(http.MethodGet, test.uri, nil)
			require.NoError(t, err)
			req.Header.Add("Content-Type", "application/json")

			router.ServeHTTP(response, req)

			assert.Equal(t, test.status, response.Code)

			if test.status == http.StatusOK {
				var ar articles.Article
				if err := json.Unmarshal(response.Body.Bytes(), &ar); err != nil {
					assert.Fail(t, "failed to unmarshal", response.Body.String(), err)
				}
				assert.Equal(t, test.response, ar)
			} else {
				var err errors.AppError
				if err := json.Unmarshal(response.Body.Bytes(), &err); err != nil {
					assert.Fail(t, "failed to unmarshal", response.Body.String(), err)
				}
				assert.Equal(t, test.response, err)
			}
		})
	}
}

func TestHandlerCreate(t *testing.T) {
	id := uuid.New().String()
	title := "some-title"
	body := "some-body"

	ar := articles.Article{
		ID:         id,
		Title:      title,
		Body:       body,
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
		DisabledAt: nil,
	}

	tests := map[string]struct {
		mockService articles.Service
		uri         string
		body        string
		response    interface{}
		status      int
	}{
		"Happy path": {
			mockService: &mockService{
				CreateResult: ar,
				CreateErr:    nil,
			},
			uri:  fmt.Sprintf("/articles/"),
			body: fmt.Sprintf(`{"title": "%s", "body": "%s"}`, title, body),
			response: articles.Article{
				ID:         id,
				Title:      title,
				Body:       body,
				CreatedAt:  time.Time{},
				UpdatedAt:  time.Time{},
				DisabledAt: nil,
			},
			status: http.StatusCreated,
		},
		"Create failure": {
			mockService: &mockService{
				CreateResult: articles.Article{},
				CreateErr:    articles.ErrArticleCreate,
			},
			uri:  fmt.Sprintf("/articles/"),
			body: fmt.Sprintf(`{"title": "%s", "body": "%s"}`, title, body),
			response: errors.AppError{
				Code:        errors.InternalServerError,
				Description: "unable to create/update article",
				Field:       "",
			},
			status: http.StatusInternalServerError,
		},
		"Malformed request": {
			mockService: &mockService{
				CreateResult: articles.Article{},
				CreateErr:    articles.ErrArticleCreate,
			},
			uri:  fmt.Sprintf("/articles/"),
			body: `{}`,
			response: errors.AppError{
				Code:        errors.BadRequest,
				Description: errors.Descriptions[errors.BadRequest],
				Field:       "",
			},
			status: http.StatusBadRequest,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			response := httptest.NewRecorder()
			router := gin.New()
			newHandler(router, test.mockService)

			req, err := http.NewRequest(http.MethodPost, test.uri, strings.NewReader(test.body))
			require.NoError(t, err)
			req.Header.Add("Content-Type", "application/json")

			router.ServeHTTP(response, req)

			assert.Equal(t, test.status, response.Code)

			if test.status == http.StatusCreated {
				var ar articles.Article
				if err := json.Unmarshal(response.Body.Bytes(), &ar); err != nil {
					assert.Fail(t, "failed to unmarshal", response.Body.String(), err)
				}
				assert.Equal(t, test.response, ar)
			} else {
				var err errors.AppError
				if err := json.Unmarshal(response.Body.Bytes(), &err); err != nil {
					assert.Fail(t, "failed to unmarshal", response.Body.String(), err)
				}
				assert.Equal(t, test.response, err)
			}
		})
	}
}
