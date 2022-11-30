package cms

import (
	"encoding/json"
)

type PackageJSON struct {
	Dependencies map[string]string `json:"dependencies"`
}

func ParsePackageJSON(data []byte) *PackageJSON {
	pkg := &PackageJSON{}
	json.Unmarshal(data, pkg)
	return pkg
}
