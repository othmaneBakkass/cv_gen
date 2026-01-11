package shell

import (
	"github.com/othmaneBakkass/cv_gen/cmd/cv_gen/root"
	apperror "github.com/othmaneBakkass/cv_gen/internal/common/appError"

	"github.com/spf13/cobra"
)

var ShellCommand = &cobra.Command{
	Use:   "shell",
	Short: "Start interactive mode",
	Long:  "Start interactive mode which allows you to run multiple commands at once in dev mode",
	RunE: func(cmd *cobra.Command, args []string) error {
		return apperror.New(
			"Error from shell command",
			"Detailed explanation",
			apperror.ErrorCodeUnknown,
			apperror.ErrorSensitivityPublic)
	},
}

func ShellCommandExecute() {
	ShellCommand.Execute()
}

func init() {
	root.RootCommand.AddCommand(ShellCommand)
}
