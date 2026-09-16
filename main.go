package main

import (
	"errors"
	"fmt"
	"os"

	_ "github.com/othmaneBakkass/cv_gen/cmd/generate"
	"github.com/othmaneBakkass/cv_gen/cmd/root"
	_ "github.com/othmaneBakkass/cv_gen/cmd/serve"
	apperror "github.com/othmaneBakkass/cv_gen/internal/common/appError"
	"github.com/othmaneBakkass/cv_gen/internal/common/logs"
)

func main() {
	err := root.RootCommandExecute()
	if err == nil {
		return
	}

	var appErr apperror.AppError
	if !errors.As(err, &appErr) {
		fmt.Fprintln(os.Stderr, logs.ErrorLog(err.Error()))
		os.Exit(2)
	}

	if appErr.Sensitivity == apperror.ErrorSensitivityPublic {
		fmt.Fprintln(os.Stderr, logs.ErrorLog(fmt.Sprintf("%s: %s", appErr.Title, appErr.Detail)))
		for _, issue := range appErr.Issues {
			if issue.Sensitivity == apperror.ErrorSensitivityPublic {
				fmt.Fprintln(os.Stderr, logs.ErrorLog(fmt.Sprintf("  - %s: %s", issue.Title, issue.Detail)))
			}
		}
	} else {
		fmt.Fprintln(os.Stderr, logs.ErrorLog("An unexpected error occurred."))
	}
	os.Exit(1)
}
