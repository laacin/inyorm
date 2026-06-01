package ddl

type PrimaryKey struct {
	Cols []string
}

func NewPrimaryKey(cols []string) *PrimaryKey {
	return &PrimaryKey{Cols: cols}
}

// --- Build

func (*PrimaryKey) Build() error {
	return nil
}
