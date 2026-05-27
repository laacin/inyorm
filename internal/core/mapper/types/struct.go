package types

import "reflect"

func readStruct(t reflect.Type) StructInfo {
	info := NewStructInfo()

	for field := range t.Fields() {
		typ, _ := DerefPtrTyp(field.Type)

		if field.Anonymous && typ.Kind() == reflect.Struct {
			baseIndex := append([]int(nil), field.Index...)
			info.Merge(baseIndex, readStruct(typ))
			continue
		}

		info.Add(field.Name, field.Tag.Get(TAG), field.Index)
	}

	return info
}

type StructInfo map[string]FieldInfo

type FieldInfo struct {
	Index []int
	Tag   TagResult
}

func NewStructInfo() StructInfo { return StructInfo(map[string]FieldInfo{}) }

func (m StructInfo) Add(fieldName, tag string, idx []int) {
	tagResult := ParseTag(fieldName, tag)

	m[tagResult.Name] = FieldInfo{
		Tag:   tagResult,
		Index: idx,
	}
}

func (m StructInfo) Merge(baseIndex []int, other StructInfo) {
	for name, info := range other {
		idx := append(baseIndex, info.Index...)

		m[name] = FieldInfo{
			Index: idx,
			Tag:   info.Tag,
		}
	}
}

func (m StructInfo) Get(name string) (FieldInfo, bool) {
	r, ok := m[name]
	return r, ok
}

func (m StructInfo) GetIndex(name string) ([]int, bool) {
	r, ok := m[name]
	return r.Index, ok
}

func (m StructInfo) IterNames(fn func(string)) {
	for name := range m {
		fn(name)
	}
}
