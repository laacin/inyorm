package expr

type sqlOp string

const (
	equal    sqlOp = "="
	notEqual sqlOp = "<>"

	greater    sqlOp = ">"
	notGreater sqlOp = "<="

	less    sqlOp = "<"
	notLess sqlOp = ">="

	and sqlOp = "AND"
	or  sqlOp = "OR"

	in    sqlOp = "IN"
	notIn sqlOp = "NOT IN"

	between    sqlOp = "BETWEEN"
	notBetween sqlOp = "NOT BETWEEN"

	isNull    sqlOp = "IS NULL"
	isNotNull sqlOp = "IS NOT NULL"

	like    sqlOp = "LIKE"
	notLike sqlOp = "NOT LIKE"
)

var negations = map[sqlOp]sqlOp{
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
