package store

import (
	"context"
	"database/sql"

	"github.com/kott/go-service-example/pkg/services/articles"
	"github.com/kott/go-service-example/pkg/utils/log"
)

const (
	selectArticle      = `SELECT * FROM articles WHERE id=$1`
	selectManyArticles = `SELECT * FROM articles LIMIT $1 OFFSET $2`
	insertArticle      = `INSERT INTO articles (title, body, created_at, updated_at) VALUES ($1, $2, now(), now()) RETURNING id`
	updateArticle      = `UPDATE articles SET title = $1, body = $2, updated_at = now() WHERE id = $3`
)

type articleRepo struct {
	DB *sql.DB
}

// New creates an instance of the accountRepo.
func New(conn *sql.DB) articles.Repo {
	return &articleRepo{conn}
}

// Get retrieves the article with the given id
func (r *articleRepo) Get(ctx context.Context, id string) (articles.Article, error) {
	var ar articles.Article

	err := r.DB.QueryRow(selectArticle, id).
		Scan(&ar.ID, &ar.Title, &ar.Body, &ar.CreatedAt, &ar.UpdatedAt, &ar.DisabledAt)
	if err != nil {
		log.Info(ctx, "select article error: %s", err.Error())
		return ar, articles.ErrArticleNotFound
	}

	return ar, nil
}

// GetAll retrieves all articles within the limit and offset. Limit defaults to 25
func (r *articleRepo) GetAll(ctx context.Context, limit, offset int) ([]articles.Article, error) {
	al := make([]articles.Article, 0)

	rows, err := r.DB.Query(selectManyArticles, limit, offset)
	if err != nil {
		log.Warn(ctx, "unable to query db: %s", err.Error())
		return al, articles.ErrArticleQuery
	}
	defer rows.Close()

	for rows.Next() {
		var ar articles.Article
		if err := rows.Scan(&ar.ID, &ar.Title, &ar.Body, &ar.CreatedAt, &ar.UpdatedAt, &ar.DisabledAt); err != nil {
			log.Error(ctx, "unable to scan db rows: %s", err.Error())
			return al, articles.ErrArticleQuery
		}

		al = append(al, ar)
	}

	return al, nil
}

// Create sets the title and body in a new db record
func (r *articleRepo) Create(ctx context.Context, ar articles.ArticleCreateUpdate) (string, error) {
	var id string
	if err := r.DB.QueryRow(insertArticle, ar.Title, ar.Body).Scan(&id); err != nil {
		log.Error(ctx, "unable to create article: %s", err.Error())
		return "", articles.ErrArticleCreate
	}

	log.Info(ctx, "created article with id=%s", id)
	return id, nil
}

// Update sets the title and body on an existing record on the requested version of the row
func (r *articleRepo) Update(ctx context.Context, ar articles.ArticleCreateUpdate, id string) error {
	_, err := r.DB.Exec(updateArticle, ar.Title, ar.Body, id)
	if err != nil {
		log.Error(ctx, "unable to update article (%s): %s", id, err.Error())
		return articles.ErrArticleUpdate
	}
	return nil
}
