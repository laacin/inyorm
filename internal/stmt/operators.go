package stmt

type Operator string

const (
	Equal      Operator = "="
	NotEqual   Operator = "<>"
	Greater    Operator = ">"
	NotGreater Operator = "<="
	Less       Operator = "<"
	NotLess    Operator = ">="

	And Operator = "AND"
	Or  Operator = "OR"
	Not Operator = "NOT"

	In         Operator = "IN"
	NotIn      Operator = "NOT IN"
	Between    Operator = "BETWEEN"
	NotBetween Operator = "NOT BETWEEN"

	IsNull    Operator = "IS NULL"
	IsNotNull Operator = "IS NOT NULL"

	Like    Operator = "LIKE"
	NotLike Operator = "NOT LIKE"
)

var negations = map[Operator]Operator{
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

func GetOp(kind Operator, negated bool) string {
	if negated {
		if op, ok := negations[kind]; ok {
			return string(op)
		}
	}

	return string(kind)
}
