package cms

import (
	"encoding/json"
	"reflect"

	jsonschemaValidate "github.com/santhosh-tekuri/jsonschema/v5"
)

const JsonSchemaName = "_schema.json"

type Schema struct {
	jsonSchema *jsonschemaValidate.Schema
}

type schemaProperties map[string]interface{}
type schemaObject struct {
	Schema     string           `json:"$schema,omitempty"`
	Type       string           `json:"type"`
	Properties schemaProperties `json:"properties,omitempty"`
	Required   []string         `json:"required,omitempty"`
}
type schemaArray struct {
	Type  string       `json:"type"`
	Items schemaObject `json:"items"`
}

const jsonSchemaTag = "http://json-schema.org/draft-07/schema#"

var schemaTypeMap = map[string]string{
	"string":  "string",
	"float64": "number",
	"bool":    "boolean",
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
	var v map[string]interface{}
	err := json.Unmarshal([]byte(contents), &v)
	if err != nil {
		return "", err
	}

	fields := make(map[string]interface{})
	fieldNames := parseJson("", v, fields)

	so := schemaObject{
		Schema:     jsonSchemaTag,
		Type:       "object",
		Properties: fields,
		Required:   fieldNames,
	}
	schemaStr, err := json.MarshalIndent(so, "", "  ")
	if err != nil {
		return "", err
	}

	return string(schemaStr), nil
}

func parseJson(parentKey string, obj interface{}, fields map[string]interface{}) []string {
	fieldNames := make([]string, 0)
	if len(parentKey) > 0 {
		parentKey = parentKey + "."
	}
	for key, v := range obj.(map[string]interface{}) {
		switch v.(type) {
		case []interface{}:
			// Array of some sort infer type from 1st elem
			arrayConversion, ok := v.([]interface{})
			if !ok {
				break
			}
			firstElement := arrayConversion[0]
			fieldNames = append(fieldNames, parentKey+key)

			subFieldMap, ok := firstElement.(map[string]interface{})
			if !ok {
				fields[parentKey+key] = schemaArray{
					Type: "array",
					Items: schemaObject{
						Type: reflect.TypeOf(firstElement).Name(),
					},
				}
			} else {
				subFields := make(map[string]interface{})
				subFieldNames := parseJson("", subFieldMap, subFields)
				fields[parentKey+key] = schemaArray{
					Type: "array",
					Items: schemaObject{
						Type:       "object",
						Properties: subFields,
						Required:   subFieldNames,
					},
				}
			}
			break
		case map[string]interface{}:
			fieldNames = append(fieldNames, parseJson(parentKey+key, v, fields)...)
			break
		default:
			// No more nested so append the type
			fieldNames = append(fieldNames, parentKey+key)
			fields[parentKey+key] = map[string]string{
				"type": schemaTypeMap[reflect.TypeOf(v).Name()],
			}
			break
		}
	}

	return fieldNames
}
