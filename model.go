package mongogen

import (
	"fmt"
	"reflect"
	"strings"
)

var findableTypes = map[string]bool{
	"string":  true,
	"int":     true,
	"int8":    true,
	"int16":   true,
	"int32":   true,
	"int64":   true,
	"uint":    true,
	"uint8":   true,
	"uint16":  true,
	"uint32":  true,
	"uint64":  true,
	"float32": true,
	"float64": true,
	"struct":  true,
}

type Model struct {
	Name       string
	Collection string
	Fields     []*Field
}

func NewModel(n string) *Model {
	return &Model{
		Name:   n,
		Fields: make([]*Field, 0),
	}
}

func (m *Model) String() string {
	fields := make([]string, 0)
	for _, f := range m.Fields {
		fields = append(fields, "\t"+f.String()+"\n")
	}

	fieldsStr := strings.Join(fields, "")
	str := fmt.Sprintf("(Model '%s' [\n %s]", m.Name, fieldsStr)

	return str
}

func (m *Model) FindableFields() []*Field {
	fields := make([]*Field, 0)
	for _, f := range m.Fields {
		if f.Findable() {
			fields = append(fields, f)
		}
	}

	return fields
}

type Field struct {
	Name   string
	Type   string
	Tag    reflect.StructTag
	Fields []*Field
}

func NewField(n string, t string) *Field {
	return &Field{
		Name:   n,
		Type:   t,
		Fields: make([]*Field, 0),
	}
}

func (f *Field) String() string {
	fields := make([]string, 0)
	for _, f := range f.Fields {
		fields = append(fields, f.String())
	}

	fieldsStr := strings.Join(fields, ", ")

	return fmt.Sprintf("%s %s %s [%s]", f.Name, f.Type, f.Tag, fieldsStr)
}
func (f *Field) GetTagValue(key string) string {
	if f.Tag == "" {
		return ""
	}

	return f.Tag.Get(key)
}

func (f *Field) DbName() string {
	name := f.GetTagValue("bson")
	endFieldName := strings.Index(name, ",")
	if endFieldName != -1 {
		name = name[:endFieldName]
	}

	if name == "" {
		name = strings.ToLower(f.Name)
	}

	return name
}

func (f *Field) FindableType() string {
	return strings.Replace(f.Type, "[]", "", 1)
}

func (f *Field) Findable() bool {
	return findableTypes[f.FindableType()]
}
