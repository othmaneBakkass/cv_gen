package apperror

// ErrorSensitivity represents how safe it is to expose an error to users.
type ErrorSensitivity string

const (
	ErrorSensitivityPublic   ErrorSensitivity = "public"   // can be logged and shown to users
	ErrorSensitivityPrivate  ErrorSensitivity = "private"  // can be logged but not shown to users
	ErrorSensitivitySanitize ErrorSensitivity = "sanitize" // needs sanitation before being shown
)

// AppErrorIssue is a single problem within an AppError (e.g. one invalid field).
type AppErrorIssue struct {
	Title       string
	Detail      string
	Sensitivity ErrorSensitivity
}

// ErrorCode is a coarse category for an AppError.
type ErrorCode string

const (
	ErrorCodeUnknown ErrorCode = "unknown_error"
	ErrorCodeArgs    ErrorCode = "args_error"
)

// AppError is an application error carrying user-facing text and,
// optionally, a list of finer-grained issues.
type AppError struct {
	Title       string
	Detail      string
	Code        ErrorCode
	Sensitivity ErrorSensitivity
	Issues      []AppErrorIssue
}

// Error implements the error interface.
func (e AppError) Error() string {
	return e.Title + ": " + e.Detail
}

// New builds an AppError. Attach field-level detail with WithIssues.
func New(title, detail string, code ErrorCode, sensitivity ErrorSensitivity) AppError {
	return AppError{
		Title:       title,
		Detail:      detail,
		Code:        code,
		Sensitivity: sensitivity,
	}
}

// WithIssues returns a copy of e with issues appended.
func (e AppError) WithIssues(issues ...AppErrorIssue) AppError {
	e.Issues = append(e.Issues, issues...)
	return e
}
