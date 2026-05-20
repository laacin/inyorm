package std_ddl

import "github.com/laacin/inyorm/internal/ir/ddl"

var mapOnAct = map[ddl.OnAction]string{
	ddl.OnActionCascade:  "CASCADE",
	ddl.OnActionSetNull:  "SET NULL",
	ddl.OnActionDefault:  "SET DEFAULT",
	ddl.OnActionRestrict: "RESTRICT",
	ddl.OnActionNoAction: "NO ACTION",
}

var mapColKind = map[ddl.ColKind]string{
	ddl.ColKindText:  "TEXT",
	ddl.ColKindInt:   "INTEGER",
	ddl.ColKindFloat: "DOUBLE",
	ddl.ColKindBool:  "BOOLEAN",
}
