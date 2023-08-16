package cms

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type Schema struct {
	jsonSchema *CollectionSchema
}

type CollectionSchema struct {
	Name   string        `json:"name"`
	Fields []SchemaField `json:"fields,omitempty"`
}

type SchemaField struct {
	Name     string            `json:"name"`
	Label    string            `json:"label"`
	Type     string            `json:"type"`
	List     bool              `json:"list"`
	Required bool              `json:"required"`
	Object   *CollectionSchema `json:"object,omitempty"`
}

var schemaTypeMap = map[string]string{
	"string":  "string",
	"float64": "number",
	"bool":    "boolean",
}

func NewSchema(schema []byte) *Schema {
	jsonSchema := &CollectionSchema{}
	err := json.Unmarshal(schema, jsonSchema)
	if err != nil {
		return &Schema{}
	}
	return &Schema{jsonSchema}
}

func (s *Schema) Available() bool {
	return s.jsonSchema != nil
}

func (s *Schema) Validate(v map[string]interface{}) error {
	if s.Available() {
		return validateJson(s.jsonSchema, v, "")
	}
	return nil
}

func (s *Schema) ValidateString(data string) error {
	var v map[string]interface{}
	err := json.Unmarshal([]byte(data), &v)
	if err != nil {
		return err
	}
	return s.Validate(v)
}

func GenerateSchema(name string, contents string) (string, error) {
	var v map[string]interface{}
	err := json.Unmarshal([]byte(contents), &v)
	if err != nil {
		return "", err
	}

	fields := parseJson("", v)
	cs := CollectionSchema{
		Name:   name,
		Fields: fields,
	}
	schemaStr, err := json.MarshalIndent(cs, "", "  ")
	if err != nil {
		return "", err
	}

	return string(schemaStr), nil
}

func parseJson(parentKey string, obj interface{}) []SchemaField {
	fields := make([]SchemaField, 0)
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

			subFieldMap, ok := firstElement.(map[string]interface{})
			if !ok {
				fields = append(fields, SchemaField{
					Name:     parentKey + key,
					Label:    parentKey + key,
					Type:     reflect.TypeOf(firstElement).Name(),
					List:     true,
					Required: true,
				})
			} else {
				subFields := parseJson("", subFieldMap)
				fields = append(fields, SchemaField{
					Name:     parentKey + key,
					Label:    parentKey + key,
					Type:     "object",
					List:     true,
					Required: true,
					Object: &CollectionSchema{
						Name:   parentKey + key,
						Fields: subFields,
					},
				})
			}
		case map[string]interface{}:
			fields = append(fields, parseJson(parentKey+key, v)...)
		default:
			fields = append(fields, SchemaField{
				Name:     parentKey + key,
				Label:    parentKey + key,
				Type:     reflect.TypeOf(v).Name(),
				Required: true,
			})
		}
	}

	return fields
}

func validateJson(schema *CollectionSchema, v map[string]interface{}, parent string) error {
	if len(parent) > 0 {
		parent += "/"
	}
	for _, f := range schema.Fields {
		if f.Required {
			fieldName := fmt.Sprintf("%s%s", parent, f.Name)
			if v[f.Name] == nil {
				return fmt.Errorf("missing field: %s", fieldName)
			}
			if f.Type == "object" {
				var objects []interface{}
				var ok bool
				if f.List {
					objects, ok = v[f.Name].([]interface{})
					if !ok {
						return fmt.Errorf("invalid input list format at field: %s", fieldName)
					}
				} else {
					objects = []interface{}{v[f.Name]}
				}
				for _, object := range objects {
					o, ok := object.(map[string]interface{})
					if !ok {
						return fmt.Errorf("invalid input format at field: %s", fieldName)
					}
					err := validateJson(f.Object, o, fieldName)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}
