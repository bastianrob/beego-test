package services

//APIError e.g: malformed JSON, missing required properties, error unmarshall
type APIError struct {
	Code    int
	Message string
}

func (e *APIError) Error() string {
	return e.Message
}
