package generate

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/othmaneBakkass/cv_gen/cmd/root"
	apperror "github.com/othmaneBakkass/cv_gen/internal/common/appError"
	"github.com/othmaneBakkass/cv_gen/internal/fsc"
	tone "github.com/othmaneBakkass/cv_gen/internal/pdf/templates/t1"
	"github.com/spf13/cobra"
)

var command = &cobra.Command{
	Use:   "generate",
	Short: "Generate a cv based on json.",
	Long:  "Turn json data to a cv in a supported format (pdf).",
	RunE:  handler,
}

func GenerateCommandExecute() {
	command.Execute()
}

func init() {
	command.Flags().StringP("output", "o", ".", "Path where the generated file will be located. defaults to the current directory.")
	command.Flags().StringP("input", "i", "", "Path where the data file is located.")
	root.RootCommand.AddCommand(command)
}

func handler(cmd *cobra.Command, args []string) error {
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		return apperror.New("Invalid output value",
			"Output value must be a valid string",
			apperror.ErrorCodeArgs,
			apperror.ErrorSensitivityPublic)
	}

	outputDir, err := fsc.EnsureDir(output)
	if err != nil {
		return err
	}

	input, err := cmd.Flags().GetString("input")
	if err != nil {
		return apperror.New("Invalid input value", "Input value must be a valid string",
			apperror.ErrorCodeArgs,
			apperror.ErrorSensitivityPublic)
	}
	if input == "" {
		return apperror.New("Invalid input value", "Input path not specified, please provide the path of the json file",
			apperror.ErrorCodeArgs,
			apperror.ErrorSensitivityPublic)
	}

	jsonInput, err := fsc.EnsureJSONFile(input)
	if err != nil {
		return err
	}

	content, err := os.ReadFile(jsonInput)
	if err != nil {
		return apperror.New("Failed to read input file",
			err.Error(),
			apperror.ErrorCodeArgs,
			apperror.ErrorSensitivityPublic)
	}

	var inputData tone.InputData
	err = json.Unmarshal(content, &inputData)
	if err != nil {
		return apperror.New("Failed to parse JSON",
			err.Error(),
			apperror.ErrorCodeArgs,
			apperror.ErrorSensitivityPublic)
	}

	var inputDataSize = len(inputData.Data)
	if inputDataSize == 0 {
		return apperror.New("No data provided",
			"The data array is empty",
			apperror.ErrorCodeArgs,
			apperror.ErrorSensitivityPublic)
	}

	for _, v := range inputData.Data {
		if _, err := tone.Validate(&v); err != nil {
			return err
		}
		if v.Template != "t1" {
			return apperror.New("Unsupported template",
				"Only t1 template is supported",
				apperror.ErrorCodeArgs,
				apperror.ErrorSensitivityPublic)
		}

		fileName := fsc.EnsureFileName(v.FileName, "", "pdf")
		path := filepath.Join(outputDir, fileName)
		tone.GenerateT1PDF(path, v)
	}

	return nil
}
