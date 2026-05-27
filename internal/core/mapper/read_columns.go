package mapper

import (
	"reflect"
	"slices"

	"github.com/laacin/inyorm/internal/core/mapper/types"
)

func (m *Mapper) ReadCols(entries ...any) []string {
	c := collector(map[string]struct{}{})

	for _, entry := range entries {
		info := types.ReadInfo(reflect.TypeOf(entry))

		if info.Kind != types.KindStruct && !info.CanBeDeref() {
			continue
		}

		switch info.Kind {
		case types.KindStruct:
			colsByStruct(info, c)

		case types.KindMap:
			colsByMap(entry, info, c)

		case types.KindString:
			colsByString(entry, info, c)

		case types.KindColumn:
			colsByCol(entry, info, c)
		}
	}

	return c.ToSlice()
}

// --- Readers
type collector map[string]struct{}

func (c collector) Add(col string) {
	if col != "" {
		c[col] = struct{}{}
	}
}

func (c collector) ToSlice() []string {
	out := make([]string, 0, len(c))

	for col := range c {
		out = append(out, col)
	}

	slices.Sort(out)
	return out
}

func colsByStruct(info types.TypeInfo, c collector) {
	info.Schema.IterNames(func(s string) {
		c.Add(s)
	})
}

func colsByMap(v any, info types.TypeInfo, c collector) {
	if info.IsSlc() {
		if info.IsSlcOfPtrs() {
			slc := types.UnwrapSlc[*map[string]any](v, info.IsPtr())

			for _, m := range slc {
				if m != nil {
					for k := range *m {
						c.Add(k)
					}
				}
			}

			return
		}

		slc := types.UnwrapSlc[map[string]any](v, info.IsPtr())
		for _, m := range slc {
			for k := range m {
				c.Add(k)
			}
		}

		return
	}

	m := types.Unwrap[map[string]any](v, info.IsPtr())
	for k := range m {
		c.Add(k)
	}
}

func colsByCol(v any, info types.TypeInfo, c collector) {
	if info.IsSlc() {
		slc := types.UnwrapSlc[interface{ BaseName() string }](v, info.IsPtr())
		for _, col := range slc {
			c.Add(col.BaseName())
		}

		return
	}
	c.Add(types.Unwrap[interface{ BaseName() string }](v, false).BaseName())
}

func colsByString(v any, info types.TypeInfo, c collector) {
	if info.IsSlc() {
		if info.IsSlcOfPtrs() {
			slc := types.UnwrapSlc[*string](v, info.IsPtr())

			for _, s := range slc {
				if s != nil {
					c.Add(*s)
				}
			}
			return
		}
		for _, s := range types.UnwrapSlc[string](v, info.IsPtr()) {
			c.Add(s)
		}
		return
	}
	c.Add(types.Unwrap[string](v, info.IsPtr()))
}
