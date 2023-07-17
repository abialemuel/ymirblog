// Package article is implements component logic.
package article

import (
	"context"
	"fmt"
	"reflect"

	"entgo.io/ent/dialect"
	"github.com/abialemuel/ymirblog/pkg/adapters"
	"github.com/abialemuel/ymirblog/pkg/entity"
	"github.com/abialemuel/ymirblog/pkg/persist/ymirblog"
	"github.com/abialemuel/ymirblog/pkg/usecase"
)

func init() {
	usecase.Register(usecase.Registration{
		Name: "article",
		Inf:  reflect.TypeOf((*T)(nil)).Elem(),
		New: func() any {
			return &impl{}
		},
	})
}

// T is the interface implemented by all article Component implementations.
type T interface {
	GetAll(ctx context.Context, request entity.GetArticlePayload) (entity.ArticlesWithPagination, error)
	GetByID(ctx context.Context, id int) (entity.Article, error)
	Create(ctx context.Context, p entity.UpsertArticlePayload) (entity.Article, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, id int, p entity.UpsertArticlePayload) (entity.Article, error)
}

type impl struct {
	adapter *adapters.Adapter
}

// Init initializes the execution of a process involved in a article Component usecase.
func (i *impl) Init(adapter *adapters.Adapter) error {
	i.adapter = adapter
	return nil
}

func WithYmirBlogPersist() adapters.Option {
	return func(a *adapters.Adapter) {
		// adapter  MySQL
		if a.YmirblogMySQL == nil {
			panic(fmt.Errorf("%s is not found", "YmirBlogMySQL"))
		}
		// persist ymirblog driver
		var c = ymirblog.Driver(
			ymirblog.WithDriver(a.YmirblogMySQL, dialect.MySQL),
		)

		a.YmirblogPersist = c
	}
}
