package cmd

import (
	"fmt"
	"reflect"

	"github.com/invopop/jsonschema"
	"github.com/spf13/cobra"

	"github.com/twpayne/chezmoi/v2/pkg/chezmoi"
)

func (c *Config) newJSONSchemaCmd() *cobra.Command {
	jsonSchemaCmd := &cobra.Command{
		Use:       "json-schema component",
		Short:     "Generate JSON schemas",
		Args:      cobra.ExactArgs(1),
		ValidArgs: []string{"config", "external"},
		// Long: mustLongHelp("json-schema"),
		Example: example("completion"),
		RunE:    c.runJSONSchemaCmd,
		Annotations: map[string]string{
			doesNotRequireValidConfig: "true",
		},
	}

	return jsonSchemaCmd
}

func (c *Config) runJSONSchemaCmd(cmd *cobra.Command, args []string) error {
	reflector := &jsonschema.Reflector{
		BaseSchemaID: "https://chezmoi.io/schemas",
		Mapper: func(t reflect.Type) *jsonschema.Schema {
			fmt.Printf("pkgPath=%s name=%s\n", t.PkgPath(), t.Name())
			switch t.PkgPath() {
			case "github.com/twpayne/chezmoi/v2/pkg/chezmoi":
				switch t.Name() {
				case "AbsPath":
					return &jsonschema.Schema{
						Type: "string",
					}
				case "EntryTypeSet":
					return &jsonschema.Schema{
						Type: "string",
					}
				case "Mode":
					return &jsonschema.Schema{
						Type: "string",
						Enum: []interface{}{
							"symlink",
						},
					}
				}
			case "github.com/twpayne/chezmoi/v2/pkg/cmd":
				switch t.Name() {
				case "autoBool":
					return &jsonschema.Schema{
						AnyOf: []*jsonschema.Schema{
							{
								Type: "bool",
							},
							{
								Type:  "string",
								Const: "auto",
							},
						},
					}
				}
			}
			return nil
		},
	}
	var schema *jsonschema.Schema
	switch args[0] {
	case "config":
		schema = reflector.Reflect(&Config{})
	case "external":
		schema = reflector.Reflect(&chezmoi.External{})
	default:
		return fmt.Errorf("%s: unknown component", args[0])
	}

	data, err := chezmoi.FormatJSON.Marshal(schema)
	if err != nil {
		return err
	}
	return c.writeOutput(data)
}
