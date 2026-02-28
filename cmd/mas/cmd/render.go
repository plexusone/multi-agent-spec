package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	multiagentspec "github.com/plexusone/multi-agent-spec/sdk/go"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/spf13/cobra"
)

const defaultSchemaURL = "https://raw.githubusercontent.com/plexusone/multi-agent-spec/main/schema/report/team-report.schema.json"

var (
	format       string
	boxOut       string
	narrativeOut string
	validate     bool
	schemaURL    string
)

func init() {
	rootCmd.AddCommand(renderCmd)

	renderCmd.Flags().StringVar(&format, "format", "box", "Output format for stdout: box or narrative")
	renderCmd.Flags().StringVar(&boxOut, "box-out", "", "Write box format to file")
	renderCmd.Flags().StringVar(&narrativeOut, "narrative-out", "", "Write narrative format to file")
	renderCmd.Flags().BoolVar(&validate, "validate", false, "Validate JSON against schema before rendering")
	renderCmd.Flags().StringVar(&schemaURL, "schema", "", "JSON Schema URL or file path for validation")
}

var renderCmd = &cobra.Command{
	Use:   "render [file.json]",
	Short: "Render TeamReport JSON to box or narrative format",
	Long: `Render a TeamReport JSON file to box format (terminal) or narrative
format (Pandoc-friendly Markdown).

If no file is provided, reads from stdin.

Examples:
  # Box format to stdout (default)
  mas render report.json

  # Narrative format to stdout
  mas render --format=narrative report.json

  # Both formats to separate files
  mas render --box-out=report.txt --narrative-out=report.md report.json

  # Validate before rendering
  mas render --validate report.json

  # Read from stdin
  cat report.json | mas render --format=narrative`,
	Args: cobra.MaximumNArgs(1),
	RunE: runRender,
}

func runRender(cmd *cobra.Command, args []string) error {
	// Read input
	var data []byte
	var err error

	if len(args) > 0 {
		data, err = os.ReadFile(args[0])
		if err != nil {
			return fmt.Errorf("reading file: %w", err)
		}
	} else {
		data, err = io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("reading stdin: %w", err)
		}
	}

	if len(data) == 0 {
		return fmt.Errorf("empty input")
	}

	// Validate if requested
	if validate {
		if err := validateJSON(data); err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}
	}

	// Parse report
	report, err := multiagentspec.ParseTeamReport(data)
	if err != nil {
		return fmt.Errorf("parsing report: %w", err)
	}

	// Determine what to render
	renderBox := boxOut != "" || (format == "box" && narrativeOut == "")
	renderNarrative := narrativeOut != "" || format == "narrative"

	// Render box format
	if renderBox {
		var w io.Writer = os.Stdout
		if boxOut != "" {
			f, err := os.Create(boxOut)
			if err != nil {
				return fmt.Errorf("creating box output file: %w", err)
			}
			defer f.Close()
			w = f
		}

		renderer := multiagentspec.NewRenderer(w)
		if err := renderer.Render(report); err != nil {
			return fmt.Errorf("rendering box format: %w", err)
		}
	}

	// Render narrative format
	if renderNarrative {
		var w io.Writer = os.Stdout
		if narrativeOut != "" {
			f, err := os.Create(narrativeOut)
			if err != nil {
				return fmt.Errorf("creating narrative output file: %w", err)
			}
			defer f.Close()
			w = f
		}

		renderer := multiagentspec.NewNarrativeRenderer(w)
		if err := renderer.Render(report); err != nil {
			return fmt.Errorf("rendering narrative format: %w", err)
		}
	}

	return nil
}

func validateJSON(data []byte) error {
	// Determine schema URL
	url := defaultSchemaURL
	if schemaURL != "" {
		url = schemaURL
	}

	// Compile schema
	compiler := jsonschema.NewCompiler()
	schema, err := compiler.Compile(url)
	if err != nil {
		return fmt.Errorf("compiling schema: %w", err)
	}

	// Unmarshal JSON for validation
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return fmt.Errorf("parsing JSON: %w", err)
	}

	// Validate
	if err := schema.Validate(v); err != nil {
		return err
	}

	return nil
}
