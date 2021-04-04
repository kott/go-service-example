package transport

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kott/go-service-example/pkg/errors"
	"github.com/kott/go-service-example/pkg/services/articles"
	"github.com/kott/go-service-example/pkg/services/articles/store"
	"github.com/kott/go-service-example/pkg/utils/context"
	"github.com/kott/go-service-example/pkg/utils/log"
)

type handler struct {
	ArticleService articles.Service
}

// Activate sets all the services required for articles and registers all the endpoints with the engine.
func Activate(router *gin.Engine, db *sql.DB) {
	articleService := articles.New(store.New(db))
	newHandler(router, articleService)
}

func newHandler(router *gin.Engine, as articles.Service) {
	h := handler{
		ArticleService: as,
	}

	router.GET("/articles/:id", h.Get)
	router.GET("/articles/", h.GetAll)
	router.POST("/articles/", h.Create)
	router.PUT("/articles/:id", h.Update)
}

func (h *handler) Get(c *gin.Context) {
	ctx := context.GetReqCtx(c)

	log.Info(ctx, "retrieving article id=%s", c.Param("id"))
	article, err := h.ArticleService.Get(ctx, c.Param("id"))
	if err != nil {
		status, appErr := handleError(err)
		c.IndentedJSON(status, appErr)
		return
	}

	c.IndentedJSON(http.StatusOK, article)
}

func (h *handler) GetAll(c *gin.Context) {
	var q struct {
		Limit  int `form:"limit,default=25"`
		Offset int `form:"offset,default=0"`
	}

	ctx := context.GetReqCtx(c)
	if err := c.BindQuery(&q); err != nil {
		log.Info(ctx, "query parse error: %s", err.Error())
		c.IndentedJSON(http.StatusBadRequest,
			errors.NewAppError(errors.BadRequest, errors.Descriptions[errors.BadRequest], ""))
		return
	}

	log.Info(ctx, "retrieving all articles: offset=%d limit=%d", q.Limit, q.Offset)
	artcls, err := h.ArticleService.GetAll(ctx, q.Limit, q.Offset)
	if err != nil {
		status, appErr := handleError(err)
		c.IndentedJSON(status, appErr)
		return
	}

	c.IndentedJSON(http.StatusOK, articles.Articles{Articles: artcls})
}

func (h *handler) Create(c *gin.Context) {
	ctx := context.GetReqCtx(c)

	var ac articles.ArticleCreateUpdate
	if err := c.ShouldBindJSON(&ac); err != nil {
		log.Info(ctx, "request parse error: %s", err.Error())
		c.IndentedJSON(http.StatusBadRequest,
			errors.NewAppError(errors.BadRequest, errors.Descriptions[errors.BadRequest], ""))
		return
	}

	log.Info(ctx, "creating article %v", ac)
	article, err := h.ArticleService.Create(ctx, ac)
	if err != nil {
		status, appErr := handleError(err)
		c.IndentedJSON(status, appErr)
		return
	}

	c.IndentedJSON(http.StatusCreated, article)
}

func (h *handler) Update(c *gin.Context) {
	ctx := context.GetReqCtx(c)

	var ac articles.ArticleCreateUpdate
	if err := c.ShouldBindJSON(&ac); err != nil {
		log.Info(ctx, "request parse error: %s", err.Error())
		c.IndentedJSON(http.StatusBadRequest,
			errors.NewAppError(errors.BadRequest, errors.Descriptions[errors.BadRequest], ""))
		return
	}

	id := c.Param("id")

	log.Info(ctx, "updating article %s with %v", id, ac)
	article, err := h.ArticleService.Update(ctx, ac, id)
	if err != nil {
		status, appErr := handleError(err)
		c.IndentedJSON(status, appErr)
		return
	}

	c.IndentedJSON(http.StatusOK, article)
}

// handleError allows us to map errors defined internally to appropriate HTTP error codes and JSON responses
func handleError(e error) (int, error) {
	switch e {
	case articles.ErrArticleNotFound:
		return http.StatusNotFound, errors.NewAppError(errors.NotFound, e.Error(), "id")
	case articles.ErrArticleUpdate:
		fallthrough
	case articles.ErrArticleCreate:
		return http.StatusInternalServerError, errors.NewAppError(errors.InternalServerError, "unable to create/update article", "")
	default:
		return http.StatusInternalServerError, errors.NewAppError(errors.InternalServerError, e.Error(), "unknown")
	}
}
