# dotsqlx

[![GoDoc](https://godoc.org/github.com/swithek/dotsqlx?status.png)](https://godoc.org/github.com/swithek/dotsqlx)
[![Build status](https://travis-ci.org/swithek/dotsqlx.svg?branch=master)](https://travis-ci.org/swithek/dotsqlx)
[![Go Report Card](https://goreportcard.com/badge/github.com/swithek/dotsqlx)](https://goreportcard.com/report/github.com/swithek/dotsqlx)

[Dotsql](https://github.com/gchaincl/dotsql) wrapper allowing seemless work with [jmoiron/sqlx](https://github.com/jmoiron/sqlx)

## Installation
```
go get github.com/swithek/dotsqlx
```

## Usage
```go
// connect to db and obtain a new sqlx db instance.
dbx, err := sqlx.Connect("postgres", "user=foo dbname=bar sslmode=disable")
if err != nil {
    // handle error
}

// load your queries and obtain a new dotsql instance.
dot, err := dotsql.LoadFromFile("./queries.sql")
if err != nil {
    // handle error
}

// wrap dotsql's instance, extend its logic and allow sqlx support.
dotx := dotsqlx.Wrap(dot)

// execute named dotsql queries in sqlx methods.
var users []Users
if err = dotx.Select(dbx, &users, "select_users"); err != nil {
    // handle error
}

// you can use dotsql's methods as well.
res, err := dotx.Exec(dbx, "insert_user")
if err != nil {
    // handle error
}

```
