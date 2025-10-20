package clause

type GroupByBuilder struct {
	targets []*GroupBy
}

func (gb *GroupByBuilder) New(value string) *GroupBy {
	group := &GroupBy{target: value}
	gb.targets = append(gb.targets, group)
	return group
}

type GroupBy struct {
	target string
	having string
}

func (gb *GroupBy) Having(value string) {
	gb.having = value
}
