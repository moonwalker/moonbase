package cache

import (
	"testing"
	"time"
)

type testy struct {
	Foo string
}

func TestCacheGeneric(t *testing.T) {
	genc := NewGeneric[*testy](time.Minute * 1)
	ty := &testy{Foo: "bar"}

	err := genc.Set("foo", ty)
	if err != nil {
		t.Error(err)
	}

	res, err := genc.Get("foo")
	if err != nil {
		t.Error(err)
	}

	if res.Foo != "bar" {
		t.Fail()
	}
}
