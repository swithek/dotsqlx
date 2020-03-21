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

type Queryerx interface {
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
}

type QueryerxContext interface {
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
}

type QueryRowerx interface {
	QueryRowx(query string, args ...interface{}) *sqlx.Row
}

type QueryRowerxContext interface {
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
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
	*dotsql.DotSql
}

func Wrap(d *dotsql.DotSql) *DotSqlx {
	return &DotSqlx{d}
}

func (d DotSqlx) Preparex(dbx Preparerx, name string) (*sqlx.Stmt, error) {
	query, err := d.Raw(name)
	if err != nil {
		return nil, err
	}

	return dbx.Preparex(query)
}

func (d DotSqlx) PreparexContext(ctx context.Context, dbx PreparerxContext, name string) (*sqlx.Stmt, error) {
	query, err := d.Raw(name)
	if err != nil {
		return nil, err
	}

	return dbx.PreparexContext(ctx, query)
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

func (d DotSqlx) Queryx(dbx Queryerx, name string, args ...interface{}) (*sqlx.Rows, error) {
	query, err := d.Raw(name)
	if err != nil {
		return nil, err
	}

	return dbx.Queryx(query, args...)
}

func (d DotSqlx) QueryxContext(ctx context.Context, dbx QueryerxContext, name string, args ...interface{}) (*sqlx.Rows, error) {
	query, err := d.Raw(name)
	if err != nil {
		return nil, err
	}

	return dbx.QueryxContext(ctx, query, args...)
}

func (d DotSqlx) QueryRowx(dbx QueryRowerx, name string, args ...interface{}) (*sqlx.Row, error) {
	query, err := d.Raw(name)
	if err != nil {
		return nil, err
	}

	return dbx.QueryRowx(query, args...), nil
}

func (d DotSqlx) QueryRowxContext(ctx context.Context, dbx QueryRowerxContext, name string, args ...interface{}) (*sqlx.Row, error) {
	query, err := d.Raw(name)
	if err != nil {
		return nil, err
	}

	return dbx.QueryRowxContext(ctx, query, args...), nil
}

func (d DotSqlx) MustExec(dbx MustExecer, name string, args ...interface{}) sql.Result {
	query, err := d.Raw(name)
	if err != nil {
		panic(err)
	}

	return dbx.MustExec(query, args...)
}

func (d DotSqlx) MustExecContext(ctx context.Context, dbx MustExecerContext, name string, args ...interface{}) sql.Result {
	query, err := d.Raw(name)
	if err != nil {
		panic(err)
	}

	return dbx.MustExecContext(ctx, query, args...)
}

func (d DotSqlx) Rebind(dbx Rebinder, name string) (string, error) {
	query, err := d.Raw(name)
	if err != nil {
		return "", err
	}

	return dbx.Rebind(query), nil
}

func (d DotSqlx) PrepareNamed(dbx NamedPreparer, name string) (*sqlx.NamedStmt, error) {
	query, err := d.Raw(name)
	if err != nil {
		return nil, err
	}

	return dbx.PrepareNamed(query)
}

func (d DotSqlx) PrepareNamedContext(ctx context.Context, dbx NamedPreparerContext, name string) (*sqlx.NamedStmt, error) {
	query, err := d.Raw(name)
	if err != nil {
		return nil, err
	}

	return dbx.PrepareNamedContext(ctx, query)
}

func (d DotSqlx) NamedQuery(dbx NamedQueryer, name string, arg interface{}) (*sqlx.Rows, error) {
	query, err := d.Raw(name)
	if err != nil {
		return nil, err
	}

	return dbx.NamedQuery(query, arg)
}

func (d DotSqlx) NamedQueryContext(ctx context.Context, dbx NamedQueryerContext, name string, arg interface{}) (*sqlx.Rows, error) {
	query, err := d.Raw(name)
	if err != nil {
		return nil, err
	}

	return dbx.NamedQueryContext(ctx, query, arg)
}

func (d DotSqlx) NamedExec(dbx NamedExecer, name string, arg interface{}) (sql.Result, error) {
	query, err := d.Raw(name)
	if err != nil {
		return nil, err
	}

	return dbx.NamedExec(query, arg)
}

func (d DotSqlx) NamedExecContext(ctx context.Context, dbx NamedExecerContext, name string, arg interface{}) (sql.Result, error) {
	query, err := d.Raw(name)
	if err != nil {
		return nil, err
	}

	return dbx.NamedExecContext(ctx, query, arg)
}

func (d DotSqlx) BindNamed(dbx NamedBinder, name string, arg interface{}) (string, []interface{}, error) {
	query, err := d.Raw(name)
	if err != nil {
		return "", nil, err
	}

	return dbx.BindNamed(query, arg)
}

func (d DotSqlx) In(name string, args ...interface{}) (string, []interface{}, error) {
	query, err := d.Raw(name)
	if err != nil {
		return "", nil, err
	}

	return sqlx.In(query, args...)
}
