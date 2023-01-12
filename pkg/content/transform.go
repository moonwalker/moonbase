package content

import (
	"github.com/moonwalker/gontentful"
)

func TransformModel(model *gontentful.ContentType) (*Schema, error) {
	panic("not implemented")
}

func FormatSchema(schema *Schema) (*gontentful.ContentType, error) {
	panic("not implemented")
}

func TransformEntry(model *gontentful.Entry) (*ContentData, error) {
	panic("not implemented")
}

func FormatData(data *ContentData) (*gontentful.Entry, error) {
	panic("not implemented")
}
