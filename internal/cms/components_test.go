package cms

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var tree = map[string]string{
	"/src/components/index.js": "import { Foo } from './foo.js'\nexport const Greet = () => <h1>Hello, world from {Foo}!</h1>",
	"/src/components/foo.js":   "export const Foo = { bar: 'baz' }",
}

func TestBundleComponents(t *testing.T) {
	data, _ := os.ReadFile(yamlPath)
	config := ParseConfig(data)
	res, err := BundleComponents(tree, config.Components, false, false)
	if err != nil {
		t.Error(err)
	}
	println(res)
}

func TestBundleComponentsFromJSON(t *testing.T) {
	td, _ := os.ReadFile("testdata/compstree.json")
	var tree map[string]string
	err := json.Unmarshal(td, &tree)
	if err != nil {
		t.Error(err)
	}

	data, _ := os.ReadFile(yamlPath)
	config := ParseConfig(data)
	res, err := BundleComponents(tree, config.Components, false, false)
	if err != nil {
		t.Error(err)
	}
	println(res)
}

func getDir(s string) string {
	if len(filepath.Ext(s)) > 0 {
		return filepath.Dir(s)
	}
	return s
}

func TestPathMethods(t *testing.T) {
	var1 := "src/components"
	var2 := "src/components/index.js"

	if getDir(var1) != var1 {
		t.Fail()
	}

	if getDir(var2) != var1 {
		t.Fail()
	}

	fmt.Println(strings.Split("next/link", "/")[0])
	fmt.Println(strings.Split("next", "/")[0])
}
