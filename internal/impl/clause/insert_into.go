package clause

import "github.com/laacin/inyorm/internal/entity"

type InsertIntoImpl[Next any] struct {
	declared bool
	emb      entity.InsertInto
}
