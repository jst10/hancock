package structs

type HttpError struct {
	Details string `json:"details"`
}

func NewHttpError(details string) *HttpError {
	he := HttpError{Details: details}
	return &he
}
