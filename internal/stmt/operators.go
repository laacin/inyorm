package stmt

type SqlOp string

const (
	Equal    SqlOp = "="
	NotEqual SqlOp = "<>"

	Greater    SqlOp = ">"
	NotGreater SqlOp = "<="

	Less    SqlOp = "<"
	NotLess SqlOp = ">="

	And SqlOp = "AND"
	Or  SqlOp = "OR"
	Not SqlOp = "NOT"

	In    SqlOp = "IN"
	NotIn SqlOp = "NOT IN"

	Between    SqlOp = "BETWEEN"
	NotBetween SqlOp = "NOT BETWEEN"

	IsNull    SqlOp = "IS NULL"
	IsNotNull SqlOp = "IS NOT NULL"

	Like    SqlOp = "LIKE"
	NotLike SqlOp = "NOT LIKE"
)

var negations = map[SqlOp]SqlOp{
	Equal:      NotEqual,
	NotEqual:   Equal,
	Greater:    NotGreater,
	NotGreater: Greater,
	Less:       NotLess,
	NotLess:    Less,
	In:         NotIn,
	NotIn:      In,
	Between:    NotBetween,
	NotBetween: Between,
	IsNull:     IsNotNull,
	IsNotNull:  IsNull,
	Like:       NotLike,
	NotLike:    Like,
}

func GetSqlOp(kind SqlOp, negated bool) string {
	if negated {
		if op, ok := negations[kind]; ok {
			return string(op)
		}
	}

	return string(kind)
}
