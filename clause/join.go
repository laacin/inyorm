package clause

import (
	"fmt"
	"strings"
)

type JoinTyp string

const (
	InnerJoin JoinTyp = "INNER"
	LeftJoin  JoinTyp = "LEFT"
	RightJoin JoinTyp = "RIGHT"
	FullJoin  JoinTyp = "FULL"
	CrossJoin JoinTyp = "CROSS"
)

func (typ JoinTyp) IsValid() bool {
	switch typ {
	case InnerJoin, LeftJoin, RightJoin, FullJoin, CrossJoin:
		return true
	default:
		return false
	}
}

type JoinBuilder struct {
	table      string
	alias      string
	key        string
	joins      []*Joined
	interJoins []*InterJoin
}

func NewJoinBuilder(table, alias, key string) *JoinBuilder {
	return &JoinBuilder{table: table, alias: alias, key: key}
}

type InterJoin struct {
	typ   JoinTyp
	table string
	alias string
	keys  map[string]string
	joins []*Joined
}

type Joined struct {
	typ        JoinTyp
	table      string
	alias      string
	foreignKey string
}

func (j *JoinBuilder) Join(typ JoinTyp, table, alias, foreignKey string) {
	j.joins = append(j.joins, newJoin(typ, table, alias, foreignKey))
}

func (j *JoinBuilder) Many(typ JoinTyp, table, alias string, keys map[string]string) *InterJoin {
	inter := &InterJoin{
		typ:   typ,
		table: table,
		alias: alias,
		keys:  keys,
	}
	j.interJoins = append(j.interJoins, inter)
	return inter
}

func (j *InterJoin) Join(typ JoinTyp, table, alias, foreignKey string) {
	j.joins = append(j.joins, newJoin(typ, table, alias, foreignKey))
}

// --- Helper
func newJoin(typ JoinTyp, table, alias, fk string) *Joined {
	return &Joined{
		typ:        typ,
		table:      table,
		alias:      alias,
		foreignKey: fk,
	}
}

func (j *JoinBuilder) Build(sb *strings.Builder) []error {
	var errs []error
	manyToMany := false

	for i, inter := range j.interJoins {
		fk, err := validateJoin(inter.typ, j.table, inter.table, inter.keys)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		if i > 0 {
			sb.WriteByte(' ')
		}

		manyToMany = true
		writeJoin(sb, j.alias, j.key, &Joined{
			typ:        inter.typ,
			table:      inter.table,
			alias:      inter.alias,
			foreignKey: fk,
		})

		for _, join := range inter.joins {
			joinedFk, err := validateJoin(join.typ, join.table, inter.table, inter.keys)
			if err != nil {
				errs = append(errs, err)
				continue
			}

			sb.WriteByte(' ')
			writeJoin(sb, inter.alias, joinedFk, join)
		}

	}

	for i, join := range j.joins {
		if _, err := validateJoin(join.typ, "", "", nil); err != nil {
			errs = append(errs, err)
			continue
		}

		if i > 0 || manyToMany {
			sb.WriteByte(' ')
		}
		writeJoin(sb, j.alias, j.key, join)
	}
	return errs
}

func writeJoin(sb *strings.Builder, alias, key string, joined *Joined) {
	sb.WriteString(string(joined.typ))
	sb.WriteString(" JOIN ")
	sb.WriteString(joined.table)

	if joined.typ != CrossJoin {
		sb.WriteByte(' ')
		sb.WriteString(joined.alias)
		sb.WriteString(" ON ")
		sb.WriteString(joined.alias)
		sb.WriteByte('.')
		sb.WriteString(joined.foreignKey)
		sb.WriteString(" = ")
		sb.WriteString(alias)
		sb.WriteByte('.')
		sb.WriteString(key)
	}
}

func validateJoin(typ JoinTyp, join, on string, keys map[string]string) (string, error) {
	if !typ.IsValid() {
		return "", fmt.Errorf("'%s' Join type is not valid", typ)
	}

	if keys != nil {
		fk, exists := keys[join]
		if !exists {
			return "", fmt.Errorf("'%s' table has no relation with '%s' table", join, on)
		}
		return fk, nil
	}

	return "", nil
}
