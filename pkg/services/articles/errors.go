package articles

import (
	"errors"
)

var (
	// ErrArticleNotFound ...
	ErrArticleNotFound = errors.New("requested article could not be found")

	// ErrArticleQuery ...
	ErrArticleQuery = errors.New("requested articles could not be retrieved base on the given criteria")

	// ErrArticleCreate ...
	ErrArticleCreate = errors.New("article could not be created")

	// ErrArticleUpdate ...
	ErrArticleUpdate = errors.New("article could not be updated")
)
