package cms

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

// ->  /\.modules?\.css$/i;

// https://github.com/evanw/esbuild/issues/594#issuecomment-744737854

var css string

func compsPlugin(tree map[string]string, config compsConfig, bundleExternals bool) api.Plugin {
	return api.Plugin{
		Name: "comps",
		Setup: func(build api.PluginBuild) {
			build.OnResolve(api.OnResolveOptions{Filter: ".*"},
				func(args api.OnResolveArgs) (api.OnResolveResult, error) {

					// TODO: improve?
					s := strings.Split(args.Path, "/")[0]
					for _, d := range config.Dependencies {
						if d == s {
							if bundleExternals {
								return api.OnResolveResult{
									Path:      fmt.Sprintf("https://unpkg.com/%s", d),
									Namespace: "http-url",
								}, nil
							} else {
								return api.OnResolveResult{
									Path:     args.Path,
									External: true,
								}, nil
							}
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

			build.OnLoad(api.OnLoadOptions{Filter: ".*", Namespace: "comps-ns"},
				func(args api.OnLoadArgs) (api.OnLoadResult, error) {
					contents := tree[args.Path]

					// if filepath.Ext(args.Path) == ".css" {
					// 	contents = "const result = 1; export defailt result;"
					// }

					return api.OnLoadResult{
						Contents: &contents,
						Loader:   api.LoaderDefault,
					}, nil
				})

			build.OnResolve(api.OnResolveOptions{Filter: ".*", Namespace: "http-url"},
				func(args api.OnResolveArgs) (api.OnResolveResult, error) {
					base, err := url.Parse(args.Importer)
					if err != nil {
						return api.OnResolveResult{}, err
					}
					relative, err := url.Parse(args.Path)
					if err != nil {
						return api.OnResolveResult{}, err
					}
					return api.OnResolveResult{
						Path:      base.ResolveReference(relative).String(),
						Namespace: "http-url",
					}, nil
				})

			build.OnLoad(api.OnLoadOptions{Filter: ".*", Namespace: "http-url"},
				func(args api.OnLoadArgs) (api.OnLoadResult, error) {
					res, err := http.Get(args.Path)
					if err != nil {
						return api.OnLoadResult{}, err
					}
					defer res.Body.Close()
					bytes, err := ioutil.ReadAll(res.Body)
					if err != nil {
						return api.OnLoadResult{}, err
					}
					contents := string(bytes)
					return api.OnLoadResult{Contents: &contents}, nil
				})
		},
	}
}

func BundleComponents(tree map[string]string, config compsConfig, bundleExternals bool, preserveJSX bool, minify bool) (string, error) {

	// TODO: improve
	entry := config.EntryDir() + "/index.js"

	opts := api.BuildOptions{
		EntryPoints: []string{entry},
		Loader: map[string]api.Loader{
			".js":         api.LoaderJSX,
			".module.css": api.LoaderText,
		},
		Plugins: []api.Plugin{compsPlugin(tree, config, bundleExternals)},
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
