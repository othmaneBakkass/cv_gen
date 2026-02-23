package main

import (
	"errors"
	"fmt"
	"os"

	_ "github.com/othmaneBakkass/cv_gen/cmd/generate"
	"github.com/othmaneBakkass/cv_gen/cmd/root"
	apperror "github.com/othmaneBakkass/cv_gen/internal/common/appError"
	"github.com/othmaneBakkass/cv_gen/internal/common/logs"
)

func main() {
	err := root.RootCommandExecute()
	if err != nil {
		var appErr apperror.AppError
		if errors.As(err, &appErr) {
			if appErr.Sensitivity == apperror.ErrorSensitivityPublic {
				fmt.Println(logs.ErrorLog(fmt.Sprintf("Title: %s Details: %s", appErr.Title, appErr.Detail)))
			}
			os.Exit(1)
		} else {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
	}
}
