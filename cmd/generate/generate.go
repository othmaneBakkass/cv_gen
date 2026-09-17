package generate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/othmaneBakkass/cv_gen/cmd/root"
	apperror "github.com/othmaneBakkass/cv_gen/internal/common/appError"
	"github.com/othmaneBakkass/cv_gen/internal/common/logs"
	"github.com/othmaneBakkass/cv_gen/internal/fsc"
	"github.com/othmaneBakkass/cv_gen/internal/pdf/i18n"
	"github.com/othmaneBakkass/cv_gen/internal/pdf/render"
	"github.com/othmaneBakkass/cv_gen/internal/pdf/templates"
	"github.com/othmaneBakkass/cv_gen/internal/schema"
	"github.com/othmaneBakkass/cv_gen/internal/settings"
)

// sectionToggleFlags map a "--no-X" boolean flag to the section it force-
// disables. Sections required by the schema (education, experience,
// languages) have no toggle — they're always present in valid input.
var sectionToggleFlags = []struct {
	flag string
	key  render.SectionKey
}{
	{"no-profile", render.Profile},
	{"no-skills", render.Skills},
	{"no-projects", render.Projects},
	{"no-certifications", render.Certifications},
}

// sectionPriorityFlags map a "--<section>" int flag to the section whose
// render-order priority it overrides. Lower renders first.
var sectionPriorityFlags = []struct {
	flag string
	key  render.SectionKey
}{
	{"profile", render.Profile},
	{"education", render.Education},
	{"experience", render.Experience},
	{"skills", render.Skills},
	{"projects", render.Projects},
	{"certifications", render.Certifications},
	{"languages", render.Languages},
}

const titleBadInput = "Invalid input value"

var command = &cobra.Command{
	Use:   "generate",
	Short: "Generate a CV from JSON.",
	Long:  "Turn JSON data into a CV in a supported format (PDF).",
	RunE:  handler,
}

func init() {
	command.Flags().StringP("output", "o", ".", "Directory where generated files are written (default: current directory).")
	command.Flags().StringP("input", "i", "", "Path to the JSON data file (required).")
	command.Flags().Bool("french", false, "Render section labels in French instead of English.")
	command.Flags().String("density", "normal", "Spacing density: dense, normal, airy, or adaptive (auto-fit to one page).")
	command.Flags().String("template", "", fmt.Sprintf("Override the template for every entry (%s). Default: use each entry's own \"template\" field.", strings.Join(templates.Names(), ", ")))
	for _, t := range sectionToggleFlags {
		command.Flags().Bool(t.flag, false, fmt.Sprintf("Omit the %s section even if present in the data.", t.key))
	}
	for _, p := range sectionPriorityFlags {
		command.Flags().Int(p.flag, 0, fmt.Sprintf("Render priority for the %s section (lower = earlier); default order is preserved unless set.", p.key))
	}
	root.RootCommand.AddCommand(command)
}

// optionsForCV builds this CV entry's render.Options: the JSON's own
// "settings" block as the base, then any CLI flag the user actually typed
// (per cmd.Flags().Changed) applied on top as an override. A flag left at
// its default never touches a setting the JSON already specified — only
// an explicitly-passed flag wins.
func optionsForCV(cmd *cobra.Command, s *schema.Settings) (render.Options, error) {
	opts, err := optionsFromSettings(s)
	if err != nil {
		return opts, err
	}
	return applyCLIOverrides(cmd, opts)
}

// optionsFromSettings converts a CV entry's JSON "settings" block into its
// base render.Options. It is a thin wrapper over settings.FromSettings so
// the CLI and the HTTP server (cmd/serve) share one conversion.
func optionsFromSettings(s *schema.Settings) (render.Options, error) {
	return settings.FromSettings(s)
}

// applyCLIOverrides applies only the flags the user actually passed
// (cmd.Flags().Changed) on top of opts, so an unset flag never clobbers a
// setting the JSON already specified.
func applyCLIOverrides(cmd *cobra.Command, opts render.Options) (render.Options, error) {
	if cmd.Flags().Changed("french") {
		french, err := cmd.Flags().GetBool("french")
		if err != nil {
			return opts, apperror.New(titleBadInput, "french must be a valid boolean", apperror.ErrorCodeArgs, apperror.ErrorSensitivityPublic)
		}
		if french {
			opts.Lang = i18n.FR
		} else {
			opts.Lang = i18n.EN
		}
	}

	if cmd.Flags().Changed("density") {
		density, err := cmd.Flags().GetString("density")
		if err != nil {
			return opts, apperror.New(titleBadInput, "density must be a valid string", apperror.ErrorCodeArgs, apperror.ErrorSensitivityPublic)
		}
		mult, ok := settings.DensityMultipliers[density]
		if !ok {
			return opts, apperror.New(titleBadInput, "density must be one of: dense, normal, airy, adaptive", apperror.ErrorCodeArgs, apperror.ErrorSensitivityPublic)
		}
		opts.Density = mult
	}

	for _, t := range sectionToggleFlags {
		if !cmd.Flags().Changed(t.flag) {
			continue
		}
		off, err := cmd.Flags().GetBool(t.flag)
		if err != nil {
			return opts, apperror.New(titleBadInput, fmt.Sprintf("%s must be a valid boolean", t.flag), apperror.ErrorCodeArgs, apperror.ErrorSensitivityPublic)
		}
		opts.Include[t.key] = !off
	}

	for _, p := range sectionPriorityFlags {
		if !cmd.Flags().Changed(p.flag) {
			continue
		}
		v, err := cmd.Flags().GetInt(p.flag)
		if err != nil {
			return opts, apperror.New(titleBadInput, fmt.Sprintf("%s must be a valid integer", p.flag), apperror.ErrorCodeArgs, apperror.ErrorSensitivityPublic)
		}
		opts.Priority[p.key] = v
	}

	return opts, nil
}

func handler(cmd *cobra.Command, _ []string) error {
	fmt.Println(logs.InfoLog("PDF generation started"))

	output, err := cmd.Flags().GetString("output")
	if err != nil {
		return apperror.New("Invalid output value", "Output value must be a valid string",
			apperror.ErrorCodeArgs, apperror.ErrorSensitivityPublic)
	}

	outputDir, err := fsc.EnsureDir(output)
	if err != nil {
		return err
	}

	input, err := cmd.Flags().GetString("input")
	if err != nil {
		return apperror.New(titleBadInput, "Input value must be a valid string",
			apperror.ErrorCodeArgs, apperror.ErrorSensitivityPublic)
	}
	if input == "" {
		return apperror.New(titleBadInput, "Input path not specified, please provide the path of the JSON file with -i",
			apperror.ErrorCodeArgs, apperror.ErrorSensitivityPublic)
	}

	jsonPath, err := fsc.EnsureJSONFile(input)
	if err != nil {
		return err
	}

	content, err := os.ReadFile(jsonPath)
	if err != nil {
		return apperror.New("Failed to read input file", err.Error(),
			apperror.ErrorCodeArgs, apperror.ErrorSensitivityPublic)
	}

	var inputData schema.InputData
	if err := json.Unmarshal(content, &inputData); err != nil {
		return apperror.New("Failed to parse JSON", err.Error(),
			apperror.ErrorCodeArgs, apperror.ErrorSensitivityPublic)
	}

	if len(inputData.Data) == 0 {
		return apperror.New("No data provided", "The data array is empty",
			apperror.ErrorCodeArgs, apperror.ErrorSensitivityPublic)
	}

	templateOverride, err := cmd.Flags().GetString("template")
	if err != nil {
		return apperror.New(titleBadInput, "template must be a valid string", apperror.ErrorCodeArgs, apperror.ErrorSensitivityPublic)
	}
	if templateOverride != "" {
		valid := false
		for _, n := range templates.Names() {
			if n == templateOverride {
				valid = true
				break
			}
		}
		if !valid {
			return apperror.New(titleBadInput,
				fmt.Sprintf("template must be one of: %s", strings.Join(templates.Names(), ", ")),
				apperror.ErrorCodeArgs, apperror.ErrorSensitivityPublic)
		}
	}

	fmt.Println(logs.InfoLog("Processing JSON data"))

	for i := range inputData.Data {
		cv := inputData.Data[i]
		if templateOverride != "" {
			cv.Template = templateOverride
		}

		if err := schema.Validate(&cv); err != nil {
			return err
		}

		opts, err := optionsForCV(cmd, cv.Settings)
		if err != nil {
			return err
		}

		fileName := fsc.EnsureFileName(cv.FileName, "cv", "pdf")
		outPath := filepath.Join(outputDir, fileName)

		if err := templates.Generate(outPath, cv, opts); err != nil {
			return err
		}
		fmt.Println(logs.InfoLog("Generated: " + fileName))
	}

	fmt.Println(logs.SuccessLog("Generation finished"))
	return nil
}
