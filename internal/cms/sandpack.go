package cms

func SandpackResolveDeps(pkgJson *PackageJSON, depNames []string) map[string]string {
	res := make(map[string]string)
	for _, depName := range depNames {
		depVersion := pkgJson.Dependencies[depName]
		if len(depVersion) > 0 {
			res[depName] = depVersion
		}
	}
	return res
}
