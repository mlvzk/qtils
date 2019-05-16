package commandhelper

func ValidateInt(key string) ValidationFunc {
	return ValidationFunc(func(value string) error {
		for _, c := range value {
			if c < '0' || c > '9' {
				return NewValueParsingError(key, "value must be an integer")
			}
		}
		return nil
	})
}
