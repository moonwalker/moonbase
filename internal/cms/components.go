package cms

import (
	"os"

	"github.com/evanw/esbuild/pkg/api"
)

func BundleComponents(config *Config) string {
	result := api.Build(api.BuildOptions{
		//EntryNames: ,
		EntryPoints: []string{"src/components"},
		Loader: map[string]api.Loader{
			".js": api.LoaderJSX,
		},
		External:          []string{"next"},
		JSXMode:           api.JSXModePreserve,
		Format:            api.FormatESModule,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		Bundle:            true,
		Write:             true,
		Outdir:            "out",
	})

	if len(result.Errors) > 0 {
		os.Exit(1)
	}

	return ""
}
