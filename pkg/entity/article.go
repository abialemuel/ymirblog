package entity

import "github.com/kubuskotak/asgard/rest"

// Result is a resource list result.
type Article struct {
	ID    int
	Title string
	Body  string
	User  *User
	Tags  []Tag
}

type Tag struct {
	ID   int
	Name string
}

type ArticlesWithPagination struct {
	Items    []*Article
	Metadata rest.Pagination
}

type GetArticlePayload struct {
	Title  *string
	UserID *int
	Limit  int
	Page   int
}

type UpsertArticlePayload struct {
	UserID int
	Title  string
	Body   string
	Tags   []Tag
}
