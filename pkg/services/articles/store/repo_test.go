package store

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/kott/go-service-example/pkg/services/articles"
)

func TestArticleRepoGet(t *testing.T) {
	columns := []string{"id", "title", "body", "created_at", "updated_at", "disabled_at"}
	id := uuid.New().String()
	now := time.Now()
	mockResult := []driver.Value{id, "title", "body", now, now, now}

	tests := map[string]struct {
		expectQueryArgs        []driver.Value
		expectQueryResultRows  []*sqlmock.Rows
		expectQueryResultError error
		input                  string
		expect                 articles.Article
		err                    error
	}{
		"Happy path": {
			expectQueryArgs:        []driver.Value{id},
			expectQueryResultRows:  []*sqlmock.Rows{sqlmock.NewRows(columns).AddRow(mockResult...)},
			expectQueryResultError: nil,
			input:                  id,
			expect:                 articles.Article{ID: id},
			err:                    nil,
		},
		"Unknown DB error": {
			expectQueryArgs:        []driver.Value{id},
			expectQueryResultRows:  []*sqlmock.Rows{sqlmock.NewRows(columns).AddRow(mockResult...)},
			expectQueryResultError: errors.New("some-db-error"),
			input:                  id,
			expect:                 articles.Article{},
			err:                    articles.ErrArticleNotFound,
		},
		"Not found error": {
			expectQueryArgs:        []driver.Value{"fake-id"},
			expectQueryResultRows:  []*sqlmock.Rows{sqlmock.NewRows(columns).AddRow(mockResult...)},
			expectQueryResultError: sql.ErrNoRows,
			input:                  "fake-id",
			expect:                 articles.Article{},
			err:                    articles.ErrArticleNotFound,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()

			mock.ExpectQuery(regexp.QuoteMeta(selectArticle)).WithArgs(test.expectQueryArgs...).WillReturnError(test.expectQueryResultError).WillReturnRows(test.expectQueryResultRows...)

			repo := New(db)
			response, err := repo.Get(context.Background(), test.input)

			assert.Equal(t, test.err, err)
			assert.Equal(t, test.expect.ID, response.ID)
		})
	}
}

func TestArticleRepoCreate(t *testing.T) {
	columns := []string{"id"}
	id := uuid.New().String()
	title := "some-title"
	body := "some-body"

	tests := map[string]struct {
		expectQueryArgs        []driver.Value
		expectQueryResultRows  []*sqlmock.Rows
		expectQueryResultError error
		input                  articles.ArticleCreateUpdate
		expect                 string
		err                    error
	}{
		"Happy path": {
			expectQueryArgs:        []driver.Value{title, body},
			expectQueryResultRows:  []*sqlmock.Rows{sqlmock.NewRows(columns).AddRow(id)},
			expectQueryResultError: nil,
			input:                  articles.ArticleCreateUpdate{Title: title, Body: body},
			expect:                 id,
			err:                    nil,
		},
		"Create error": {
			expectQueryArgs:        []driver.Value{title, body},
			expectQueryResultRows:  []*sqlmock.Rows{sqlmock.NewRows(columns).AddRow(id)},
			expectQueryResultError: errors.New("some-db-error"),
			input:                  articles.ArticleCreateUpdate{Title: title, Body: body},
			expect:                 "",
			err:                    articles.ErrArticleCreate,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()

			mock.ExpectQuery(regexp.QuoteMeta(insertArticle)).WithArgs(test.expectQueryArgs...).WillReturnError(test.expectQueryResultError).WillReturnRows(test.expectQueryResultRows...)

			repo := New(db)
			response, err := repo.Create(context.Background(), test.input)

			assert.Equal(t, test.err, err)
			assert.Equal(t, test.expect, response)
		})
	}
}
