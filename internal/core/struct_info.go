package core

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
