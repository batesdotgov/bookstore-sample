package usersregister

import (
	"context"
	"database/sql"

	"github.com/diegoholiveira/bookstore-sample/users"
)

type persister struct {
	db *sql.DB
}

func NewUserPersister(db *sql.DB) UserPersister {
	return &persister{
		db: db,
	}
}

func (p persister) Persist(ctx context.Context, u users.User) error {
	_, err := p.db.ExecContext(
		ctx,
		`INSERT INTO users (name, email, password) VALUES (?, ?, ?)`,
		u.Name,
		u.Email,
		u.Password,
	)

	return err
}
