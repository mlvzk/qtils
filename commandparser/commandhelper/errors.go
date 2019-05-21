package commandhelper

type _error struct {
	msg string
}

func NewError(msg string) *_error {
	return &_error{msg}
}

func (e *_error) Error() string {
	return e.msg
}

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

type InvalidValue struct {
	key    string
	suffix string
}

func NewInvalidValue(key string, information string) *InvalidValue {
	return &InvalidValue{
		key:    key,
		suffix: information,
	}
}

func (e *InvalidValue) Error() string {
	return "invalid value of key '" + e.key + "'; " + e.suffix
}
