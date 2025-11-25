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

Inyorm is at its core a query builder and object mapper. It doesn’t manage connections or drivers,
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

### Guides

<details>
<summary>Column builder</summary>

```go
// The second return value of each new statement is the Column Builder,
// where you can write all non-literal values.
_, c := qe.NewSelect(ctx, "table")

// ----- Col -----

// Col is the most common method. It references a table column
// and accepts two parameters: the first is the main table’s column name.
c.Col("id")

// To reference another table, pass a second parameter.
c.Col("id", "posts")

// ----- All -----

// All references the wildcard '*'.
// In a joined query, the default All() references the main table.
c.All()

// To reference another table, pass a parameter.
c.All("posts")

// ----- Param -----

// Param writes a placeholder.
// In Inyorm you MUST write explicit parameters each time.
c.Param("id")

// You can also skip parameters for lazy values, useful for prepared statements.
c.Param()

// ----- Concat -----

// Now it's time for complex columns.

// Concat writes a CONCAT() in SQL.
// You can include any value.
c.Concat(c.Col("firstname"), " ", c.Col("lastname"))

// Expected output:
// CONCAT(firstname, ' ', lastname)

// ----- Switch -----

// Switch is a reference to a simple CASE.
// It accepts two parameters: a comparable value and a callback
// where you can compare against literal values.
c.Switch(c.Col("banned"), func(cs inyorm.Case) {
	cs.When(true).Then("Invalid user")
	cs.Else("Valid")
})

// Expected output:
// CASE banned WHEN 1 THEN 'Invalid user' ELSE 'Valid' END

// ----- Search -----

// Search is a reference to a searched CASE.
// It accepts one callback
// where each When accepts a boolean expression,
// such as an expression (built below) or literals.
c.Search(func(cs inyorm.Case) {
	cond := c.Cond(c.Cond("age")).Greater(17)
	cs.When(cond).Then("Adult")
	cs.Else("Kid")
})

// Expected output:
// CASE WHEN (age > 17) THEN 'Adult' ELSE 'Kid' END

// ----- Cond -----

// Cond is a helper used to create a condition.
// Currently, only used in the Search method.
c.Cond(c.Cond("banned")).IsNull().And("age").Greater(17)
```

</details>

<details>
<summary>Columns</summary>

```go
// When you have a column built by the Column Builder,
// you can also modify it afterward.

// samples
var (
    age   = c.Col("age")
    stock = c.Col("stock")
    price = c.Col("price")
    fname = c.Col("firstname")
)

// ----- Aggregation ------

// You can call aggregation methods such as Count(), Sum(), Max(), Min(), or Avg().
// Only one aggregation method will be applied, and it will always be the last one called.
// Each method accepts an optional bool that enables DISTINCT inside the aggregation.
age.Count()   // COUNT(age)
age.Sum(true) // SUM(DISTINCT age)
age.Max()     // MAX(age)
age.Min(true) // MIN(DISTINCT age)
age.Avg()     // AVG(age)

// ----- Arithmetic operations -----

// You can perform arithmetic operations with other columns or literal values.
price.Add(10) // price + 10
price.Sub(10) // price - 10
price.Mul(10) // price * 10
price.Div(10) // price / 10
price.Mod(10) // price % 10

stock.Mul(price) // stock * price

// If you want to wrap an expression, use Wrap()
stock.Mul(price).Wrap().Mul(100) // (stock * price) * 100

// ----- Scalar functions ------

// You can use scalar functions that modify the resulting expression.
fname.Upper() // UPPER(firstname)
fname.Lower() // LOWER(firstname)
fname.Trim()  // TRIM(firstname)
price.Round() // ROUND(price)
age.Avg()     // AVG(age)
age.Abs()     // ABS(age)

// ----- Alias -----

// Some methods change the column’s default name.
// To preserve control over the final name, use As() to assign an alias.
age.Avg().As("avg_age") // AVG(age) AS avg_age

// ----- Written methods ------

// Columns can be written differently depending on the clause.
// If you want to force a specific written form, you can use:

// sample:
age.Count().As("avg_age") // (COUNT(age) AS avg_age)

age.Base()  // age
age.Expr()  // COUNT(age)
age.Alias() // avg_age
age.Def()   // COUNT(age) AS avg_age
```

</details>

<details>
<summary>Select statement</summary>

```go
// Start a new SELECT statement.
// Returns (q) the statement and (c) the column builder.
// Requires a context and a default table, which will be used
// as the default table for the column builder and the Auto From() clause.
q, c := qe.NewSelect(ctx, "users")

// Sample columns
var (
	id       = c.Col("id")
	fk       = c.Col("user_id", "posts")
	age      = c.Col("age")
	banned   = c.Col("banned")
	example  = c.Col("example")
	postsNum = c.Col("id", "posts").Count()
)

// ----- SELECT -----

// Accepts a variadic parameter; you can select any columns you want.
q.Select(c.All()) // SELECT *

// You can also start with Distinct() for this purpose.
q.Distinct().Select(id) // SELECT DISTINCT id

// ----- FROM -----

// From() selects the target table.
// If omitted, the default table is used automatically.
q.From("users")

// ----- JOIN -----

// Always start a JOIN clause with Join().
// You can chain On() directly to set the condition, and it will default to INNER JOIN.
q.Join("posts").On(fk).Equal(id)

// You can also explicitly set the join type.
q.Join("posts").Left().On(fk).Equal(id)
q.Join("posts").Full().On(fk).Equal(id)

// CROSS JOIN does not accept a condition.
q.Join("posts").Cross()

// ----- WHERE -----

// Where() starts a condition. Each condition is wrapped in parentheses.
q.Where(age).Greater(17)                     // WHERE (age > 17)
q.Where(age).Less(40).And(banned).IsNull()  // WHERE (age < 40 AND banned IS NULL)

// Multiple Where() calls are combined with AND at the top level.
// Example: WHERE (age > 17) AND (age < 40 AND banned IS NULL)

// You can negate any condition using Not()
q.Where(example).Not().Equal(100)          // WHERE (example <> 100)
q.Where(example).Not().Like("%something%") // WHERE (example NOT LIKE '%something%')
q.Where(example).Not().Greater(40)         // WHERE (example <= 40)
q.Where(example).Not().IsNull()            // WHERE (example IS NOT NULL)

// You can also use parameters for external inputs.
q.Where(id).Equal(c.Param("uuid")) // WHERE (id = ?)

// ----- GROUP BY -----

// GROUP BY behaves similarly to the SELECT clause.
// You can pass multiple columns as variadic parameters.
q.GroupBy(postsNum)

// ----- HAVING -----

// HAVING behaves like WHERE. Start with an identifier to define a condition.
q.Having(postsNum).Greater(10)

// ----- ORDER BY -----

// ORDER BY accepts an identifier.
q.OrderBy(age)

// You can chain Desc() for descending order.
// Multiple OrderBy() calls are allowed.
q.OrderBy(age).Desc()
q.OrderBy(id)

// ----- LIMIT & OFFSET -----

// Pass an integer to set the limit or offset.
q.Limit(10)
q.Offset(5)
```

</details>
