// Package rest is port handler.
package rest

import "github.com/abialemuel/ymirblog/pkg/entity"

// GetArticleRequest Get a Article request.  /** PLEASE EDIT THIS EXAMPLE, request handler */.
type GetArticleRequest struct {
	Title  *string
	UserID *int
	Limit  int `validate:"gte=0,default=10"`
	Page   int `validate:"gte=0,default=1"`
}

// GetArticleResponse Get a Article response.  /** PLEASE EDIT THIS EXAMPLE, return handler response */.
type GetArticleResponse struct {
	Items []*ArticleResponse `json:"items"`
}

type ArticleResponse struct {
	ID    int                 `json:"id"`
	Title string              `json:"title"`
	Body  string              `json:"body"`
	User  *SimpleUserResponse `json:"user,omitempty"`
	Tags  []SimpleTagResponse `json:"tags"`
}

type SimpleUserResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type SimpleTagResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type GetOrDeleteArticleByIDRequest struct {
	ID int `json:"id" validate:"required"`
}

type UpdateArticleRequest struct {
	// from path parameter
	ID string `json:"id" validate:"required"`
	// from body
	UserID int      `json:"user_id" validate:"required"`
	Title  string   `json:"title" validate:"required"`
	Body   string   `json:"body" validate:"required"`
	Tags   []string `json:"tags,omitempty" validate:"required"`
}

type CreateArticleRequest struct {
	UserID int      `json:"user_id" validate:"required"`
	Title  string   `json:"title" validate:"required"`
	Body   string   `json:"body" validate:"required"`
	Tags   []string `json:"tags,omitempty" validate:"required"`
}

type ResponseArticle struct {
	Message  string           `json:"message"`
	Article  *entity.Article  `json:"article,omitempty"`
	Articles []entity.Article `json:"articles,omitempty"`
}
