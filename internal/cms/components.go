package cms

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

// ->  /\.modules?\.css$/i;

func compsPlugin(tree map[string]string, config compsConfig) api.Plugin {
	return api.Plugin{
		Name: "comps",
		Setup: func(build api.PluginBuild) {
			build.OnResolve(api.OnResolveOptions{Filter: `.*`},
				func(args api.OnResolveArgs) (api.OnResolveResult, error) {

					// TODO: improve
					s := strings.Split(args.Path, "/")[0]
					for _, d := range config.Dependencies {
						if d == s {
							return api.OnResolveResult{
								Path:     args.Path,
								External: true,
							}, nil
						}
					}

					path := args.Path

					if len(filepath.Ext(path)) == 0 {
						path += "/index.js"
					}

					if args.Kind == api.ResolveJSImportStatement {
						dir := filepath.Dir(args.Importer)
						path = filepath.Join(dir, path)
					}

					return api.OnResolveResult{
						Path:      path,
						Namespace: "comps-ns",
					}, nil
				})
			build.OnLoad(api.OnLoadOptions{Filter: `.*`, Namespace: "comps-ns"},
				func(args api.OnLoadArgs) (api.OnLoadResult, error) {
					contents := tree[args.Path]

					if filepath.Ext(args.Path) == ".css" {
						contents = "/* css module */"
					}

					return api.OnLoadResult{
						Contents: &contents,
						Loader:   api.LoaderDefault,
					}, nil
				})
		},
	}
}

func BundleComponents(tree map[string]string, config compsConfig, preserveJSX bool, minify bool) (string, error) {

	// TODO: improve
	entry := config.EntryDir() + "/index.js"

	opts := api.BuildOptions{
		EntryPoints: []string{entry},
		// External:    config.Dependencies,
		Loader: map[string]api.Loader{
			".js": api.LoaderJSX,
		},
		Plugins: []api.Plugin{compsPlugin(tree, config)},
		Format:  api.FormatESModule,
		Bundle:  true,
	}

	if preserveJSX {
		opts.JSXMode = api.JSXModePreserve
	}

	if minify {
		opts.MinifyWhitespace = true
		opts.MinifyIdentifiers = true
		opts.MinifySyntax = true
	}

	result := api.Build(opts)

	if len(result.Errors) > 0 {
		esb := new(strings.Builder)
		for _, e := range result.Errors {
			fmt.Fprintf(esb, "%s; ", e.Text)
		}
		return "", errors.New(strings.TrimSpace(esb.String()))
	}

	rsb := new(strings.Builder)
	for _, out := range result.OutputFiles {
		fmt.Fprintf(rsb, "%s", out.Contents)
	}

	return rsb.String(), nil
}
