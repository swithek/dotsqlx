package dotsqlx

import (
	"context"
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

	preparerStub := func(err error) *PreparerxMock {
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
	p := preparerStub(nil)
	stmt, err := dot.Preparex(p, "insert123")
	assert.Nil(t, stmt)
	assert.NotNil(t, err)
	assert.Zero(t, len(p.PreparexCalls()))

	// error returned by db
	p = preparerStub(assert.AnError)
	stmt, err = dot.Preparex(p, "select")
	assert.Nil(t, stmt)
	assert.NotNil(t, err)
	ff := p.PreparexCalls()
	require.Equal(t, 1, len(ff))
	assert.NotZero(t, ff[0].Query)

	// successful call
	p = preparerStub(nil)
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

	preparerCtxStub := func(err error) *PreparerxContextMock {
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
	p := preparerCtxStub(nil)
	stmt, err := dot.PreparexContext(context.Background(), p, "insert123")
	assert.Nil(t, stmt)
	assert.NotNil(t, err)
	assert.Zero(t, len(p.PreparexContextCalls()))

	// error returned by db
	p = preparerCtxStub(assert.AnError)
	stmt, err = dot.PreparexContext(context.Background(), p, "select")
	assert.Nil(t, stmt)
	assert.NotNil(t, err)
	ff := p.PreparexContextCalls()
	require.Equal(t, 1, len(ff))
	assert.NotNil(t, ff[0].Ctx)
	assert.NotZero(t, ff[0].Query)

	// successful call
	p = preparerCtxStub(nil)
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
