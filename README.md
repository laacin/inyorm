# Inyorm

[![Go Reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/laacin/inyorm)

##### Inyorm is a fully declarative ORM for Go, designed for clarity, type safety, and predictable SQL generation.

## Overview

- SQL-like declarative API
- Strong typing
- Fast mapping with minimal reflection
- Lightweight, explicit, and fast
- Highly customizable

```go
q, c := qe.NewSelect(ctx, "users")

var (
    id = c.Col("id")
    fk = c.Col("user_id", "posts")
    postNum = c.Col("id", "posts").Count()
)

q.Select(c.All(), postNum)
q.Join("posts").On(fk).Equal(id)
q.Where(id).Equal(c.Param(43))
q.Limit(1)

u := &User{}
if err := q.Run(u); err != nil {
    log.Fatal(err)
}
```

## Getting started

### Install

```bash
go get -u github.com/laacin/inyorm
```

### Minimal setup

Inyorm is at its core a declarative query builder and object mapper. It doesn’t manage connections or drivers,
so you remain fully in control of the database layer.

```go
package main

import (
	"database/sql"
	"log"

	"github.com/laacin/inyorm"
)

func main() {
    db, err := sql.Open("sqlite3", "./data.db")
    if err != nil {
        log.Fatal(err)
    }

    qe := inyorm.New("sqlite3", db, &inyorm.Options{})

    ctx := context.Background()

    q, c := qe.NewSelect(ctx, "table")
    // q, c := qe.NewInsert(ctx, "table")
    // q, c := qe.NewUpdate(ctx, "table")
    // q, c := qe.NewDelete(ctx, "table")
}
```
