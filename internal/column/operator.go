package column

// Aggregation constant
const (
	countAggr = "COUNT"
	sumAggr   = "SUM"
	minAggr   = "MIN"
	maxAggr   = "MAX"
	avgAggr   = "AVG"
)

// Arithmetical operation
const (
	addOp byte = '+'
	subOp byte = '-'
	mulOp byte = '*'
	divOp byte = '/'
	modOp byte = '%'
)

// Scalar function
const (
	lowerFunc = "LOWER"
	upperFunc = "UPPER"
	trimFunc  = "TRIM"
	roundFunc = "ROUND"
	absFunc   = "ABS"
)
