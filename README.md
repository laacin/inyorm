# Inyorm

[![Go Reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/laacin/inyorm)

##### Inyorm is a fully declarative ORM for Go, designed for clarity, type safety, and predictable SQL generation.

- SQL-like declarative API
- Strong typing
- Fast mapping with minimal reflection
- Lightweight, explicit, and predictable
- Highly customizable

## Minimal setup

### Install

```bash
go get -u github.com/laacin/inyorm
```

Inyorm is just a query builder and object mapper. It doesn’t manage connections or drivers,
so you stay fully in control of the database layer.

```go
package main

import (
	"database/sql"
	"log"

	"github.com/laacin/inyorm"
)

func main() {
    ctx := context.Background()

    db, err := sql.Open("sqlite3", "./data.db")
    if err != nil {
        log.Fatal(err)
    }

    qe := inyorm.New("sqlite3", db, &inyorm.Options{})

    q, c := qe.NewSelect(ctx, "table")
    // q, c := qe.NewInsert(ctx, "table")
    // q, c := qe.NewUpdate(ctx, "table")
    // q, c := qe.NewDelete(ctx, "table")
}
```

## Why inyorm?

Inyorm is a fully declarative, explicit, and strongly typed ORM,
letting you write complex queries without falling back to raw SQL.

### Examples:

<details>
<summary>Simple</summary>

```go
q, c := qe.NewSelect(ctx, "users")

var (
	id    = c.Col("id")
	fk    = c.Col("user_id", "posts")
	posts = c.Col("id", "posts").Count()
)

q.Select(c.All(), posts)
q.Join("posts").On(fk).Equal(id)
q.Where(id).Equal(c.Param("uuid"))
q.Limit(1)
```

SQL:

```sql
SELECT a.*, COUNT(b.id)
FROM users a
INNER JOIN posts b ON (b.user_id = a.id)
WHERE (a.id = ?)
LIMIT 1;
```

</details>

<details>
<summary>Complex</summary>

```go
q, c := qe.NewSelect(context.Background(), "users")

// helper
isNull := func(cond inyorm.Column, then, els any) inyorm.Column {
	return c.Search(func(cs inyorm.Case) {
		cond := c.Cond(cond).IsNull()
		cs.When(cond).Then(then)
		cs.Else(els)
	})
}

var (
	// base columns
	fname     = c.Col("firstname")
	lname     = c.Col("lastname")
	lastLogin = c.Col("last_login")
	banned    = c.Col("banned")
	roleName  = c.Col("name", "roles")
	age       = c.Col("age")
	active    = c.Col("active")

	// summary columns
	posts    = c.Col("id", "posts").Count()
	comments = c.Col("id", "comments").Count()
	role     = isNull(roleName, "No role", roleName)
	lastLog  = isNull(lastLogin, "Never", lastLogin)
	status   = isNull(banned, "Active", c.Concat("Banned at: ", banned))

	// join keys
	userId        = c.Col("id")
	roleId        = c.Col("id", "roles")
	postUserFk    = c.Col("user_id", "posts")
	commentUserFk = c.Col("user_id", "comments")

	interUserFk = c.Col("user_id", "user_roles")
	interRoleFk = c.Col("role_id", "user_roles")

	// select columns
	totalPost    = c.Col("id", "posts").Count(true).As("total_posts")
	totalComment = c.Col("id", "comments").Count(true).As("total_comments")
	lastPost     = c.Col("created_at", "posts").Max().As("last_post_date")
	summary      = c.Concat(
		"User: ", fname, " ", lname, " | ",
		"Role: ", role, " | ",
		"Posts: ", posts, " | ",
		"Comments: ", comments, " | ",
		"Last login: ", lastLog, " | ",
		status,
	).As("user_summary")
)

// statement building
q.Select(summary, totalPost, totalComment, lastPost)

q.Join("user_roles").Left().On(interUserFk).Equal(userId)
q.Join("roles").Left().On(roleId).Equal(interRoleFk)
q.Join("posts").Left().On(postUserFk).Equal(userId)
q.Join("comments").Left().On(commentUserFk).Equal(userId)

q.Where(age).Between(18, 60).And(active).Equal(true)

q.GroupBy(userId, fname, lname, lastLogin, banned, roleName)
q.Having(posts).Greater(5)

q.OrderBy(totalPost).Desc()
q.OrderBy(age)
q.Limit(50)
```

SQL:

```sql
SELECT CONCAT(
        'User: ', a.firstname, ' ', a.lastname, ' | ',
        'Role: ', CASE WHEN (b.name IS NULL) THEN 'No role' ELSE b.name END, ' | ',
        'Posts: ', COUNT(c.id), ' | ',
        'Comments: ', COUNT(d.id), ' | ',
        'Last login: ', CASE WHEN (a.last_login IS NULL) THEN 'Never' ELSE a.last_login END, ' | ',
        CASE WHEN (a.banned IS NULL) THEN 'Active' ELSE CONCAT('Banned at: ', a.banned) END
    ) AS user_summary,
    COUNT(DISTINCT c.id) AS total_posts,
    COUNT(DISTINCT d.id) AS total_comments,
    MAX(c.created_at) AS last_post_date
FROM users a
LEFT JOIN user_roles e ON (e.user_id = a.id)
LEFT JOIN roles b ON (b.id = e.role_id)
LEFT JOIN posts c ON (c.user_id = a.id)
LEFT JOIN comments d ON (d.user_id = a.id)
WHERE (a.age BETWEEN 18 AND 60 AND a.active = 1)
GROUP BY a.id, a.firstname, a.lastname, a.last_login, a.banned, b.name
HAVING (COUNT(c.id) > 5)
ORDER BY total_posts DESC, a.age
LIMIT 50;
```

</details>

## Guides

<details>
<summary>Column builder</summary>

```go
// The second return value of each new statement is the Column Builder,
// where you can write all non-literal values.
_, c := qe.NewSelect(ctx, "table")

// ----- Col -----

// Col is the most common method. It references a table column
// and accepts two parameters: the first is column name of the main table.
c.Col("id")

// To reference another table, pass a second parameter with the table name.
c.Col("id", "posts")

// ----- All -----

// All references the wildcard '*'.
// In a joined query, the default All() references the main table.
c.All()

// To reference another table, pass a parameter with the table name.
c.All("posts")

// ----- Param -----

// Param writes a placeholder.
// In Inyorm you should write explicit parameters (except in explicit values clause).
// such as Insert().Values() and Update().Values()
c.Param("id")

// You can also skip parameters for lazy values, useful for prepared statements.
c.Param()

// ----- Concat -----

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
// such as a condition (built below) or literals.
c.Search(func(cs inyorm.Case) {
	cond := c.Cond(c.Col("age")).Greater(17)
	cs.When(cond).Then("Adult")
	cs.Else("Kid")
})

// Expected output:
// CASE WHEN (age > 17) THEN 'Adult' ELSE 'Kid' END

// ----- Cond -----

// Cond is a helper used to create a condition.
// Currently, only used in the Search method.
c.Cond(c.Col("banned")).IsNull().And(c.Col("age")).Greater(17)
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
<summary>Select</summary>

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

<details>
<summary>Insert & Update</summary>

```go

// ----- Insert & Update -----

// Both share the same logic.
// I would use Insert() for the examples, but both are the same.

type User struct {
	Email    string `col:"account"`
	Password string `col:"password"`
	Age      int    `col:"age"`
}

// We have Insert() and Update(). Both share the same purpose, which is
// to receive an example of which columns will be affected.
// Both methods return another method, Values(),
// where the real values are provided. You can ignore this method for prepared statements.
q.Insert().Values()
q.Update().Values()

// For Insert/Update you can pass different values:

// A struct – it will read the tags and use each one.
// In this case, EACH field that has the correct tag will be read,
// even if it doesn't have values.
q.Insert(User{}).Values(User{
	Email:    "example@email.com",
	Password: "1234",
	Age:      34,
})

// Columns –
// Only the specified ones will be inserted/updated, even if the struct/map has more fields.
// This is the best way to limit the columns.
// In this case:
// - password and age will be inserted/updated and email ignored.
q.Insert(c.Col("password"), c.Col("age")).Values(User{
	Email:    "example@email.com",
	Password: "1234",
	Age:      34,
})

// Map, in which EVERY key of the map will be one column.
m := map[string]any{
	"email":    "example@email.com",
	"password": "1234",
	"age":      34,
}
q.Insert(m).Values(m)
```

</details>
