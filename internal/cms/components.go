package cms

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

func compsPlugin(tree map[string]string) api.Plugin {
	return api.Plugin{
		Name: "comps",
		Setup: func(build api.PluginBuild) {
			build.OnResolve(api.OnResolveOptions{Filter: `.*`},
				func(args api.OnResolveArgs) (api.OnResolveResult, error) {
					path := args.Path

					if args.Kind == api.ResolveEntryPoint {
						// path = "/" + args.Path
					}

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
						contents = "/* css */"
						fmt.Println(contents)
					}

					if len(contents) == 0 {
						contents = "/* external */"
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

	// TODO: this needs to be improved
	entry := config.EntryDir() + "/index.js"

	opts := api.BuildOptions{
		EntryPoints: []string{entry},
		External:    config.Dependencies,
		Loader: map[string]api.Loader{
			".js": api.LoaderJSX,
		},
		Plugins: []api.Plugin{compsPlugin(tree)},
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
