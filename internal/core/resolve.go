package core

func ResolveColumnWriter(opt *ColumnWriter) {
	resolve := func(provided *ColumnType, dflt ColumnType) {
		if *provided != ColTypUnset {
			return
		}
		*provided = dflt
	}

	dflt := &DefaultColumnWriter

	resolve(&opt.Select, dflt.Select)
	resolve(&opt.Join, dflt.Join)
	resolve(&opt.Where, dflt.Where)
	resolve(&opt.GroupBy, dflt.GroupBy)
	resolve(&opt.Having, dflt.Having)
	resolve(&opt.OrderBy, dflt.OrderBy)
}
