package cms

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

var tree = map[string]string{
	"/src/components/index.js": "import { Foo } from './foo.js'\nexport const Greet = () => <h1>Hello, world from {Foo}!</h1>",
	"/src/components/foo.js":   "export const Foo = { bar: 'baz' }",
}

var compsPlugin = api.Plugin{
	Name: "comps",
	Setup: func(build api.PluginBuild) {
		build.OnResolve(api.OnResolveOptions{Filter: `.*`},
			func(args api.OnResolveArgs) (api.OnResolveResult, error) {
				var path string

				if args.Kind == api.ResolveEntryPoint {
					path = "/" + args.Path
				}

				if args.Kind == api.ResolveJSImportStatement {
					dir := filepath.Dir(args.Importer)
					path = filepath.Join(dir, args.Path)
				}

				return api.OnResolveResult{
					Path:      path,
					Namespace: "ns-comps",
				}, nil
			})
		build.OnLoad(api.OnLoadOptions{Filter: `.*`, Namespace: "ns-comps"},
			func(args api.OnLoadArgs) (api.OnLoadResult, error) {
				contents := tree[args.Path]
				return api.OnLoadResult{
					Contents: &contents,
					Loader:   api.LoaderDefault,
				}, nil
			})
	},
}

func BundleComponents(enyry string, external []string, preserveJSX bool, minify bool) (string, error) {
	opts := api.BuildOptions{
		EntryPoints: []string{enyry},
		External:    external,
		Loader: map[string]api.Loader{
			".js": api.LoaderJSX,
		},
		Plugins: []api.Plugin{compsPlugin},
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
			fmt.Fprintf(esb, "%s", e.Text)
		}
		return "", errors.New(esb.String())
	}

	rsb := new(strings.Builder)
	for _, out := range result.OutputFiles {
		fmt.Fprintf(rsb, "%s", out.Contents)
	}

	return rsb.String(), nil
}
