package articles

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type repoMock struct {
	GetResult Article
	GetError  error

	GetAllResult []Article
	GetAllError  error

	CreateResult string
	CreateError  error

	UpdateError error
}

func (r *repoMock) Get(ctx context.Context, id string) (Article, error) {
	return r.GetResult, r.GetError
}

func (r *repoMock) GetAll(ctx context.Context, limit, offset int) ([]Article, error) {
	return r.GetAllResult, r.GetAllError
}

func (r *repoMock) Create(ctx context.Context, ar ArticleCreateUpdate) (string, error) {
	return r.CreateResult, r.CreateError
}

func (r *repoMock) Update(ctx context.Context, ar ArticleCreateUpdate, id string) error {
	return r.UpdateError
}

func TestServiceGet(t *testing.T) {
	id := uuid.New().String()
	tests := map[string]struct {
		repo   Repo
		result Article
		err    error
	}{
		"Happy path": {
			repo: &repoMock{
				GetResult: Article{ID: id},
				GetError:  nil,
			},
			result: Article{ID: id},
			err:    nil,
		},
		"Not found from repo": {
			repo: &repoMock{
				GetResult: Article{},
				GetError:  ErrArticleNotFound,
			},
			result: Article{},
			err:    ErrArticleNotFound,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			service := New(test.repo)
			response, err := service.Get(context.Background(), id)

			assert.Equal(t, test.err, err)
			assert.Equal(t, test.result.ID, response.ID)
		})
	}
}

func TestServiceCreate(t *testing.T) {
	id := uuid.New().String()
	title := "some-title"
	body := "some-body"
	ac := ArticleCreateUpdate{Title: title, Body: body}
	ar := Article{
		ID:         id,
		Title:      title,
		Body:       body,
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
		DisabledAt: nil,
	}

	tests := map[string]struct {
		repo   Repo
		result Article
		input  ArticleCreateUpdate
		err    error
	}{
		"Happy path": {
			repo: &repoMock{
				CreateResult: id,
				CreateError:  nil,
				GetResult:    ar,
				GetError:     nil,
			},
			input:  ac,
			result: ar,
			err:    nil,
		},
		"Not found from repo after create": {
			repo: &repoMock{
				CreateResult: id,
				CreateError:  nil,
				GetResult:    Article{},
				GetError:     ErrArticleNotFound,
			},
			input:  ac,
			result: Article{},
			err:    ErrArticleNotFound,
		},
		"Creation failure": {
			repo: &repoMock{
				CreateResult: id,
				CreateError:  ErrArticleCreate,
			},
			input:  ac,
			result: Article{},
			err:    ErrArticleCreate,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			service := New(test.repo)
			response, err := service.Create(context.Background(), test.input)

			assert.Equal(t, test.err, err)
			assert.Equal(t, test.result, response)
		})
	}
}
