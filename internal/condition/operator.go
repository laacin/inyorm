package condition

const (
	equal    = "="
	notEqual = "<>"

	greater    = ">"
	notGreater = "<="

	less    = "<"
	notLess = ">="

	in    = "IN"
	notIn = "NOT IN"

	between    = "BETWEEN"
	notBetween = "NOT BETWEEN"

	isNull    = "IS NULL"
	isNotNull = "IS NOT NULL"

	like    = "LIKE"
	notLike = "NOT LIKE"

	and = "AND"
	or  = "OR"
)

var negations = map[string]string{
	equal:    notEqual,
	notEqual: equal,

	greater:    notGreater,
	notGreater: greater,

	less:    notLess,
	notLess: less,

	in:    notIn,
	notIn: in,

	between:    notBetween,
	notBetween: between,

	isNull:    isNotNull,
	isNotNull: isNull,

	like:    notLike,
	notLike: like,
}

func getOp(kind string, negated bool) string {
	if negated {
		return negations[kind]
	}
	return kind
}
