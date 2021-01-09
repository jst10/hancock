package custom_errors

func GetNotFoundError() *CustomError {
	return &CustomError{Code: 1, Info: "Not found"}
}

func GetMissingQueryParamError(details string) *CustomError {
	return &CustomError{Code: 2, Info: "Missing query param", Details: details}
}

func GetNotFoundErrorr(originalError error) *CustomError {
	return &CustomError{Code: 1, OriginalError: originalError, Details: "Not found"}
}
