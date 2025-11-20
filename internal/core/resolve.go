package core

func ResolveAutoPlaceholder(opt **AutoPlaceholder) {
	ptr := *opt
	if ptr == nil {
		ptr = &DefaultAutoPlaceholder
	}
	*opt = ptr
}

func ResolveColumnWriter(opt **ColumnWriter) {
	ptr := *opt
	if ptr == nil {
		ptr = &ColumnWriter{}
	}

	resolve := func(provided *ColumnType, dflt ColumnType) {
		if *provided != ColTypUnset {
			return
		}
		*provided = dflt
	}

	dflt := &DefaultColumnWriter

	resolve(&ptr.Select, dflt.Select)
	resolve(&ptr.Join, dflt.Join)
	resolve(&ptr.Where, dflt.Where)
	resolve(&ptr.GroupBy, dflt.GroupBy)
	resolve(&ptr.Having, dflt.Having)
	resolve(&ptr.OrderBy, dflt.OrderBy)

	*opt = ptr
}
