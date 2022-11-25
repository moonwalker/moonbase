package cms

import (
	"encoding/json"

	jsonschemaGenerate "github.com/invopop/jsonschema"
	jsonschemaValidate "github.com/santhosh-tekuri/jsonschema/v5"
)

type Schema struct {
	jsonSchema *jsonschemaValidate.Schema
}

func NewSchema(schema []byte) *Schema {
	jsonSchema, err := jsonschemaValidate.CompileString("", string(schema))
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

func GenerateSchema(contents string) (string, error) {
	var v any
	err := json.Unmarshal([]byte(contents), &v)
	if err != nil {
		return "", err
	}
	s := jsonschemaGenerate.Reflect(v)
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
