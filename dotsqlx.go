package dotsqlx

import (
	"context"
	"database/sql"

	"github.com/gchaincl/dotsql"
	"github.com/jmoiron/sqlx"
)

// Preparerx is an interface used by Preparex.
type Preparerx interface {
	Preparex(query string) (*sqlx.Stmt, error)
}

// PreparerxContext is an interface used by PreparexContext.
type PreparerxContext interface {
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
}

// Getter is an interface used by Get.
type Getter interface {
	Get(dest interface{}, query string, args ...interface{}) error
}

// GetterContext is an interface used by GetContext.
type GetterContext interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

// Selecter is an interface used by Select.
type Selecter interface {
	Select(dest interface{}, query string, args ...interface{}) error
}

// SelecterContext is an interface used by SelectContext.
type SelecterContext interface {
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

// Queryerx is an interface used by Queryx.
type Queryerx interface {
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
}

// QueryerxContext is an interface used by QueryxContext.
type QueryerxContext interface {
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
}

// QueryRowerx is an interface used by QueryRowx.
type QueryRowerx interface {
	QueryRowx(query string, args ...interface{}) *sqlx.Row
}

// QueryRowerxContext is an interface used by QueryRowxContext.
type QueryRowerxContext interface {
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
}

// MustExecer is an interface used by MustExec.
type MustExecer interface {
	MustExec(query string, args ...interface{}) sql.Result
}

// MustExecerContext is an interface used by MustExecContext.
type MustExecerContext interface {
	MustExecContext(ctx context.Context, query string, args ...interface{}) sql.Result
}

// Rebinder is an interface used by Rebind.
type Rebinder interface {
	Rebind(query string) string
}

// NamedPreparer is an interface used by PrepareNamed.
type NamedPreparer interface {
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
}

// NamedPreparerContext is an interface used by PrepareNamedContext.
type NamedPreparerContext interface {
	PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)
}

// NamedQueryer is an interface used by NamedQuery.
type NamedQueryer interface {
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
}

// NamedQueryerContext is an interface used by NamedQueryContext.
type NamedQueryerContext interface {
	NamedQueryContext(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error)
}

// NamedExecer is an interface used by NamedExec.
type NamedExecer interface {
	NamedExec(query string, arg interface{}) (sql.Result, error)
}

// NamedExecerContext is an interface used by NamedExecContext.
type NamedExecerContext interface {
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}

// NamedBinder is an interface used by BindNamed.
type NamedBinder interface {
	BindNamed(query string, arg interface{}) (string, []interface{}, error)
}

// DotSqlx wraps dotsql.DotSql instance and allows seemless work with
// jmoiron/sqlx.
type DotSqlx struct {
	*dotsql.DotSql
}

// Wrap creates a new DotSqlx instance and embeds provided dotsql.DotSql
// instance into it.
func Wrap(d *dotsql.DotSql) *DotSqlx {
	return &DotSqlx{d}
}

// Preparex is a wrapper for jmoiron/sqlx's Preparex(), using dotsql named
// query.
func (d DotSqlx) Preparex(dbx Preparerx, name string) (*sqlx.Stmt, error) {
	query, err := d.Raw(name)
	if err != nil {
		return nil, err
	}

	return dbx.Preparex(query)
}

// PreparexContext is a wrapper for jmoiron/sqlx's PreparexContext(), using
// dotsql named query.
func (d DotSqlx) PreparexContext(ctx context.Context, dbx PreparerxContext, name string) (*sqlx.Stmt, error) {
	query, err := d.Raw(name)
	if err != nil {
		return nil, err
	}

	return dbx.PreparexContext(ctx, query)
}

// Get is a wrapper for jmoiron/sqlx's Get(), using dotsql named query.
func (d DotSqlx) Get(dbx Getter, dest interface{}, name string, args ...interface{}) error {
	query, err := d.Raw(name)
	if err != nil {
		return err
	}

	return dbx.Get(dest, query, args...)
}

// GetContext is a wrapper for jmoiron/sqlx's GetContext(), using dotsql
// named query.
func (d DotSqlx) GetContext(ctx context.Context, dbx GetterContext, dest interface{}, name string, args ...interface{}) error {
	query, err := d.Raw(name)
	if err != nil {
		return err
	}

	return dbx.GetContext(ctx, dest, query, args...)
}

// Select is a wrapper for jmoiron/sqlx's Select(), using dotsql named query.
func (d DotSqlx) Select(dbx Selecter, dest interface{}, name string, args ...interface{}) error {
	query, err := d.Raw(name)
	if err != nil {
		return err
	}

	return dbx.Select(dest, query, args...)
}

// SelectContext is a wrapper for jmoiron/sqlx's SelectContext(), using
// dotsql named query.
func (d DotSqlx) SelectContext(ctx context.Context, dbx SelecterContext, dest interface{}, name string, args ...interface{}) error {
	query, err := d.Raw(name)
	if err != nil {
		return err
	}

	return dbx.SelectContext(ctx, dest, query, args...)
}

// Queryx is a wrapper for jmoiron/sqlx's Queryx(), using dotsql named query.
func (d DotSqlx) Queryx(dbx Queryerx, name string, args ...interface{}) (*sqlx.Rows, error) {
	query, err := d.Raw(name)
	if err != nil {
		return nil, err
	}

	return dbx.Queryx(query, args...)
}

// QueryxContext is a wrapper for jmoiron/sqlx's QueryxContext(), using
// dotsql named query.
func (d DotSqlx) QueryxContext(ctx context.Context, dbx QueryerxContext, name string, args ...interface{}) (*sqlx.Rows, error) {
	query, err := d.Raw(name)
	if err != nil {
		return nil, err
	}

	return dbx.QueryxContext(ctx, query, args...)
}

// QueryRowx is a wrapper for jmoiron/sqlx's QueryRowx(), using dotsql
// named query.
func (d DotSqlx) QueryRowx(dbx QueryRowerx, name string, args ...interface{}) (*sqlx.Row, error) {
	query, err := d.Raw(name)
	if err != nil {
		return nil, err
	}

	return dbx.QueryRowx(query, args...), nil
}

// QueryRowxContext is a wrapper for jmoiron/sqlx's QueryRowxContext(), using
// dotsql named query.
func (d DotSqlx) QueryRowxContext(ctx context.Context, dbx QueryRowerxContext, name string, args ...interface{}) (*sqlx.Row, error) {
	query, err := d.Raw(name)
	if err != nil {
		return nil, err
	}

	return dbx.QueryRowxContext(ctx, query, args...), nil
}

// MustExec is a wrapper for jmoiron/sqlx's MustExec(), using dotsql named
// query.
func (d DotSqlx) MustExec(dbx MustExecer, name string, args ...interface{}) sql.Result {
	query, err := d.Raw(name)
	if err != nil {
		panic(err)
	}

	return dbx.MustExec(query, args...)
}

// MustExecContext is a wrapper for jmoiron/sqlx's MustExecContext(), using
// dotsql named query.
func (d DotSqlx) MustExecContext(ctx context.Context, dbx MustExecerContext, name string, args ...interface{}) sql.Result {
	query, err := d.Raw(name)
	if err != nil {
		panic(err)
	}

	return dbx.MustExecContext(ctx, query, args...)
}

// Rebind is a wrapper for jmoiron/sqlx's Rebind(), using dotsql named query.
func (d DotSqlx) Rebind(dbx Rebinder, name string) (string, error) {
	query, err := d.Raw(name)
	if err != nil {
		return "", err
	}

	return dbx.Rebind(query), nil
}

// PrepareNamed is a wrapper for jmoiron/sqlx's PrepareNamed(), using dotsql
// named query.
func (d DotSqlx) PrepareNamed(dbx NamedPreparer, name string) (*sqlx.NamedStmt, error) {
	query, err := d.Raw(name)
	if err != nil {
		return nil, err
	}

	return dbx.PrepareNamed(query)
}

// PrepareNamedContext is a wrapper for jmoiron/sqlx's PrepareNamedContext(),
// using dotsql named query.
func (d DotSqlx) PrepareNamedContext(ctx context.Context, dbx NamedPreparerContext, name string) (*sqlx.NamedStmt, error) {
	query, err := d.Raw(name)
	if err != nil {
		return nil, err
	}

	return dbx.PrepareNamedContext(ctx, query)
}

// NamedQuery is a wrapper for jmoiron/sqlx's NamedQuery(), using dotsql
// named query.
func (d DotSqlx) NamedQuery(dbx NamedQueryer, name string, arg interface{}) (*sqlx.Rows, error) {
	query, err := d.Raw(name)
	if err != nil {
		return nil, err
	}

	return dbx.NamedQuery(query, arg)
}

// NamedQueryContext is a wrapper for jmoiron/sqlx's NamedQueryContext(),
// using dotsql named query.
func (d DotSqlx) NamedQueryContext(ctx context.Context, dbx NamedQueryerContext, name string, arg interface{}) (*sqlx.Rows, error) {
	query, err := d.Raw(name)
	if err != nil {
		return nil, err
	}

	return dbx.NamedQueryContext(ctx, query, arg)
}

// NamedExec is a wrapper for jmoiron/sqlx's NamedExec(), using dotsql
// named query.
func (d DotSqlx) NamedExec(dbx NamedExecer, name string, arg interface{}) (sql.Result, error) {
	query, err := d.Raw(name)
	if err != nil {
		return nil, err
	}

	return dbx.NamedExec(query, arg)
}

// NamedExecContext is a wrapper for jmoiron/sqlx's NamedExecContext(),
// using dotsql named query.
func (d DotSqlx) NamedExecContext(ctx context.Context, dbx NamedExecerContext, name string, arg interface{}) (sql.Result, error) {
	query, err := d.Raw(name)
	if err != nil {
		return nil, err
	}

	return dbx.NamedExecContext(ctx, query, arg)
}

// BindNamed is a wrapper for jmoiron/sqlx's BindNamed(), using dotsql named
// query.
func (d DotSqlx) BindNamed(dbx NamedBinder, name string, arg interface{}) (string, []interface{}, error) {
	query, err := d.Raw(name)
	if err != nil {
		return "", nil, err
	}

	return dbx.BindNamed(query, arg)
}

// In is a wrapper for jmoiron/sqlx's In(), using dotsql named query.
func (d DotSqlx) In(name string, args ...interface{}) (string, []interface{}, error) {
	query, err := d.Raw(name)
	if err != nil {
		return "", nil, err
	}

	return sqlx.In(query, args...)
}
