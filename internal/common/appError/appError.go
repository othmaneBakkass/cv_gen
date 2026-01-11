package apperror

// ErrorSensitivity represents how safe it is to expose an error to clients
type ErrorSensitivity string

const (
	ErrorSensitivityPublic   ErrorSensitivity = "public"   // can be logged and viewed by users
	ErrorSensitivityPrivate  ErrorSensitivity = "private"  // can be logged but not viewed by users
	ErrorSensitivitySanitize ErrorSensitivity = "sanitize" // needs Sanitation to be public
)

// AppErrorIssue represents a specific issue within an error
type AppErrorIssue struct {
	Title       string           // Optional short title
	Detail      string           // Required detailed message
	Sensitivity ErrorSensitivity // How sensitive the issue is
}

type ErrorCode string

const (
	ErrorCodeUnknown ErrorCode = "unknown_error"
)

type AppError struct {
	Title       string
	Detail      string
	Code        ErrorCode
	Sensitivity ErrorSensitivity
	Issues      []AppErrorIssue
}

// implement Error interface for AppError
func (e AppError) Error() string {
	return e.Title + ": " + e.Detail
}

// constructor for AppError struct
func New(title, detail string, code ErrorCode, sensitivity ErrorSensitivity, issues ...[]AppErrorIssue) AppError {
	var errs []AppErrorIssue
	if len(issues) > 0 {
		errs = issues[0]
	}
	return AppError{
		Title:       title,
		Detail:      detail,
		Code:        code,
		Sensitivity: sensitivity,
		Issues:      errs,
	}
}
