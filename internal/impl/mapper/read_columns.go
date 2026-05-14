package mapper

import (
	"reflect"
	"slices"

	"github.com/laacin/inyorm/internal/impl/exprimpl"
	"github.com/laacin/inyorm/internal/impl/mapper/types"
)

func ReadColumns(entries []any) []string {
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
	for _, fi := range info.Fis {
		c.Add(fi.Meta.Name)
	}
}

func colsByMap(v any, info types.TypeInfo, c collector) {
	if info.IsSlc() {
		if info.IsSlcOfPtrs() {
			slc := unwrapSlc[*map[string]any](v, info.IsPtr())

			for _, m := range slc {
				if m != nil {
					for k := range *m {
						c.Add(k)
					}
				}
			}

			return
		}

		slc := unwrapSlc[map[string]any](v, info.IsPtr())
		for _, m := range slc {
			for k := range m {
				c.Add(k)
			}
		}

		return
	}

	m := unwrap[map[string]any](v, info.IsPtr())
	for k := range m {
		c.Add(k)
	}
}

func colsByCol(v any, info types.TypeInfo, c collector) {
	if info.IsSlc() {
		slc := unwrapSlc[*exprimpl.ColumnImpl](v, info.IsPtr())
		for _, col := range slc {
			c.Add(col.Name)
		}

		return
	}
	c.Add(unwrap[*exprimpl.ColumnImpl](v, false).Name)
}

func colsByString(v any, info types.TypeInfo, c collector) {
	if info.IsSlc() {
		if info.IsSlcOfPtrs() {
			slc := unwrapSlc[*string](v, info.IsPtr())

			for _, s := range slc {
				if s != nil {
					c.Add(*s)
				}
			}
			return
		}
		for _, s := range unwrapSlc[string](v, info.IsPtr()) {
			c.Add(s)
		}
		return
	}
	c.Add(unwrap[string](v, info.IsPtr()))
}

// --- Helpers

func unwrap[T any](v any, ptr bool) T {
	if ptr {
		return *v.(*T)
	}
	return v.(T)
}

func unwrapSlc[T any](v any, ptr bool) []T {
	if ptr {
		return *v.(*[]T)
	}
	return v.([]T)
}
