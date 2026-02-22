package root

import "github.com/spf13/cobra"

var RootCommand = &cobra.Command{
	Use:           "cv_gen",
	Short:         "cv_gen is a tool for generating CVs.",
	Long:          "cv_gen is a tool for generating CVs based on a json file.",
	SilenceErrors: true,
}

func RootCommandExecute() error {
	return RootCommand.Execute()
}
