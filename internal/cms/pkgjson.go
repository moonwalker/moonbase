package cms

import (
	"encoding/json"
)

const PackageJSONFile = "package.json"

type PackageJSON struct {
	Dependencies map[string]string `json:"dependencies"`
}

func ParsePackageJSON(data []byte) *PackageJSON {
	pkg := &PackageJSON{}
	json.Unmarshal(data, pkg)
	return pkg
}

func ResolveDepsVersions(pkgJson *PackageJSON, depNames []string) map[string]string {
	res := make(map[string]string)
	for _, depName := range depNames {
		depVersion := pkgJson.Dependencies[depName]
		if len(depVersion) > 0 {
			res[depName] = depVersion
		}
	}
	return res
}
