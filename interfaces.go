package inyorm

//	type Column interface {
//		Switch(target any, fn func(w Case[string]))
//		Search(target any) Case[func(Expression)]
//	}

// type Case interface {
// 	When(value any) Then
// 	Else(do any)
// }
//
// type Search interface {
// 	When(values ...any) *Expression
// 	Else(do any)
// }
//
// type Then interface {
// 	Then(do any) Case
// }
//
// type FieldT interface {
// 	Set() string
//
// 	Greater(any) Field
// 	Less(any) Field
//
// 	Add(any) Field
// 	Sub(any) Field
// 	Mul(any) Field
// 	Div(any) Field
// 	Mod(any) Field
// 	Get() Field
//
// 	Substring(start, end int) Field
// 	Upper() Field
// 	Lower() Field
// 	Trim() Field
// }
//
// type FieldBuilder interface {
// 	Target(values ...string) Field
// 	Switch(literal string, fn func(c Case)) Field
// 	Search(func(c Case, e *Expression)) Field
// }
//
// type
//
