package commandhelper

func ValidateInt(key string) ValidationFunc {
	return ValidationFunc(func(value string) error {
		for _, c := range value {
			if c < '0' || c > '9' {
				return NewInvalidValue(key, "value must be an integer")
			}
		}
		return nil
	})
}

func ValidateSelection(options ...string) func(key string) ValidationFunc {
	return func(key string) ValidationFunc {
		return ValidationFunc(func(value string) error {
			optionsJoined := options[0]
			for i, option := range options {
				if option == value {
					return nil
				}

				if i > 0 {
					optionsJoined += ", " + option
				}
			}
			return NewInvalidValue(key, "value must be one of: "+optionsJoined)
		})
	}
}
