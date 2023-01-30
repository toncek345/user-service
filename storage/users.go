package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

var _ UserStorage = (*UserStorageSQL)(nil)
var _ UserStorage = (*MockUser)(nil)

type UserStorage interface {
	InsertUser(ctx context.Context, user *InsertUser) (*UserModel, error)
	DeleteUser(ctx context.Context, id string) error
	UpdateUser(ctx context.Context, user *UpdateUser) (*UserModel, error)
	SearchUser(ctx context.Context, filters *Filters, offset, limit int64) ([]*UserModel, error)
}

// ErrNotFound is returned as an error if object doesn't exist in DB.
var ErrNotFound = errors.New("object doesn't exist in db")

type UserStorageSQL struct {
	DB *sqlx.DB
}

type InsertUser struct {
	FirstName string
	LastName  string
	Email     string
	Country   string
	Password  string
}

func (us *UserStorageSQL) InsertUser(ctx context.Context, user *InsertUser) (*UserModel, error) {
	u := &UserModel{}

	tx, err := us.DB.Beginx()
	if err != nil {
		return nil, fmt.Errorf("starting transaction: %w", err)
	}

	if err := tx.GetContext(
		ctx,
		u,
		`INSERT INTO users (id, first_name, last_name, email, country, password, created_at, updated_at) VALUES
		(uuid_generate_v4(), $1, $2, $3, $4, $5, NOW(), NOW()) RETURNING
		id, first_name, last_name, email, country, password, created_at, updated_at`,
		user.FirstName, user.LastName, user.Email, user.Country, user.Password); err != nil {

		tx.Rollback()
		return nil, fmt.Errorf("inserting user: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("transaction commit: %w", err)
	}

	return u, nil
}

type UpdateUser struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
	Country   string
	Password  string
}

func (us *UserStorageSQL) UpdateUser(ctx context.Context, user *UpdateUser) (*UserModel, error) {
	u := &UserModel{}

	tx, err := us.DB.Beginx()
	if err != nil {
		return nil, fmt.Errorf("starting transaction: %w", err)
	}

	if err := tx.GetContext(
		ctx,
		u,
		`UPDATE users SET first_name = $1, last_name = $2, email = $3, country = $4, password = $5, updated_at = NOW()
		WHERE users.id = $6 RETURNING
		id, first_name, last_name, email, country, password, updated_at,
		(SELECT created_at FROM users WHERE id = $6) AS created_at`,
		user.FirstName, user.LastName, user.Email, user.Country, user.Password, user.ID); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		tx.Rollback()
		return nil, fmt.Errorf("updating user: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("transaction commit: %w", err)
	}

	return u, nil
}

type UserModel struct {
	// ID is represented in UUID.
	ID        string    `db:"id"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Email     string    `db:"email"`
	Country   string    `db:"country"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (us *UserStorageSQL) DeleteUser(ctx context.Context, id string) error {
	tx, err := us.DB.Beginx()
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}

	if _, err := tx.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id); err != nil {
		return fmt.Errorf("sql deleting: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit: %w", err)
	}

	return nil
}

type Filters struct {
	Country string
}

func (us *UserStorageSQL) SearchUser(ctx context.Context, filters *Filters, offset, limit int64) ([]*UserModel, error) {
	query := sq.Select("*").From("users").Offset(uint64(offset)).Limit(uint64(limit))

	if filters.Country != "" {
		query = query.Where("country ILIKE $1", filters.Country)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building sql: %w", err)
	}

	users := []*UserModel{}
	if err := us.DB.SelectContext(ctx, &users, sql, args...); err != nil {
		return nil, fmt.Errorf("user searching: %w", err)
	}

	return users, nil
}
