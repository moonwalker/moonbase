package cms

import (
	jsonschema "github.com/santhosh-tekuri/jsonschema/v5"
)

type Schema struct {
	jsonSchema *jsonschema.Schema
}

func NewSchema(schema []byte) *Schema {
	jsonSchema, err := jsonschema.CompileString("", string(schema))
	if err != nil {
		return &Schema{}
	}
	return &Schema{jsonSchema}
}

func (s *Schema) Validate(v any) error {
	if s.jsonSchema != nil {
		return s.jsonSchema.Validate(v)
	}
	return nil
}
