package commandparser

type _error struct {
	msg string
}

func NewError(msg string) *_error {
	return &_error{msg}
}

func (e *_error) Error() string {
	return e.msg
}

var EmptyArgvError = NewError("argument list is empty")

type InvalidKeyError struct {
	key string
}

func NewInvalidKeyError(key string) *InvalidKeyError {
	return &InvalidKeyError{
		key: key,
	}
}

func (e InvalidKeyError) Error() string {
	return "invalid key '" + e.key + "'"
}

type MissingValueError struct {
	key string
}

func NewMissingValueError(key string) *MissingValueError {
	return &MissingValueError{
		key: key,
	}
}

func (e *MissingValueError) Error() string {
	return "missing value for key '" + e.key + "'"
}
