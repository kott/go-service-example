package articles

import (
	"context"
)

// Repo defines the DB level interaction of articles
type Repo interface {
	Get(ctx context.Context, id string) (Article, error)
	GetAll(ctx context.Context, limit, offset int) ([]Article, error)
	Create(ctx context.Context, ar ArticleCreateUpdate) (string, error)
	Update(ctx context.Context, ar ArticleCreateUpdate, id string) error
}

// Service defines the service level contract that other services
// outside this package can use to interact with Article resources
type Service interface {
	Get(ctx context.Context, id string) (Article, error)
	GetAll(ctx context.Context, limit, offset int) ([]Article, error)
	Create(ctx context.Context, ar ArticleCreateUpdate) (Article, error)
	Update(ctx context.Context, ar ArticleCreateUpdate, id string) (Article, error)
}

type article struct {
	repo Repo
}

// New Service instance
func New(repo Repo) Service {
	return &article{repo}
}

// Get sends the request straight to the repo
func (s *article) Get(ctx context.Context, id string) (Article, error) {
	return s.repo.Get(ctx, id)
}

// GetAll sends the request straight to the repo
func (s *article) GetAll(ctx context.Context, limit, offset int) ([]Article, error) {
	return s.repo.GetAll(ctx, limit, offset)
}

// Create passes of the created to the repo and retrieves the newly created record
func (s *article) Create(ctx context.Context, ar ArticleCreateUpdate) (Article, error) {
	id, err := s.repo.Create(ctx, ar)
	if err != nil {
		return Article{}, err
	}
	return s.repo.Get(ctx, id)
}

// Update the requested resource
func (s *article) Update(ctx context.Context, ar ArticleCreateUpdate, id string) (Article, error) {
	if err := s.repo.Update(ctx, ar, id); err != nil {
		return Article{}, err
	}
	return s.Get(ctx, id)
}
