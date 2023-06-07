// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: query.sql

package querytest

import (
	"context"
	"database/sql"
)

const test = `-- name: Test :one
select txt from Demo
where txt ~~ '%' || $1 || '%'
`

func (q *Queries) Test(ctx context.Context, val string) (string, error) {
	row := q.db.QueryRowContext(ctx, test, val)
	var txt string
	err := row.Scan(&txt)
	return txt, err
}

const test2 = `-- name: Test2 :one
select txt from Demo
where txt like '%' || $1 || '%'
`

func (q *Queries) Test2(ctx context.Context, val sql.NullString) (string, error) {
	row := q.db.QueryRowContext(ctx, test2, val)
	var txt string
	err := row.Scan(&txt)
	return txt, err
}

const test3 = `-- name: Test3 :one
select txt from Demo
where txt like concat('%', $1, '%')
`

func (q *Queries) Test3(ctx context.Context, val interface{}) (string, error) {
	row := q.db.QueryRowContext(ctx, test3, val)
	var txt string
	err := row.Scan(&txt)
	return txt, err
}
