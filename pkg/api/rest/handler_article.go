// Package rest is port handler.
package rest

import (
	"net/http"
	"strconv"

	"github.com/abialemuel/ymirblog/pkg/entity"
	"github.com/abialemuel/ymirblog/pkg/usecase/article"
	"github.com/go-chi/chi/v5"
	"github.com/kubuskotak/asgard/rest"
	"github.com/rs/zerolog/log"
)

// ArticleOption is a struct holding the handler options.
type ArticleOption func(Article *Article)

// Article handler instance data.
type Article struct {
	UcArticle article.T
}

// NewArticle creates a new Article handler instance.
//
//	var ArticleHandler = rest.NewArticle()
//
//	You can pass optional configuration options by passing a Config struct:
//
//	var adaptor = &adapters.Adapter{}
//	var ArticleHandler = rest.NewArticle(rest.WithArticleAdapter(adaptor))
func NewArticle(opts ...ArticleOption) *Article {
	// Create a new handler.
	var handler = &Article{}

	// Assign handler options.
	for o := range opts {
		var opt = opts[o]
		opt(handler)
	}

	// Return handler.
	return handler
}

func WithArticleUsecase(u article.T) ArticleOption {
	return func(a *Article) {
		a.UcArticle = u
	}
}

// Register is endpoint group for handler.
func (a *Article) Register(router chi.Router) {
	// PLEASE EDIT THIS EXAMPLE, how to register handler to router
	router.Get("/articles", rest.HandlerAdapter[GetArticleRequest](a.GetArticle).JSON)
	router.Post("/articles", rest.HandlerAdapter[CreateArticleRequest](a.CreateArticle).JSON)
	router.Get("/articles/{id}", rest.HandlerAdapter[GetOrDeleteArticleByIDRequest](a.GetByIDArticle).JSON)
	router.Delete("/articles/{id}", rest.HandlerAdapter[GetOrDeleteArticleByIDRequest](a.DeleteArticle).JSON)
	router.Patch("/articles/{id}", rest.HandlerAdapter[UpdateArticleRequest](a.UpdateArticle).JSON)

}

// GetArticle endpoint func. /** PLEASE EDIT THIS EXAMPLE, return handler response */.
func (a *Article) GetArticle(w http.ResponseWriter, r *http.Request) (GetArticleResponse, error) {
	var (
		request GetArticleRequest
	)
	request, err := rest.GetBind[GetArticleRequest](r)
	if err != nil {
		return GetArticleResponse{}, rest.ErrBadRequest(w, r, err)
	}

	payload := entity.GetArticlePayload{
		Title:  request.Title,
		UserID: request.UserID,
		Limit:  request.Limit,
		Page:   request.Page,
	}
	articlePagination, err := a.UcArticle.GetAll(r.Context(), payload)
	if err != nil {
		return GetArticleResponse{}, err
	}
	// articlePagination := entity.ArticlesWithPagination{}

	rest.Paging(r, rest.Pagination{
		Page:  articlePagination.Metadata.Page,
		Limit: articlePagination.Metadata.Limit,
		Total: articlePagination.Metadata.Total,
	})

	result := GetArticleResponse{}
	for _, article := range articlePagination.Items {
		articleResponse := ArticleResponse{
			ID:    article.ID,
			Title: article.Title,
			Body:  article.Body,
		}
		if article.User != nil {
			articleResponse.User = &SimpleUserResponse{
				ID:   article.User.ID,
				Name: article.User.Name,
			}
		}
		for _, tag := range article.Tags {
			articleResponse.Tags = append(articleResponse.Tags, SimpleTagResponse{
				ID:   tag.ID,
				Name: tag.Name,
			})
		}
		result.Items = append(result.Items, &articleResponse)
	}

	return result, nil
}

func (a *Article) CreateArticle(w http.ResponseWriter, r *http.Request) (res ResponseArticle, err error) {
	// binding and validate request body
	request, err := rest.GetBind[CreateArticleRequest](r)
	if err != nil {
		return res, rest.ErrBadRequest(w, r, err)
	}

	// mapping request to entity
	// tags
	tags := []entity.Tag{}
	for _, tag := range request.Tags {
		t := entity.Tag{
			Name: tag,
		}

		tags = append(tags, t)
	}
	// user
	// u, err := a.UcUser.GetUserID(r.Context(), request.UserID)
	// if err != nil {
	// 	return res, rest.ErrNotFound(w, r, err)
	// }

	payload := entity.UpsertArticlePayload{
		Title: request.Title,
		Body:  request.Body,
		Tags:  tags,
		// User:  &u,
	}

	//create entity with usecase
	ent, err := a.UcArticle.Create(r.Context(), payload)
	if err != nil {
		log.Error().Err(err).Msg("CreateArticle4")
		return res, rest.ErrBadRequest(w, r, err)
	}

	// mapping entity to response
	res.Message = "success create article"
	res.Article = &ent

	return res, nil
}

func (a *Article) DeleteArticle(w http.ResponseWriter, r *http.Request) (res ResponseArticle, err error) {
	request, err := rest.GetBind[GetOrDeleteArticleByIDRequest](r)
	if err != nil {
		return res, rest.ErrBadRequest(w, r, err)
	}

	err = a.UcArticle.Delete(r.Context(), request.ID)
	if err != nil {
		return res, rest.ErrBadRequest(w, r, err)
	}

	res.Message = "success delete article"

	return res, nil
}

func (a *Article) UpdateArticle(w http.ResponseWriter, r *http.Request) (res ResponseArticle, err error) {
	// req := RequestUpdateArticle{}
	req, err := rest.GetBind[UpdateArticleRequest](r)
	if err != nil {
		return res, rest.ErrBadRequest(w, r, err)
	}

	id, _ := strconv.Atoi(req.ID)

	// get user
	// u, err := a.UcUser.GetUserID(r.Context(), req.UserID)
	// if err != nil {
	// 	return res, rest.ErrNotFound(w, r, err)
	// }

	tags := []entity.Tag{}
	for _, tag := range req.Tags {
		t := entity.Tag{
			Name: tag,
		}

		tags = append(tags, t)
	}

	payload := entity.UpsertArticlePayload{
		// ID:    id,
		Title: req.Title,
		Body:  req.Body,
		// User:  &u,
		Tags: tags,
	}

	article, err := a.UcArticle.Update(r.Context(), id, payload)
	if err != nil {
		return res, rest.ErrBadRequest(w, r, err)
	}

	// mapping request to entity
	res.Message = "success update article"
	res.Article = &article

	return res, nil
}

func (a *Article) GetByIDArticle(w http.ResponseWriter, r *http.Request) (res ResponseArticle, err error) {
	request, err := rest.GetBind[GetOrDeleteArticleByIDRequest](r)
	if err != nil {
		return res, rest.ErrBadRequest(w, r, err)
	}

	e, err := a.UcArticle.GetByID(r.Context(), request.ID)
	if err != nil {
		return res, rest.ErrNotFound(w, r, err)
	}

	// mapping entity to response
	res.Message = "success get article"
	res.Article = &e

	return res, nil
}
