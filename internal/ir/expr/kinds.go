package expr

type ExprKind int

const (
	// Literals
	ExprString ExprKind = iota
	ExprNumber
	ExprFloat
	ExprBool
	ExprNull

	// Specials
	ExprParameter
	ExprCondition
	ExprConcat
	ExprCaseSwitch
	ExprCaseSearch

	// SQL Values
	ExprTable
	ExprColumn
)
