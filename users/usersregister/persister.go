package usersregister

import (
	"context"
	"database/sql"

	"github.com/diegoholiveira/bookstore-sample/users"
)

type (
	DatabaseExecutor interface {
		ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	}

	persister struct {
		db DatabaseExecutor
	}
)

func NewDatabasePersister(db DatabaseExecutor) DatabasePersister {
	return &persister{
		db: db,
	}
}

func (p persister) Persist(ctx context.Context, u *users.User) error {
	result, err := p.db.ExecContext(
		ctx,
		`INSERT INTO users (name, email, password) VALUES (?, ?, ?)`,
		u.Name,
		u.Email,
		u.Password,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	u.ID = uint64(id)

	return nil
}
