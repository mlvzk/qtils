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

type ValueParsingError struct {
	key    string
	suffix string
}

func NewValueParsingError(key string, suffix string) *ValueParsingError {
	return &ValueParsingError{
		key:    key,
		suffix: suffix,
	}
}

func (e *ValueParsingError) Error() string {
	return "missing required argument '" + e.key + "'; " + e.suffix
}
