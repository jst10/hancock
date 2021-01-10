package custom_errors

func GetNotFoundError() *CustomError {
	return &CustomError{Code: 1, Info: "Not found"}
}
func GetMissingQueryParamError(details string) *CustomError {
	return &CustomError{Code: 2, Info: "Missing query param", Details: details}
}
func GetCookieNotPresentError(details string) *CustomError {
	return &CustomError{Code: 3, Info: "Missing cookie", Details: details}
}
func GetErrorDecodingPostBodyError(details string) *CustomError {
	return &CustomError{Code: 4, Info: "Error at decoding post body", Details: details}
}
func GetNotSupportedMediaTypeInRequest(details string) *CustomError {
	return &CustomError{Code: 5, Info: "Not supported media type in request", Details: details}
}
func GetJsonDecodingError(originalError error, details string) *CustomError {
	return &CustomError{Code: 6, OriginalError: originalError, Info: "Error decoding request body into json", Details: details}
}
func GetDbOpenError(originalError error) *CustomError {
	return &CustomError{Code: 7, OriginalError: originalError, Info: "Db open error"}
}
func GetDbPingError(originalError error) *CustomError {
	return &CustomError{Code: 8, OriginalError: originalError, Info: "Db ping error"}
}
func GetDbExecError(originalError error) *CustomError {
	return &CustomError{Code: 9, OriginalError: originalError, Info: "Db exec error"}
}
func GetDbQueryError(originalError error) *CustomError {
	return &CustomError{Code: 10, OriginalError: originalError, Info: "Db query error"}
}
func GetDbScanError(originalError error) *CustomError {
	return &CustomError{Code: 11, OriginalError: originalError, Info: "Db scan error"}
}
func GetDbGetLastInsertedIdError(originalError error) *CustomError {
	return &CustomError{Code: 12, OriginalError: originalError, Info: "Db get last inserted id error"}
}
func GetNotValidDataError(details string) *CustomError {
	return &CustomError{Code: 13, Info: "Not valid data", Details: details}
}
func GetNotAllowed(details string) *CustomError {
	return &CustomError{Code: 14, Info: "Not allowed", Details: details}
}
func GetErrorCreateSalt(originalError error) *CustomError {
	return &CustomError{Code: 15, OriginalError: originalError, Info: "Error create salt"}
}
func GetErrorCreateJWT(originalError error, details string) *CustomError {
	return &CustomError{Code: 16, OriginalError: originalError, Info: "Error create  JWT", Details: details}
}
func GetErrorVerifyJWT(originalError error, details string) *CustomError {
	return &CustomError{Code: 17, OriginalError: originalError, Info: "Error verify JWT", Details: details}
}
func GetErrorGettingCookie(originalError error, details string) *CustomError {
	return &CustomError{Code: 18, OriginalError: originalError, Info: "Error getting cookie", Details: details}
}
func GetInvalidJWT(details string) *CustomError {
	return &CustomError{Code: 19, Info: "Invalid JWT", Details: details}
}
func GetErrorLoadingConfigs(originalError error, details string) *CustomError {
	return &CustomError{Code: 20, OriginalError: originalError, Info: "Error loading config", Details: details}
}

func GetNotFoundErrorr(originalError error) *CustomError {
	return &CustomError{Code: 20, OriginalError: originalError, Info: "Not found"}
}
