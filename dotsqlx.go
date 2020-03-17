package dotsqlx

import (
	"context"
	"database/sql"

	"github.com/gchaincl/dotsql"
	"github.com/jmoiron/sqlx"
)

type Preparerx interface {
	Preparex(query string) (*sqlx.Stmt, error)
}

type PreparerxContext interface {
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
}

type Getter interface {
	Get(dest interface{}, query string, args ...interface{}) error
}

type GetterContext interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type Selecter interface {
	Select(dest interface{}, query string, args ...interface{}) error
}

type SelecterContext interface {
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type QueryRowerx interface {
	QueryRowx(query string, args ...interface{}) *sqlx.Row
}

type QueryRowerxContext interface {
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
}

type Queryerx interface {
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
}

type QueryerxContext interface {
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
}

type MustExecer interface {
	MustExec(query string, args ...interface{}) sql.Result
}

type MustExecerContext interface {
	MustExecContext(ctx context.Context, query string, args ...interface{}) sql.Result
}

type Rebinder interface {
	Rebind(query string) string
}

type NamedPreparer interface {
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
}

type NamedPreparerContext interface {
	PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)
}

type NamedQueryer interface {
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
}

type NamedQueryerContext interface {
	NamedQueryContext(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error)
}

type NamedExecer interface {
	NamedExec(query string, arg interface{}) (sql.Result, error)
}

type NamedExecerContext interface {
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}

type NamedBinder interface {
	BindNamed(query string, arg interface{}) (string, []interface{}, error)
}

type DotSqlx struct {
	dotsql.DotSql
}

func (d DotSqlx) Get(dbx Getter, dest interface{}, name string, args ...interface{}) error {
	query, err := d.Raw(name)
	if err != nil {
		return err
	}

	return dbx.Get(dest, query, args...)
}

func (d DotSqlx) GetContext(ctx context.Context, dbx GetterContext, dest interface{}, name string, args ...interface{}) error {
	query, err := d.Raw(name)
	if err != nil {
		return err
	}

	return dbx.GetContext(ctx, dest, query, args...)
}

func (d DotSqlx) Select(dbx Selecter, dest interface{}, name string, args ...interface{}) error {
	query, err := d.Raw(name)
	if err != nil {
		return err
	}

	return dbx.Select(dest, query, args...)
}

func (d DotSqlx) SelectContext(ctx context.Context, dbx SelecterContext, dest interface{}, name string, args ...interface{}) error {
	query, err := d.Raw(name)
	if err != nil {
		return err
	}

	return dbx.SelectContext(ctx, dest, query, args...)
}
