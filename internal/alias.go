package internal

var abc = []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}

type Alias struct {
	num  int
	list map[string]byte
}

func (a *Alias) Get(table string) byte {
	if a.list == nil {
		a.list = make(map[string]byte)
	}
	if al, exists := a.list[table]; exists {
		return al
	}
	val := abc[a.num]
	a.list[table] = val
	a.num++
	return val
}
