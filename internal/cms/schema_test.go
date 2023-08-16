package cms

import (
	"encoding/json"
	"os"
	"testing"
)

func TestValidateActualSchema(t *testing.T) {
	sch, err := os.ReadFile("testdata/schema.json")
	if err != nil {
		t.Error(err)
	}

	data, err := os.ReadFile("testdata/payload.json")
	if err != nil {
		t.Error(err)
	}

	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		t.Error(err)
	}

	s := NewSchema(sch)
	err = s.Validate(v)
	if err != nil {
		t.Error(err)
	}
}

func TestValidateEmptySchema(t *testing.T) {
	data, err := os.ReadFile("testdata/payload.json")
	if err != nil {
		t.Error(err)
	}

	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		t.Error(err)
	}

	s := NewSchema([]byte{})
	if s.Available() {
		t.Fail()
	}
	err = s.Validate(v)
	if err != nil {
		t.Error(err)
	}
}

func TestGenerateSchema(t *testing.T) {
	data, err := os.ReadFile("testdata/payload.json")
	if err != nil {
		t.Error(err)
	}

	schema, err := GenerateSchema("test", string(data))
	if err != nil {
		t.Error(err)
	}
	t.Log(schema)
}
