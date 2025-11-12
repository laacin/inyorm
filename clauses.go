package inyorm

import "github.com/laacin/inyorm/clause"

type (
	selectCls struct{ selectWrap *clause.Select }
	fromCls   struct{ fromWrap *clause.From }

	joinCls struct {
		joinWrap *clause.Join
		joinNext *joinNext
	}
	joinNext struct {
		ctx  *joinCls
		wrap **clause.Join
	}

	whereCls   struct{ whereWrap *clause.Where }
	groupByCls struct{ groupByWrap *clause.GroupBy }
	havingCls  struct{ havingWrap *clause.Having }

	orderByCls struct {
		orderByWrap *clause.OrderBy
		orderByNext *orderByNext
	}
	orderByNext struct {
		wrap **clause.OrderBy
		ctx  *orderByCls
	}

	limitCls  struct{ limitWrap *clause.Limit }
	offsetCls struct{ offsetWrap *clause.Offset }
)

// ----- Constructors -----

func wrapSelect() *selectCls   { return &selectCls{} }
func wrapFrom() *fromCls       { return &fromCls{} }
func wrapWhere() *whereCls     { return &whereCls{} }
func wrapGroupBy() *groupByCls { return &groupByCls{} }
func wrapHaving() *havingCls   { return &havingCls{} }
func wrapLimit() *limitCls     { return &limitCls{} }
func wrapOffset() *offsetCls   { return &offsetCls{} }

func wrapJoin() *joinCls {
	wrapper := &joinCls{}
	next := &joinNext{
		ctx:  wrapper,
		wrap: &wrapper.joinWrap,
	}
	wrapper.joinNext = next
	return wrapper
}

func wrapOrderBy() *orderByCls {
	wrapper := &orderByCls{}
	next := &orderByNext{
		ctx:  wrapper,
		wrap: &wrapper.orderByWrap,
	}
	wrapper.orderByNext = next
	return wrapper
}

func lazyInit[T any](target **T) *T {
	if *target == nil {
		*target = new(T)
	}
	return *target
}

// ------ Methods ------

// Distinct writes DISTINCT in the SELECT clause
//
// @SQL: SELECT DISTINCT ...
func (cls *selectCls) Distinct() *selectCls {
	wrp := lazyInit(&cls.selectWrap)
	wrp.Distinct()
	return cls
}

// Select writes the SELECT clause values
//
// @SQL: SELECT `DISTINCT?` `val1`, `val2`, `val3` ...
func (cls *selectCls) Select(values ...Value) {
	wrp := lazyInit(&cls.selectWrap)
	wrp.Select(vMany(values))
}

// From writes the FROM clause
//
// # This method is auto-generated for the statement’s default table.
// Only use it for complex FROM clauses (such as subqueries)
//
// @SQL: FROM `table`
func (cls *fromCls) From(table string) {
	wrp := lazyInit(&cls.fromWrap)
	wrp.From(table)
}

// Join writes the JOIN clause
//
// @SQL: INNER JOIN `table 'alias'` ... [JoinNext]
func (cls *joinCls) Join(table string) *joinNext {
	wrp := lazyInit(&cls.joinWrap)
	wrp.Join(clause.InnerJoin, table)
	return cls.joinNext
}

// On writes the join condition
//
// @SQL: [join] ... ON `on` ... [Condition]
func (cls *joinNext) On(on Value) *Condition {
	ptr := *cls.wrap
	cond := ptr.On(vOne(on))
	return wrapCondition(cond)
}

// Where writes the WHERE clause
//
// # Can be called multiple times,
// Conditions are combined using the logical "AND".
// e.g: (cond1) AND (cond2) AND (cond3) ...
//
// @SQL: WHERE `identifier` ... [Condition]
func (cls *whereCls) Where(identifier Value) *Condition {
	wrp := lazyInit(&cls.whereWrap)
	cond := wrp.Where(vOne(identifier))
	return wrapCondition(cond)
}

// GroupBy writes the GROUP BY clause
//
// @SQL: GROUP BY `val1`, `val2`, `val3` ...
func (cls *groupByCls) GroupBy(values ...Value) {
	wrp := lazyInit(&cls.groupByWrap)
	wrp.GroupBy(vMany(values))
}

// Having writes the HAVING clause
//
// @SQL: HAVING `on` ... [Condition]
func (cls *havingCls) Having(on Value) *Condition {
	wrp := lazyInit(&cls.havingWrap)
	cond := wrp.Having(vOne(on))
	return wrapCondition(cond)
}

// OrderBy writes the ORDER BY clause
//
// # Can be called multiple times for multiple orderings
//
// @SQL: ORDER BY `value` ... [OrderByNext]
func (cls *orderByCls) OrderBy(value Value) *orderByNext {
	wrp := lazyInit(&cls.orderByWrap)
	wrp.OrderBy(vOne(value))
	return cls.orderByNext
}

// Desc sets the descending direction for the current order
//
// @SQL: [OrderBy] ... DESC
func (cls *orderByNext) Desc() {
	ptr := *cls.wrap
	ptr.Desc()
}

// Limit writes the LIMIT clause value
//
// # Values less than 1 will be ignored
//
// @SQL: LIMIT `value`
func (cls *limitCls) Limit(value int) {
	if value > 0 {
		wrp := lazyInit(&cls.limitWrap)
		wrp.Limit(value)
	}
}

// Offset writes the OFFSET clause value
//
// # Values less than 1 will be ignored
//
// @SQL: OFFSET `value`
func (cls *offsetCls) Offset(value int) {
	if value > 0 {
		wrp := lazyInit(&cls.offsetWrap)
		wrp.Offset(value)
	}
}
