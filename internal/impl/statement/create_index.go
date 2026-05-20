package statement

// type CreateIndexStmtImpl struct {
// 	DefaultRef string
// 	Dialect    ir.Dialect
//
// 	table.ConsDeclImpl
// }
//
// func (s *CreateIndexStmtImpl) Start(ctx context.Context, eng *ir.Engine, ref string) api.CreateIndex {
// 	s.DefaultRef = ref
// 	s.Dialect = eng.Dialect
// 	return s
// }
//
// // --- Build
//
// func (s *CreateIndexStmtImpl) Build() string {
// 	w := &writer.WriterImpl{Syntax: s.Dialect}
// 	// s.TableBuilderImpl.Build(w, s.Dialect)
// 	return w.ToString()
// }
