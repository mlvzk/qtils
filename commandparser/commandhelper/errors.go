package commandhelper

type MissingRequiredError struct {
	key string
}

func NewMissingRequiredError(key string) *MissingRequiredError {
	return &MissingRequiredError{
		key: key,
	}
}

func (e *MissingRequiredError) Error() string {
	return "missing required argument '" + e.key + "'"
}
