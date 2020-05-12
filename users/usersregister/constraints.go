package usersregister

import (
	"context"
	"database/sql"
)

type (
	DatabaseQuerier interface {
		QueryRowContext(context.Context, string, ...interface{}) *sql.Row
	}

	dbConstraints struct {
		db DatabaseQuerier
	}
)

func NewUserConstraints(db DatabaseQuerier) UserConstraints {
	return dbConstraints{
		db: db,
	}
}

func (c dbConstraints) EmailIsUsed(ctx context.Context, email string) bool {
	var (
		err error
		ok  uint32
	)
	err = c.db.QueryRowContext(
		ctx,
		"SELECT 1 FROM users WHERE email = ?",
		email,
	).Scan(&ok)
	if err == sql.ErrNoRows {
		return false
	}
	return true
}
