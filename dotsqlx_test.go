package dotsqlx

import (
	"context"
	"database/sql"
	"testing"

	"github.com/gchaincl/dotsql"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const queries = `
-- name: insert
INSERT INTO numbers (nr1) VALUES(?)

-- name: select
SELECT nr FROM numbers WHERE nr = ?`

type number struct {
	Nr int `db:"nr"`
}

type sqlResult struct{}

func (s sqlResult) LastInsertId() (int64, error) {
	return 1, nil
}

func (s sqlResult) RowsAffected() (int64, error) {
	return 1, nil
}

func newDot(t *testing.T, q string) *DotSqlx {
	d, err := dotsql.LoadFromString(q)
	require.Nil(t, err)
	return &DotSqlx{d}
}

func TestWrap(t *testing.T) {
	d := Wrap(&dotsql.DotSql{})
	assert.NotNil(t, d)
	assert.NotNil(t, d.DotSql)
}

func TestPreparex(t *testing.T) {
	require.Implements(t, (*Preparerx)(nil), new(sqlx.DB))

	dot := newDot(t, queries)

	preparerxStub := func(err error) *PreparerxMock {
		return &PreparerxMock{
			PreparexFunc: func(_ string) (*sqlx.Stmt, error) {
				if err != nil {
					return nil, err
				}
				return &sqlx.Stmt{}, nil
			},
		}
	}

	// query not found
	p := preparerxStub(nil)
	stmt, err := dot.Preparex(p, "insert123")
	assert.Nil(t, stmt)
	assert.NotNil(t, err)
	assert.Zero(t, len(p.PreparexCalls()))

	// error returned by db
	p = preparerxStub(assert.AnError)
	stmt, err = dot.Preparex(p, "select")
	assert.Nil(t, stmt)
	assert.NotNil(t, err)
	ff := p.PreparexCalls()
	require.Equal(t, 1, len(ff))
	assert.NotZero(t, ff[0].Query)

	// successful call
	p = preparerxStub(nil)
	stmt, err = dot.Preparex(p, "select")
	assert.NotNil(t, stmt)
	assert.Nil(t, err)
	ff = p.PreparexCalls()
	require.Equal(t, 1, len(ff))
	assert.NotZero(t, ff[0].Query)
}

func TestPreparexContext(t *testing.T) {
	require.Implements(t, (*PreparerxContext)(nil), new(sqlx.DB))

	dot := newDot(t, queries)

	preparerxCtxStub := func(err error) *PreparerxContextMock {
		return &PreparerxContextMock{
			PreparexContextFunc: func(_ context.Context, _ string) (*sqlx.Stmt, error) {
				if err != nil {
					return nil, err
				}
				return &sqlx.Stmt{}, nil
			},
		}
	}

	// query not found
	p := preparerxCtxStub(nil)
	stmt, err := dot.PreparexContext(context.Background(), p, "insert123")
	assert.Nil(t, stmt)
	assert.NotNil(t, err)
	assert.Zero(t, len(p.PreparexContextCalls()))

	// error returned by db
	p = preparerxCtxStub(assert.AnError)
	stmt, err = dot.PreparexContext(context.Background(), p, "select")
	assert.Nil(t, stmt)
	assert.NotNil(t, err)
	ff := p.PreparexContextCalls()
	require.Equal(t, 1, len(ff))
	assert.NotNil(t, ff[0].Ctx)
	assert.NotZero(t, ff[0].Query)

	// successful call
	p = preparerxCtxStub(nil)
	stmt, err = dot.PreparexContext(context.Background(), p, "select")
	assert.NotNil(t, stmt)
	assert.Nil(t, err)
	ff = p.PreparexContextCalls()
	require.Equal(t, 1, len(ff))
	assert.NotNil(t, ff[0].Ctx)
	assert.NotZero(t, ff[0].Query)
}

func TestGet(t *testing.T) {
	require.Implements(t, (*Getter)(nil), new(sqlx.DB))

	dot := newDot(t, queries)

	getterStub := func(v int, err error) *GetterMock {
		return &GetterMock{
			GetFunc: func(nr interface{}, _ string, _ ...interface{}) error {
				if err != nil {
					return err
				}

				nr.(*number).Nr = v
				return nil
			},
		}
	}

	// query not found
	nr := &number{}
	g := getterStub(5, nil)
	err := dot.Get(g, nr, "insert123", 1)
	assert.Zero(t, nr.Nr)
	assert.NotNil(t, err)
	assert.Zero(t, len(g.GetCalls()))

	// error returned by db
	nr = &number{}
	g = getterStub(0, assert.AnError)
	err = dot.Get(g, nr, "select", 1)
	assert.Zero(t, nr.Nr)
	assert.NotNil(t, err)
	ff := g.GetCalls()
	require.NotZero(t, len(ff))
	assert.NotNil(t, ff[0].Dest)
	assert.NotZero(t, ff[0].Query)
	require.NotEmpty(t, ff[0].Args)
	assert.Equal(t, 1, ff[0].Args[0])

	// successful call
	nr = &number{}
	g = getterStub(5, nil)
	err = dot.Get(g, nr, "select", 1)
	assert.Equal(t, 5, nr.Nr)
	assert.Nil(t, err)
	ff = g.GetCalls()
	require.NotZero(t, len(ff))
	assert.NotNil(t, ff[0].Dest)
	assert.NotZero(t, ff[0].Query)
	require.NotEmpty(t, ff[0].Args)
	assert.Equal(t, 1, ff[0].Args[0])
}

func TestGetContext(t *testing.T) {
	require.Implements(t, (*GetterContext)(nil), new(sqlx.DB))

	dot := newDot(t, queries)

	getterCtxStub := func(v int, err error) *GetterContextMock {
		return &GetterContextMock{
			GetContextFunc: func(_ context.Context, nr interface{}, _ string, _ ...interface{}) error {
				if err != nil {
					return err
				}

				nr.(*number).Nr = v
				return nil
			},
		}
	}

	// query not found
	nr := &number{}
	g := getterCtxStub(5, nil)
	err := dot.GetContext(context.Background(), g, nr, "insert123", 1)
	assert.Zero(t, nr.Nr)
	assert.NotNil(t, err)
	assert.Zero(t, len(g.GetContextCalls()))

	// error returned by db
	nr = &number{}
	g = getterCtxStub(0, assert.AnError)
	err = dot.GetContext(context.Background(), g, nr, "select", 1)
	assert.Zero(t, nr.Nr)
	assert.NotNil(t, err)
	ff := g.GetContextCalls()
	require.NotZero(t, len(ff))
	assert.NotNil(t, ff[0].Ctx)
	assert.NotNil(t, ff[0].Dest)
	assert.NotZero(t, ff[0].Query)
	require.NotEmpty(t, ff[0].Args)
	assert.Equal(t, 1, ff[0].Args[0])

	// successful call
	nr = &number{}
	g = getterCtxStub(5, nil)
	err = dot.GetContext(context.Background(), g, nr, "select", 1)
	assert.Equal(t, 5, nr.Nr)
	assert.Nil(t, err)
	ff = g.GetContextCalls()
	require.NotZero(t, len(ff))
	assert.NotNil(t, ff[0].Ctx)
	assert.NotNil(t, ff[0].Dest)
	assert.NotZero(t, ff[0].Query)
	require.NotEmpty(t, ff[0].Args)
	assert.Equal(t, 1, ff[0].Args[0])
}

func TestSelect(t *testing.T) {
	require.Implements(t, (*Selecter)(nil), new(sqlx.DB))

	dot := newDot(t, queries)

	selecterStub := func(vv []int, err error) *SelecterMock {
		return &SelecterMock{
			SelectFunc: func(nrs interface{}, _ string, _ ...interface{}) error {
				if err != nil {
					return err
				}

				nrsArr := nrs.(*[]number)
				for _, v := range vv {
					*nrsArr = append(*nrsArr, number{v})
				}

				return nil
			},
		}
	}

	// query not found
	nrs := []number{}
	s := selecterStub([]int{5, 10, 20}, nil)
	err := dot.Select(s, &nrs, "insert123", 1)
	assert.Empty(t, nrs)
	assert.NotNil(t, err)
	assert.Zero(t, len(s.SelectCalls()))

	// error returned by db
	nrs = []number{}
	s = selecterStub(nil, assert.AnError)
	err = dot.Select(s, &nrs, "select", 1)
	assert.Empty(t, nrs)
	assert.NotNil(t, err)
	ff := s.SelectCalls()
	require.Equal(t, 1, len(ff))
	assert.NotNil(t, ff[0].Dest)
	assert.NotZero(t, ff[0].Query)
	require.NotEmpty(t, ff[0].Args)
	assert.Equal(t, 1, ff[0].Args[0])

	// successful call
	nrs = []number{}
	s = selecterStub([]int{5, 10, 20}, nil)
	err = dot.Select(s, &nrs, "select", 1)
	assert.Equal(t, []number{{5}, {10}, {20}}, nrs)
	assert.Nil(t, err)
	ff = s.SelectCalls()
	require.Equal(t, 1, len(ff))
	assert.NotNil(t, ff[0].Dest)
	assert.NotZero(t, ff[0].Query)
	require.NotEmpty(t, ff[0].Args)
	assert.Equal(t, 1, ff[0].Args[0])
}

func TestSelectContext(t *testing.T) {
	require.Implements(t, (*SelecterContext)(nil), new(sqlx.DB))

	dot := newDot(t, queries)

	selecterCtxStub := func(vv []int, err error) *SelecterContextMock {
		return &SelecterContextMock{
			SelectContextFunc: func(_ context.Context, nrs interface{}, _ string, _ ...interface{}) error {
				if err != nil {
					return err
				}

				nrsArr := nrs.(*[]number)
				for _, v := range vv {
					*nrsArr = append(*nrsArr, number{v})
				}

				return nil
			},
		}
	}

	// query not found
	nrs := []number{}
	s := selecterCtxStub([]int{5, 10, 20}, nil)
	err := dot.SelectContext(context.Background(), s, &nrs, "insert123", 1)
	assert.Empty(t, nrs)
	assert.NotNil(t, err)
	assert.Zero(t, len(s.SelectContextCalls()))

	// error returned by db
	nrs = []number{}
	s = selecterCtxStub(nil, assert.AnError)
	err = dot.SelectContext(context.Background(), s, &nrs, "select", 1)
	assert.Empty(t, nrs)
	assert.NotNil(t, err)
	ff := s.SelectContextCalls()
	require.Equal(t, 1, len(ff))
	assert.NotNil(t, ff[0].Ctx)
	assert.NotNil(t, ff[0].Dest)
	assert.NotZero(t, ff[0].Query)
	require.NotEmpty(t, ff[0].Args)
	assert.Equal(t, 1, ff[0].Args[0])

	// successful call
	nrs = []number{}
	s = selecterCtxStub([]int{5, 10, 20}, nil)
	err = dot.SelectContext(context.Background(), s, &nrs, "select", 1)
	assert.Equal(t, []number{{5}, {10}, {20}}, nrs)
	assert.Nil(t, err)
	ff = s.SelectContextCalls()
	require.Equal(t, 1, len(ff))
	assert.NotNil(t, ff[0].Ctx)
	assert.NotNil(t, ff[0].Dest)
	assert.NotZero(t, ff[0].Query)
	require.NotEmpty(t, ff[0].Args)
	assert.Equal(t, 1, ff[0].Args[0])
}

func TestQueryx(t *testing.T) {
	require.Implements(t, (*Queryerx)(nil), new(sqlx.DB))

	dot := newDot(t, queries)

	queryerxStub := func(err error) *QueryerxMock {
		return &QueryerxMock{
			QueryxFunc: func(_ string, _ ...interface{}) (*sqlx.Rows, error) {
				if err != nil {
					return nil, err
				}

				return &sqlx.Rows{}, nil
			},
		}
	}

	// query not found
	q := queryerxStub(nil)
	rows, err := dot.Queryx(q, "insert123", 1)
	assert.Nil(t, rows)
	assert.NotNil(t, err)
	assert.Zero(t, len(q.QueryxCalls()))

	// error returned by db
	q = queryerxStub(assert.AnError)
	rows, err = dot.Queryx(q, "select", 1)
	assert.Nil(t, rows)
	assert.NotNil(t, err)
	ff := q.QueryxCalls()
	require.Equal(t, 1, len(ff))
	assert.NotZero(t, ff[0].Query)
	require.NotEmpty(t, ff[0].Args)
	assert.Equal(t, 1, ff[0].Args[0])

	// successful call
	q = queryerxStub(nil)
	rows, err = dot.Queryx(q, "select", 1)
	assert.NotNil(t, rows)
	assert.Nil(t, err)
	ff = q.QueryxCalls()
	require.Equal(t, 1, len(ff))
	assert.NotZero(t, ff[0].Query)
	require.NotEmpty(t, ff[0].Args)
	assert.Equal(t, 1, ff[0].Args[0])
}

func TestQueryxContext(t *testing.T) {
	require.Implements(t, (*QueryerxContext)(nil), new(sqlx.DB))

	dot := newDot(t, queries)

	queryerxCtxStub := func(err error) *QueryerxContextMock {
		return &QueryerxContextMock{
			QueryxContextFunc: func(_ context.Context, _ string, _ ...interface{}) (*sqlx.Rows, error) {
				if err != nil {
					return nil, err
				}

				return &sqlx.Rows{}, nil
			},
		}
	}

	// query not found
	q := queryerxCtxStub(nil)
	rows, err := dot.QueryxContext(context.Background(), q, "insert123", 1)
	assert.Nil(t, rows)
	assert.NotNil(t, err)
	assert.Zero(t, len(q.QueryxContextCalls()))

	// error returned by db
	q = queryerxCtxStub(assert.AnError)
	rows, err = dot.QueryxContext(context.Background(), q, "select", 1)
	assert.Nil(t, rows)
	assert.NotNil(t, err)
	ff := q.QueryxContextCalls()
	require.Equal(t, 1, len(ff))
	assert.NotNil(t, ff[0].Ctx)
	assert.NotZero(t, ff[0].Query)
	require.NotEmpty(t, ff[0].Args)
	assert.Equal(t, 1, ff[0].Args[0])

	// successful call
	q = queryerxCtxStub(nil)
	rows, err = dot.QueryxContext(context.Background(), q, "select", 1)
	assert.NotNil(t, rows)
	assert.Nil(t, err)
	ff = q.QueryxContextCalls()
	require.Equal(t, 1, len(ff))
	assert.NotNil(t, ff[0].Ctx)
	assert.NotZero(t, ff[0].Query)
	require.NotEmpty(t, ff[0].Args)
	assert.Equal(t, 1, ff[0].Args[0])
}

func TestQueryRowx(t *testing.T) {
	require.Implements(t, (*QueryRowerx)(nil), new(sqlx.DB))

	dot := newDot(t, queries)

	queryRowerxStub := func() *QueryRowerxMock {
		return &QueryRowerxMock{
			QueryRowxFunc: func(_ string, _ ...interface{}) *sqlx.Row {
				return &sqlx.Row{}
			},
		}
	}

	// query not found
	q := queryRowerxStub()
	row, err := dot.QueryRowx(q, "insert123", 1)
	assert.Nil(t, row)
	assert.NotNil(t, err)
	assert.Zero(t, len(q.QueryRowxCalls()))

	// successful call
	q = queryRowerxStub()
	row, err = dot.QueryRowx(q, "select", 1)
	assert.NotNil(t, row)
	assert.Nil(t, err)
	ff := q.QueryRowxCalls()
	require.Equal(t, 1, len(ff))
	assert.NotZero(t, ff[0].Query)
	require.NotEmpty(t, ff[0].Args)
	assert.Equal(t, 1, ff[0].Args[0])
}

func TestQueryRowxContext(t *testing.T) {
	require.Implements(t, (*QueryRowerxContext)(nil), new(sqlx.DB))

	dot := newDot(t, queries)

	queryRowerxCtxStub := func() *QueryRowerxContextMock {
		return &QueryRowerxContextMock{
			QueryRowxContextFunc: func(_ context.Context, _ string, _ ...interface{}) *sqlx.Row {
				return &sqlx.Row{}
			},
		}
	}

	// query not found
	q := queryRowerxCtxStub()
	row, err := dot.QueryRowxContext(context.Background(), q, "insert123", 1)
	assert.Nil(t, row)
	assert.NotNil(t, err)
	assert.Zero(t, len(q.QueryRowxContextCalls()))

	// successful call
	q = queryRowerxCtxStub()
	row, err = dot.QueryRowxContext(context.Background(), q, "select", 1)
	assert.NotNil(t, row)
	assert.Nil(t, err)
	ff := q.QueryRowxContextCalls()
	require.Equal(t, 1, len(ff))
	assert.NotZero(t, ff[0].Ctx)
	assert.NotZero(t, ff[0].Query)
	require.NotEmpty(t, ff[0].Args)
	assert.Equal(t, 1, ff[0].Args[0])
}

func TestMustExec(t *testing.T) {
	require.Implements(t, (*MustExecer)(nil), new(sqlx.DB))

	dot := newDot(t, queries)

	mustExecerStub := func() *MustExecerMock {
		return &MustExecerMock{
			MustExecFunc: func(_ string, _ ...interface{}) sql.Result {
				return sqlResult{}
			},
		}
	}

	// query not found
	e := mustExecerStub()
	var res sql.Result
	assert.Panics(t, func() {
		res = dot.MustExec(e, "insert123", 1)
	})
	assert.Nil(t, res)
	assert.Zero(t, len(e.MustExecCalls()))

	// successful call
	e = mustExecerStub()
	assert.NotPanics(t, func() {
		res = dot.MustExec(e, "insert", 1)
	})
	assert.NotNil(t, res)
	ff := e.MustExecCalls()
	require.Equal(t, 1, len(ff))
	assert.NotZero(t, ff[0].Query)
	require.NotEmpty(t, ff[0].Args)
	assert.Equal(t, 1, ff[0].Args[0])
}

func TestMustExecContext(t *testing.T) {
	require.Implements(t, (*MustExecerContext)(nil), new(sqlx.DB))

	dot := newDot(t, queries)

	mustExecerCtxStub := func() *MustExecerContextMock {
		return &MustExecerContextMock{
			MustExecContextFunc: func(_ context.Context, _ string, _ ...interface{}) sql.Result {
				return sqlResult{}
			},
		}
	}

	// query not found
	e := mustExecerCtxStub()
	var res sql.Result
	assert.Panics(t, func() {
		res = dot.MustExecContext(context.Background(), e, "insert123", 1)
	})
	assert.Nil(t, res)
	assert.Zero(t, len(e.MustExecContextCalls()))

	// successful call
	e = mustExecerCtxStub()
	assert.NotPanics(t, func() {
		res = dot.MustExecContext(context.Background(), e, "insert", 1)
	})
	assert.NotNil(t, res)
	ff := e.MustExecContextCalls()
	require.Equal(t, 1, len(ff))
	assert.NotNil(t, ff[0].Ctx)
	assert.NotZero(t, ff[0].Query)
	require.NotEmpty(t, ff[0].Args)
	assert.Equal(t, 1, ff[0].Args[0])
}

func TestRebind(t *testing.T) {
	require.Implements(t, (*Rebinder)(nil), new(sqlx.DB))

	dot := newDot(t, queries)

	rebinderStub := func() *RebinderMock {
		return &RebinderMock{
			RebindFunc: func(_ string) string {
				return "test"
			},
		}
	}

	// query not found
	r := rebinderStub()
	str, err := dot.Rebind(r, "insert123")
	assert.Zero(t, str)
	assert.NotNil(t, err)
	assert.Zero(t, len(r.RebindCalls()))

	// successful call
	r = rebinderStub()
	str, err = dot.Rebind(r, "select")
	assert.NotZero(t, str)
	assert.Nil(t, err)
	ff := r.RebindCalls()
	require.Equal(t, 1, len(ff))
	assert.NotZero(t, ff[0].Query)
}

func TestPrepareNamed(t *testing.T) {
	require.Implements(t, (*NamedPreparer)(nil), new(sqlx.DB))

	dot := newDot(t, queries)

	namedPreparerStub := func(err error) *NamedPreparerMock {
		return &NamedPreparerMock{
			PrepareNamedFunc: func(_ string) (*sqlx.NamedStmt, error) {
				if err != nil {
					return nil, err
				}
				return &sqlx.NamedStmt{}, nil
			},
		}
	}

	// query not found
	p := namedPreparerStub(nil)
	stmt, err := dot.PrepareNamed(p, "insert123")
	assert.Nil(t, stmt)
	assert.NotNil(t, err)
	assert.Zero(t, len(p.PrepareNamedCalls()))

	// error returned by db
	p = namedPreparerStub(assert.AnError)
	stmt, err = dot.PrepareNamed(p, "select")
	assert.Nil(t, stmt)
	assert.NotNil(t, err)
	ff := p.PrepareNamedCalls()
	require.Equal(t, 1, len(ff))
	assert.NotZero(t, ff[0].Query)

	// successful call
	p = namedPreparerStub(nil)
	stmt, err = dot.PrepareNamed(p, "select")
	assert.NotNil(t, stmt)
	assert.Nil(t, err)
	ff = p.PrepareNamedCalls()
	require.Equal(t, 1, len(ff))
	assert.NotZero(t, ff[0].Query)
}

func TestPrepareNamedContext(t *testing.T) {
	require.Implements(t, (*NamedPreparerContext)(nil), new(sqlx.DB))

	dot := newDot(t, queries)

	namedPreparerCtxStub := func(err error) *NamedPreparerContextMock {
		return &NamedPreparerContextMock{
			PrepareNamedContextFunc: func(_ context.Context, _ string) (*sqlx.NamedStmt, error) {
				if err != nil {
					return nil, err
				}
				return &sqlx.NamedStmt{}, nil
			},
		}
	}

	// query not found
	p := namedPreparerCtxStub(nil)
	stmt, err := dot.PrepareNamedContext(context.Background(), p, "insert123")
	assert.Nil(t, stmt)
	assert.NotNil(t, err)
	assert.Zero(t, len(p.PrepareNamedContextCalls()))

	// error returned by db
	p = namedPreparerCtxStub(assert.AnError)
	stmt, err = dot.PrepareNamedContext(context.Background(), p, "select")
	assert.Nil(t, stmt)
	assert.NotNil(t, err)
	ff := p.PrepareNamedContextCalls()
	require.Equal(t, 1, len(ff))
	assert.NotNil(t, ff[0].Ctx)
	assert.NotZero(t, ff[0].Query)

	// successful call
	p = namedPreparerCtxStub(nil)
	stmt, err = dot.PrepareNamedContext(context.Background(), p, "select")
	assert.NotNil(t, stmt)
	assert.Nil(t, err)
	ff = p.PrepareNamedContextCalls()
	require.Equal(t, 1, len(ff))
	assert.NotNil(t, ff[0].Ctx)
	assert.NotZero(t, ff[0].Query)
}

func TestNamedQuery(t *testing.T) {
	require.Implements(t, (*NamedQueryer)(nil), new(sqlx.DB))

	dot := newDot(t, queries)

	namedQueryerStub := func(err error) *NamedQueryerMock {
		return &NamedQueryerMock{
			NamedQueryFunc: func(_ string, _ interface{}) (*sqlx.Rows, error) {
				if err != nil {
					return nil, err
				}

				return &sqlx.Rows{}, nil
			},
		}
	}

	// query not found
	q := namedQueryerStub(nil)
	rows, err := dot.NamedQuery(q, "insert123", 1)
	assert.Nil(t, rows)
	assert.NotNil(t, err)
	assert.Zero(t, len(q.NamedQueryCalls()))

	// error returned by db
	q = namedQueryerStub(assert.AnError)
	rows, err = dot.NamedQuery(q, "select", 1)
	assert.Nil(t, rows)
	assert.NotNil(t, err)
	ff := q.NamedQueryCalls()
	require.Equal(t, 1, len(ff))
	assert.NotZero(t, ff[0].Query)
	assert.NotNil(t, ff[0].Arg)

	// successful call
	q = namedQueryerStub(nil)
	rows, err = dot.NamedQuery(q, "select", 1)
	assert.NotNil(t, rows)
	assert.Nil(t, err)
	ff = q.NamedQueryCalls()
	require.Equal(t, 1, len(ff))
	assert.NotZero(t, ff[0].Query)
	assert.NotNil(t, ff[0].Arg)
}

func TestNamedQueryContext(t *testing.T) {
	require.Implements(t, (*NamedQueryerContext)(nil), new(sqlx.DB))

	dot := newDot(t, queries)

	namedQueryerCtxStub := func(err error) *NamedQueryerContextMock {
		return &NamedQueryerContextMock{
			NamedQueryContextFunc: func(_ context.Context, _ string, _ interface{}) (*sqlx.Rows, error) {
				if err != nil {
					return nil, err
				}

				return &sqlx.Rows{}, nil
			},
		}
	}

	// query not found
	q := namedQueryerCtxStub(nil)
	rows, err := dot.NamedQueryContext(context.Background(), q, "insert123", 1)
	assert.Nil(t, rows)
	assert.NotNil(t, err)
	assert.Zero(t, len(q.NamedQueryContextCalls()))

	// error returned by db
	q = namedQueryerCtxStub(assert.AnError)
	rows, err = dot.NamedQueryContext(context.Background(), q, "select", 1)
	assert.Nil(t, rows)
	assert.NotNil(t, err)
	ff := q.NamedQueryContextCalls()
	require.Equal(t, 1, len(ff))
	assert.NotNil(t, ff[0].Ctx)
	assert.NotZero(t, ff[0].Query)
	assert.NotNil(t, ff[0].Arg)

	// successful call
	q = namedQueryerCtxStub(nil)
	rows, err = dot.NamedQueryContext(context.Background(), q, "select", 1)
	assert.NotNil(t, rows)
	assert.Nil(t, err)
	ff = q.NamedQueryContextCalls()
	require.Equal(t, 1, len(ff))
	assert.NotNil(t, ff[0].Ctx)
	assert.NotZero(t, ff[0].Query)
	assert.NotNil(t, ff[0].Arg)
}

func TestNamedExec(t *testing.T) {
	require.Implements(t, (*NamedExecer)(nil), new(sqlx.DB))

	dot := newDot(t, queries)

	namedExecerStub := func(err error) *NamedExecerMock {
		return &NamedExecerMock{
			NamedExecFunc: func(_ string, _ interface{}) (sql.Result, error) {
				if err != nil {
					return nil, err
				}

				return sqlResult{}, nil
			},
		}
	}

	// query not found
	e := namedExecerStub(nil)
	res, err := dot.NamedExec(e, "insert123", 1)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Zero(t, len(e.NamedExecCalls()))

	// error returned by db
	e = namedExecerStub(assert.AnError)
	res, err = dot.NamedExec(e, "insert", 1)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	ff := e.NamedExecCalls()
	require.Equal(t, 1, len(ff))
	assert.NotZero(t, ff[0].Query)
	assert.NotNil(t, ff[0].Arg)

	// successful call
	e = namedExecerStub(nil)
	res, err = dot.NamedExec(e, "insert", 1)
	assert.NotNil(t, res)
	assert.Nil(t, err)
	ff = e.NamedExecCalls()
	require.Equal(t, 1, len(ff))
	assert.NotZero(t, ff[0].Query)
	assert.NotNil(t, ff[0].Arg)
}

func TestNamedExecContext(t *testing.T) {
	require.Implements(t, (*NamedExecerContext)(nil), new(sqlx.DB))

	dot := newDot(t, queries)

	namedExecerCtxStub := func(err error) *NamedExecerContextMock {
		return &NamedExecerContextMock{
			NamedExecContextFunc: func(_ context.Context, _ string, _ interface{}) (sql.Result, error) {
				if err != nil {
					return nil, err
				}

				return sqlResult{}, nil
			},
		}
	}

	// query not found
	e := namedExecerCtxStub(nil)
	res, err := dot.NamedExecContext(context.Background(), e, "insert123", 1)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Zero(t, len(e.NamedExecContextCalls()))

	// error returned by db
	e = namedExecerCtxStub(assert.AnError)
	res, err = dot.NamedExecContext(context.Background(), e, "insert", 1)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	ff := e.NamedExecContextCalls()
	require.Equal(t, 1, len(ff))
	assert.NotNil(t, ff[0].Ctx)
	assert.NotZero(t, ff[0].Query)
	assert.NotNil(t, ff[0].Arg)

	// successful call
	e = namedExecerCtxStub(nil)
	res, err = dot.NamedExecContext(context.Background(), e, "insert", 1)
	assert.NotNil(t, res)
	assert.Nil(t, err)
	ff = e.NamedExecContextCalls()
	require.Equal(t, 1, len(ff))
	assert.NotNil(t, ff[0].Ctx)
	assert.NotZero(t, ff[0].Query)
	assert.NotNil(t, ff[0].Arg)
}

func TestBindNamed(t *testing.T) {
	require.Implements(t, (*NamedBinder)(nil), new(sqlx.DB))

	dot := newDot(t, queries)

	namedBinderStub := func(err error) *NamedBinderMock {
		return &NamedBinderMock{
			BindNamedFunc: func(_ string, _ interface{}) (string, []interface{}, error) {
				if err != nil {
					return "", nil, err
				}

				return "test", []interface{}{1, 2, 3}, nil
			},
		}
	}

	// query not found
	b := namedBinderStub(nil)
	res1, res2, err := dot.BindNamed(b, "insert123", 1)
	assert.Zero(t, res1)
	assert.Nil(t, res2)
	assert.NotNil(t, err)
	assert.Zero(t, len(b.BindNamedCalls()))

	// error returned by db
	b = namedBinderStub(assert.AnError)
	res1, res2, err = dot.BindNamed(b, "insert", 1)
	assert.Zero(t, res1)
	assert.Nil(t, res2)
	assert.NotNil(t, err)
	ff := b.BindNamedCalls()
	require.Equal(t, 1, len(ff))
	assert.NotZero(t, ff[0].Query)
	assert.NotNil(t, ff[0].Arg)

	// successful call
	b = namedBinderStub(nil)
	res1, res2, err = dot.BindNamed(b, "insert", 1)
	assert.NotZero(t, res1)
	assert.NotNil(t, res2)
	assert.Nil(t, err)
	ff = b.BindNamedCalls()
	require.Equal(t, 1, len(ff))
	assert.NotZero(t, ff[0].Query)
	assert.NotNil(t, ff[0].Arg)
}

func TestIn(t *testing.T) {
	dot := newDot(t, `
--name: select
SELECT * FROM numbers WHERE nr IN (?)

--name: select2
SELECT * FROM numbers WHERE nr`)

	// query not found
	res1, res2, err := dot.In("select123", []int{1, 2, 3})
	assert.Zero(t, res1)
	assert.Nil(t, res2)
	assert.NotNil(t, err)

	// error returned by sqlx
	res1, res2, err = dot.In("select2", []int{1, 2, 3})
	assert.Zero(t, res1)
	assert.Nil(t, res2)
	assert.NotNil(t, err)

	// successful call
	res1, res2, err = dot.In("select", []int{1, 2, 3})
	assert.Equal(t, "SELECT * FROM numbers WHERE nr IN (?, ?, ?)", res1)
	assert.NotNil(t, res2)
	assert.Nil(t, err)
}
