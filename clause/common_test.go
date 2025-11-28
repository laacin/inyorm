package clause_test

import (
	"reflect"
	"testing"

	"github.com/laacin/inyorm"
	"github.com/laacin/inyorm/clause"
	"github.com/laacin/inyorm/internal/column"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/writer"
)

type (
	runner = func(t *testing.T, cls string, vals []any)

	colBuilder = column.ColBuilder[
		inyorm.Column, inyorm.Condition, inyorm.ConditionNext,
		inyorm.Case, inyorm.CaseNext,
	]

	clsSelect  = clause.Select[inyorm.SelectNext]
	clsJoin    = clause.Join[inyorm.JoinNext, inyorm.JoinEnd, inyorm.Condition, inyorm.ConditionNext]
	clsWhere   = clause.Where[inyorm.Condition, inyorm.ConditionNext]
	clsGroupBy = clause.GroupBy
	clsHaving  = clause.Having[inyorm.Condition, inyorm.ConditionNext]
	clsOrderBy = clause.OrderBy[inyorm.OrderByNext]
)

func New[T any](cls core.Clause, dialect []string) (T, inyorm.ColumnBuilder, runner) {
	d := ""
	if len(dialect) > 0 {
		d = dialect[0]
	}
	tbl := "users"

	q := writer.Query{Config: &core.Config{
		Dialect:   d,
		ColumnTag: core.DefaultColumnTag,
		ColWrite:  core.DefaultColumnWriter,
	}}

	q.PreBuild(func(cfg *core.Config) (useAliases bool) {
		return cls.Name() == "JOIN"
	})

	c := &colBuilder{Table: tbl}
	run := func(t *testing.T, expect string, vals []any) {
		q.SetClauses([]core.Clause{cls})
		statement, values, err := q.Build()
		if err != nil {
			t.Fatal(err)
		}

		if statement != expect {
			if len(statement) == len(expect) {
				for i := range statement {
					have, exp := statement[i], expect[i]
					if exp != have && exp != 'x' {
						t.Errorf("\nmismatch result:\nExpect:\n%s\nHave:\n%s\n", expect, statement)
						return
					}
				}
			} else {
				t.Errorf("\nmismatch result:\nExpect:\n%s\nHave:\n%s\n", expect, statement)
			}
		}
		if !reflect.DeepEqual(vals, values) {
			t.Errorf("\nmissmatch values:\nExpect:\n%#v\nHave:\n%#v\n", vals, values)
		}
	}

	return any(cls).(T), c, run
}

func NewSelect(dialect ...string) (inyorm.Select, inyorm.ColumnBuilder, runner) {
	return New[inyorm.Select](&clsSelect{}, dialect)
}

func NewJoin(dialect ...string) (inyorm.Join, inyorm.ColumnBuilder, runner) {
	return New[inyorm.Join](&clsJoin{}, dialect)
}

func NewWhere(dialect ...string) (inyorm.Where, inyorm.ColumnBuilder, runner) {
	return New[inyorm.Where](&clsWhere{}, dialect)
}

func NewGroupBy(dialect ...string) (inyorm.GroupBy, inyorm.ColumnBuilder, runner) {
	return New[inyorm.GroupBy](&clsGroupBy{}, dialect)
}

func NewHaving(dialect ...string) (inyorm.Having, inyorm.ColumnBuilder, runner) {
	return New[inyorm.Having](&clsHaving{}, dialect)
}

func NewOrderBy(dialect ...string) (inyorm.OrderBy, inyorm.ColumnBuilder, runner) {
	return New[inyorm.OrderBy](&clsOrderBy{}, dialect)
}
