package cms

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestValidateActualSchema(t *testing.T) {
	sch, err := ioutil.ReadFile("testdata/schema.json")
	if err != nil {
		t.Error(err)
	}

	data, err := ioutil.ReadFile("testdata/payload.json")
	if err != nil {
		t.Error(err)
	}

	var v interface{}
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
	data, err := ioutil.ReadFile("testdata/payload.json")
	if err != nil {
		t.Error(err)
	}

	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		t.Error(err)
	}

	s := NewSchema([]byte{})
	err = s.Validate(v)
	if err != nil {
		t.Error(err)
	}
}
