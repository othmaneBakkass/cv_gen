package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/othmaneBakkass/cv_gen/cmd/cv_gen/root"
	_ "github.com/othmaneBakkass/cv_gen/cmd/cv_gen/shell"
	apperror "github.com/othmaneBakkass/cv_gen/internal/common/appError"
)

func main() {
	err := root.RootCommandExecute()
	if err != nil {
		var appErr apperror.AppError
		if errors.As(err, &appErr) {
			if appErr.Sensitivity == apperror.ErrorSensitivityPublic {
				fmt.Fprintln(os.Stderr, "Title:", appErr.Title, "Details:", appErr.Detail)
			}
			os.Exit(1)
		} else {
			fmt.Fprintln(os.Stderr, "Unexpected error:", err)
			os.Exit(2)
		}
	}
}
