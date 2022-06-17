package queries

import "github.com/volatiletech/null/v9"

type Pagination struct {
	Page  int `query:"page" validate:"int" default:"1"`
	Limit int `query:"limit" validate:"int" default:"100"`
}

type Period struct {
	TimeFrom null.Int64 `query:"time_from" validate:"int"`
	TimeTo   null.Int64 `query:"time_to" validate:"int"`
}

type Ordering string

var (
	OrderingAsc  Ordering = "asc"
	OrderingDesc Ordering = "desc"
)

type Order struct {
	OrderBy  string   `query:"order_by" default:"created_at"`
	Ordering Ordering `query:"ordering" default:"asc"`
}
