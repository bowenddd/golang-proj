package dialect

import "reflect"

var dialetsMap = map[string]Dialect{}

type Dialect interface {
	DataTypeOf(typ reflect.Value) string
	TableExistSQL(tableName string) (string, []interface{})
}

func RegisterDialect(name string, dialect Dialect) {
	dialetsMap[name] = dialect
}

func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialetsMap[name]
	return
}
