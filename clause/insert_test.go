package clause_test

import "testing"

func TestInsert(t *testing.T) {
	type User struct {
		Account string `col:"account"`
		Age     int    `col:"age"`
		Active  bool   `col:"active"`
		Score   int    `col:"score"`
		Country string `col:"country"`
	}

	writeVals := func(exp *string, rows, perRow int) {
		*exp += "VALUES "
		for i := range rows {
			if i > 0 {
				*exp += ", "
			}
			*exp += "("
			for pi := range perRow {
				if pi > 0 {
					*exp += ", "
				}
				*exp += "?"
			}
			*exp += ")"
		}
	}

	t.Run("basic", func(t *testing.T) {
		cls, _, run := NewInsert("users")

		cls.Insert(User{}).Values(User{
			Account: "acc1",
			Age:     26,
			Active:  true,
			Score:   200,
			Country: "AR",
		})

		exp := "INSERT INTO users (account, active, age, country, score) "
		exp += "VALUES (?, ?, ?, ?, ?)"
		run(t, exp, []any{"acc1", true, 26, "AR", 200})
	})

	t.Run("many", func(t *testing.T) {
		cls, _, run := NewInsert("users")

		cls.Insert(User{}).Values([]User{
			{Account: "acc1", Age: 26, Active: true, Score: 200, Country: "AR"},
			{Account: "acc2", Age: 30, Active: false, Score: 150, Country: "US"},
			{Account: "acc3", Age: 22, Active: true, Score: 320, Country: "BR"},
			{Account: "acc4", Age: 41, Active: false, Score: 400, Country: "DE"},
			{Account: "acc5", Age: 19, Active: true, Score: 180, Country: "UK"},
		})

		vals := []any{
			"acc1", true, 26, "AR", 200,
			"acc2", false, 30, "US", 150,
			"acc3", true, 22, "BR", 320,
			"acc4", false, 41, "DE", 400,
			"acc5", true, 19, "UK", 180,
		}

		exp := "INSERT INTO users (account, active, age, country, score) "
		writeVals(&exp, 5, 5)
		run(t, exp, vals)
	})

	t.Run("ptr", func(t *testing.T) {
		cls, _, run := NewInsert("users")

		cls.Insert(&User{}).Values(&User{
			Account: "acc1",
			Age:     26,
			Active:  true,
			Score:   200,
			Country: "AR",
		})

		exp := "INSERT INTO users (account, active, age, country, score) "
		exp += "VALUES (?, ?, ?, ?, ?)"
		run(t, exp, []any{"acc1", true, 26, "AR", 200})
	})

	t.Run("many_ptr", func(t *testing.T) {
		cls, _, run := NewInsert("users")

		cls.Insert(&User{}).Values(&[]User{
			{Account: "acc1", Age: 26, Active: true, Score: 200, Country: "AR"},
			{Account: "acc2", Age: 30, Active: false, Score: 150, Country: "US"},
			{Account: "acc3", Age: 22, Active: true, Score: 320, Country: "BR"},
			{Account: "acc4", Age: 41, Active: false, Score: 400, Country: "DE"},
			{Account: "acc5", Age: 19, Active: true, Score: 180, Country: "UK"},
		})

		vals := []any{
			"acc1", true, 26, "AR", 200,
			"acc2", false, 30, "US", 150,
			"acc3", true, 22, "BR", 320,
			"acc4", false, 41, "DE", 400,
			"acc5", true, 19, "UK", 180,
		}

		exp := "INSERT INTO users (account, active, age, country, score) "
		writeVals(&exp, 5, 5)
		run(t, exp, vals)
	})

	t.Run("literals", func(t *testing.T) {
		cls, c, run := NewInsert("users")

		cls.Insert(c.Col("age")).Values(22)

		exp := "INSERT INTO users (age) "
		exp += "VALUES (?)"
		run(t, exp, []any{22})
	})

	t.Run("literals_slc", func(t *testing.T) {
		cls, c, run := NewInsert("users")

		cls.Insert(c.Col("age"), c.Col("name")).Values([]any{22, "marie", 54, "robert"})

		exp := "INSERT INTO users (age, name) "
		writeVals(&exp, 2, 2)
		run(t, exp, []any{22, "marie", 54, "robert"})
	})

	t.Run("map", func(t *testing.T) {
		cls, _, run := NewInsert("users")

		m := map[string]any{
			"account": "acc1", "age": 26, "active": true, "score": 200, "country": "AR",
		}

		cls.Insert(m).Values(m)

		exp := "INSERT INTO users (account, active, age, country, score) "
		exp += "VALUES (?, ?, ?, ?, ?)"
		run(t, exp, []any{"acc1", true, 26, "AR", 200})
	})

	t.Run("many_map", func(t *testing.T) {
		cls, _, run := NewInsert("users")

		m := []map[string]any{
			{"account": "acc1", "age": 26, "active": true, "score": 200, "country": "AR"},
			{"account": "acc2", "age": 30, "active": false, "score": 150, "country": "US"},
			{"account": "acc3", "age": 22, "active": true, "score": 320, "country": "BR"},
			{"account": "acc4", "age": 41, "active": false, "score": 400, "country": "DE"},
			{"account": "acc5", "age": 19, "active": true, "score": 180, "country": "UK"},
		}

		cls.Insert(m[0]).Values(m)

		vals := []any{
			"acc1", true, 26, "AR", 200,
			"acc2", false, 30, "US", 150,
			"acc3", true, 22, "BR", 320,
			"acc4", false, 41, "DE", 400,
			"acc5", true, 19, "UK", 180,
		}

		exp := "INSERT INTO users (account, active, age, country, score) "
		writeVals(&exp, 5, 5)
		run(t, exp, vals)
	})

	t.Run("map_ptr", func(t *testing.T) {
		cls, _, run := NewInsert("users")

		m := map[string]any{
			"account": "acc1", "age": 26, "active": true, "score": 200, "country": "AR",
		}

		cls.Insert(&m).Values(&m)

		exp := "INSERT INTO users (account, active, age, country, score) "
		exp += "VALUES (?, ?, ?, ?, ?)"
		run(t, exp, []any{"acc1", true, 26, "AR", 200})
	})

	t.Run("many_map_ptr", func(t *testing.T) {
		cls, _, run := NewInsert("users")

		m := []map[string]any{
			{"account": "acc1", "age": 26, "active": true, "score": 200, "country": "AR"},
			{"account": "acc2", "age": 30, "active": false, "score": 150, "country": "US"},
			{"account": "acc3", "age": 22, "active": true, "score": 320, "country": "BR"},
			{"account": "acc4", "age": 41, "active": false, "score": 400, "country": "DE"},
			{"account": "acc5", "age": 19, "active": true, "score": 180, "country": "UK"},
		}

		cls.Insert(m[0]).Values(m)

		vals := []any{
			"acc1", true, 26, "AR", 200,
			"acc2", false, 30, "US", 150,
			"acc3", true, 22, "BR", 320,
			"acc4", false, 41, "DE", 400,
			"acc5", true, 19, "UK", 180,
		}

		exp := "INSERT INTO users (account, active, age, country, score) "
		writeVals(&exp, 5, 5)
		run(t, exp, vals)
	})

	t.Run("omit_values", func(t *testing.T) {
		cls, c, run := NewInsert("users")

		cls.Insert(c.Col("score"), c.Col("country")).Values([]User{
			{Account: "acc1", Age: 26, Active: true, Score: 200, Country: "AR"},
			{Account: "acc2", Age: 30, Active: false, Score: 150, Country: "US"},
			{Account: "acc3", Age: 22, Active: true, Score: 320, Country: "BR"},
			{Account: "acc4", Age: 41, Active: false, Score: 400, Country: "DE"},
			{Account: "acc5", Age: 19, Active: true, Score: 180, Country: "UK"},
		})

		vals := []any{
			"AR", 200,
			"US", 150,
			"BR", 320,
			"DE", 400,
			"UK", 180,
		}

		exp := "INSERT INTO users (country, score) "
		writeVals(&exp, 5, 2)
		run(t, exp, vals)
	})

	t.Run("partially_filled_struct", func(t *testing.T) {
		cls, _, run := NewInsert("users")

		cls.Insert(User{}).Values(User{Account: "acc1"})
		exp := "INSERT INTO users (account, active, age, country, score) VALUES (?, ?, ?, ?, ?)"
		run(t, exp, []any{"acc1", false, 0, "", 0})
	})

	t.Run("nested_struct", func(t *testing.T) {
		type Nested struct {
			User
			Extra string `col:"extra"`
		}

		cls, _, run := NewInsert("users")
		cls.Insert(Nested{}).Values(Nested{
			User:  User{Account: "acc1", Age: 26, Active: true, Score: 200, Country: "AR"},
			Extra: "X",
		})

		exp := "INSERT INTO users (account, active, age, country, extra, score) VALUES (?, ?, ?, ?, ?, ?)"
		run(t, exp, []any{"acc1", true, 26, "AR", "X", 200})
	})
}
