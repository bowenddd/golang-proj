package schema

import (
	"geeorm/dialet"
	"go/ast"
	"reflect"
)

type Field struct {
	Type string
	Name string
	Tag  string
}

type Schema struct {
	Model      interface{}
	Name       string
	Fields     []*Field
	FieldNames []string
	fieldMap   map[string]*Field
}

func (s *Schema) GetField(name string) *Field {
	return s.fieldMap[name]
}

func Parse(dest interface{}, d dialet.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schame := &Schema{
		Model:    dest,
		Name:     modelType.Name(),
		fieldMap: map[string]*Field{},
	}

	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			if v, ok := p.Tag.Lookup("geeorm"); ok {
				field.Tag = v
			}
			schame.FieldNames = append(schame.FieldNames, p.Name)
			schame.Fields = append(schame.Fields, field)
			schame.fieldMap[p.Name] = field
		}
	}
	return schame
}
