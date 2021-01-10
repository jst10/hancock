package custom_errors

import "fmt"

type CustomError struct {
	Code          int    `json:"code"`
	Info          string `json:"info"`
	Details       string `json:"details"`
	StackTrace    string `json:"-"`
	OriginalError error  `json:"-"`
}

func (err *CustomError) Error() string {
	return fmt.Sprintf("Custom error (%d): %s %s", err.Code, err.Info, err.Details)
}

func (err *CustomError) AST(trace string) *CustomError {
	if len(err.StackTrace) > 0 {
		err.StackTrace = trace + " -> " + err.StackTrace
	} else {
		err.StackTrace = trace
	}
	return err
}
