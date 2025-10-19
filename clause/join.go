package clause

import (
	"fmt"
	"strings"
)

// TODO: better use for aliases

// -- JOIN TYPES
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

// ----- BUILDER ------

type JoinBuilder struct {
	Table        string
	PrimaryKey   string
	joins        []*Join
	intermediate []*JoinIntermediate
}

func (j *JoinBuilder) Simple(typ JoinTyp, table, foreignKey string) {
	join := &Join{typ, table, foreignKey}
	j.joins = append(j.joins, join)
}

// keys is a map which contains table reference (key) and his foreign key
func (j *JoinBuilder) ManyToMany(typ JoinTyp, relTable string, keys map[string]string) *JoinIntermediate {
	inter := &JoinIntermediate{typ: typ, table: relTable, keys: keys}
	j.intermediate = append(j.intermediate, inter)
	return inter
}

// write: type JOIN joinedTable joinAlias
// ON joinedAlias.joinedFk = mainAlias.mainPk ...
func (j *JoinBuilder) Build(sb *strings.Builder) []error {
	var errs []error
	hasRelation := false

	for i, rel := range j.intermediate {
		if !rel.typ.IsValid() {
			errs = append(errs, fmt.Errorf("'%s' Join type is not valid", rel.typ))
			continue
		}

		fk, exists := rel.keys[j.Table]
		if !exists {
			errs = append(errs, fmt.Errorf("intermediate '%s' table has no relation with the main table", rel.table))
			continue
		}

		if i > 0 {
			sb.WriteByte(' ')
		}

		hasRelation = true
		writeJoin(sb, rel.typ, j.Table, j.PrimaryKey, rel.table, fk)

		for _, join := range rel.joins {
			if !join.typ.IsValid() {
				errs = append(errs, fmt.Errorf("'%s' Join type is not valid", join.typ))
				continue
			}

			joinedFk, exists := rel.keys[join.table]
			if !exists {
				errs = append(errs, fmt.Errorf("intermediate '%s' table has no relation with '%s' table", rel.table, join.table))
				continue
			}

			sb.WriteByte(' ')
			writeJoin(sb, join.typ, rel.table, joinedFk, join.table, join.key)
		}

	}

	for i, join := range j.joins {
		if !join.typ.IsValid() {
			errs = append(errs, fmt.Errorf("'%s' Join type is not valid", join.typ))
			continue
		}

		if i > 0 || hasRelation {
			sb.WriteByte(' ')
		}
		writeJoin(sb, join.typ, j.Table, j.PrimaryKey, join.table, join.key)
	}
	return errs
}

type JoinIntermediate struct {
	typ   JoinTyp
	table string
	keys  map[string]string
	joins []*Join
}

func (j *JoinIntermediate) Join(typ JoinTyp, table, key string) *JoinIntermediate {
	j.joins = append(j.joins, &Join{typ, table, key})
	return j
}

type Join struct {
	typ   JoinTyp
	table string
	key   string
}

func writeJoin(sb *strings.Builder, typ JoinTyp, srcTable, srcKey, joinTable, joinKey string) {
	sb.WriteString(string(typ))
	sb.WriteString(" JOIN ")
	sb.WriteString(joinTable)

	if typ != CrossJoin {
		sb.WriteByte(' ')
		sb.WriteString(joinTable)
		sb.WriteString(" ON ")
		sb.WriteString(joinTable)
		sb.WriteByte('.')
		sb.WriteString(joinKey)
		sb.WriteString(" = ")
		sb.WriteString(srcTable)
		sb.WriteByte('.')
		sb.WriteString(srcKey)
	}
}
