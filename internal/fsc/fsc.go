package fsc

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	apperror "github.com/othmaneBakkass/cv_gen/internal/common/appError"
)

// Error titles reused across this package.
const (
	titleBadOutput = "Invalid output value"
	titleBadInput  = "Invalid input value"
)

// EnsureDir validates and ensures a directory exists at the specified path.
//
// This function performs the following operations:
// 1. Normalizes the input path using filepath.Clean
// 2. Checks if the path exists and is accessible
// 3. Creates the directory (including parent directories) if it doesn't exist
// 4. Validates that the final path points to a directory, not a file
//
// Parameters:
//   - path: The directory path to validate/create (string)
//
// Returns:
//   - string: The cleaned, absolute path to the directory if successful
//   - error: An apperror.Error if the operation fails
//
// Error conditions:
//   - Permission denied when accessing or creating the directory
//   - Unable to create the directory due to system constraints
//   - Path exists but points to a file instead of a directory
//   - Directory creation succeeded but subsequent access failed
//
// Usage example:
//
//	cleanPath, err := EnsureDir("/path/to/output")
//	if err != nil {
//	    // handle error
//	}
//	// use cleanPath for further operations
func EnsureDir(path string) (string, error) {
	// Normalize the path
	cleanedPath := filepath.Clean(path)

	info, err := os.Stat(cleanedPath)
	if err != nil {
		switch {
		case os.IsNotExist(err):
			// Directory doesn't exist, try creating it
			if mkdirErr := os.MkdirAll(cleanedPath, 0755); mkdirErr != nil {
				return "", apperror.New(
					titleBadOutput,
					fmt.Sprintf("Failed to create directory: %v", mkdirErr),
					apperror.ErrorCodeArgs,
					apperror.ErrorSensitivityPublic,
				)
			}

			// Refresh FileInfo after creation
			info, err = os.Stat(cleanedPath)
			if err != nil {
				return "", apperror.New(
					titleBadOutput,
					fmt.Sprintf("Directory was created but could not be accessed: %v", err),
					apperror.ErrorCodeArgs,
					apperror.ErrorSensitivityPublic,
				)
			}

		case os.IsPermission(err):
			return "", apperror.New(
				titleBadOutput,
				"Permission denied while accessing the output path",
				apperror.ErrorCodeArgs,
				apperror.ErrorSensitivityPublic,
			)

		default:
			return "", apperror.New(
				titleBadOutput,
				fmt.Sprintf("Unable to access output path: %v", err),
				apperror.ErrorCodeArgs,
				apperror.ErrorSensitivityPublic,
			)
		}
	}

	// Last element must be a directory
	if !info.IsDir() {
		return "", apperror.New(titleBadOutput, "Output location must be a directory", apperror.ErrorCodeArgs, apperror.ErrorSensitivityPublic)
	}

	return cleanedPath, nil
}

// EnsureJSONFile validates that a JSON file exists and is readable.
//
// This function performs the following checks:
// 1. Normalizes the input path using filepath.Clean
// 2. Validates the file exists and is accessible
// 3. Confirms the path points to a file (not a directory)
// 4. Verifies the file has a .json extension
//
// It does not create anything and it does not validate the file contents.
//
// Parameters:
//   - path: Path to the JSON file to validate (string)
//
// Returns:
//   - string: The cleaned path to the JSON file if valid
//   - error: An apperror.AppError if validation fails
//
// Error conditions:
//   - File does not exist at the specified path
//   - Permission denied when accessing the file
//   - Path points to a directory instead of a file
//   - File does not have a .json extension
//
// Usage example:
//
//	jsonPath, err := EnsureJSONFile("/path/to/config.json")
//	if err != nil {
//	    // handle validation error
//	}
//	// jsonPath is now safe to use for JSON operations
func EnsureJSONFile(path string) (string, error) {
	cleanedPath := filepath.Clean(path)

	// check if file exists
	fileInfo, err := os.Stat(cleanedPath)

	if err != nil {
		if os.IsNotExist(err) {
			return "", apperror.New(titleBadInput, "Input file does not exist", apperror.ErrorCodeArgs, apperror.ErrorSensitivityPublic)
		}

		if os.IsPermission(err) {
			return "", apperror.New(
				titleBadInput,
				"Permission denied while accessing the input path",
				apperror.ErrorCodeArgs,
				apperror.ErrorSensitivityPublic,
			)
		}
		return "", apperror.New(
			titleBadInput,
			fmt.Sprintf("Unable to access input path: %v", err),
			apperror.ErrorCodeArgs,
			apperror.ErrorSensitivityPublic,
		)
	}

	// check if it's a file
	if fileInfo.IsDir() {
		return "", apperror.New(titleBadInput, "Input location must be a file", apperror.ErrorCodeArgs, apperror.ErrorSensitivityPublic)
	}

	// check if it's a json file
	if filepath.Ext(cleanedPath) != ".json" {
		return "", apperror.New(titleBadInput, "Input file must be a json file", apperror.ErrorCodeArgs, apperror.ErrorSensitivityPublic)
	}

	return cleanedPath, nil
}

// EnsureFileName creates a properly formatted filename with extension.
//
// This function handles filename normalization and fallback logic:
// 1. Normalizes the file extension (adds leading dot if missing)
// 2. Strips any directory components from name (it may be user-controlled)
// 3. Uses a timestamped default name if name is empty, "." or ".."
// 4. Strips any existing extension from the provided name
// 5. Appends the specified extension to create the final filename
//
// Parameters:
//   - name: The desired filename. Directory components are removed.
//   - defaultName: Fallback name used when name resolves to nothing usable.
//   - ext: File extension to append. Leading dot is optional.
//
// Returns:
//   - string: A bare filename (no path separators) with the given extension.
//
// Usage examples:
//
//	filename := EnsureFileName("report", "output", "pdf")     // "report.pdf"
//	filename := EnsureFileName("", "backup", ".zip")          // "backup_20240115_143022.zip"
//	filename := EnsureFileName("data.old", "export", "json")  // "data.json"
//	filename := EnsureFileName("../../etc/passwd", "cv", pdf) // "passwd.pdf"
func EnsureFileName(name, defaultName, ext string) string {
	// Normalize extension
	if ext != "" && ext[0] != '.' {
		ext = "." + ext
	}

	// Strip any directory components: the name may be user-controlled and
	// must not be able to escape the output directory.
	name = filepath.Base(strings.TrimSpace(name))
	if name == "" || name == "." || name == ".." {
		return TimestampFileName(defaultName, ext)
	}

	// Strip existing extension if present
	if filepath.Ext(name) != "" {
		name = strings.TrimSuffix(name, filepath.Ext(name))
	}

	return name + ext
}

// TimestampFileName generates a filename with a timestamp suffix.
//
// This function creates a unique filename by appending a timestamp
// in YYYYMMDD_HHMMSS format to the provided prefix. This is useful
// for creating unique file names that won't conflict with existing files.
//
// Parameters:
//   - prefix: The base name for the file (string)
//   - ext: File extension to append. Should include the leading dot (string)
//
// Returns:
//   - string: A filename in the format "prefix_YYYYMMDD_HHMMSS.ext"
//
// Timestamp format:
//   - Uses Go's reference time: "20060102_150405"
//   - Year: 4 digits (2006)
//   - Month: 2 digits (01-12)
//   - Day: 2 digits (01-31)
//   - Hour: 2 digits, 24-hour format (00-23)
//   - Minute: 2 digits (00-59)
//   - Second: 2 digits (00-59)
//
// Usage examples:
//
//	filename := TimestampFileName("backup", ".zip")    // returns "backup_20240115_143022.zip"
//	filename := TimestampFileName("log", ".txt")       // returns "log_20240115_143022.txt"
//	filename := TimestampFileName("export", ".json")   // returns "export_20240115_143022.json"
func TimestampFileName(prefix, ext string) string {
	ts := time.Now().Format("20060102_150405") // YYYYMMDD_HHMMSS
	return fmt.Sprintf("%s_%s%s", prefix, ts, ext)
}
