package commandhelper

func ValidateInt(key string) ValidationFunc {
	err := NewInvalidValue(key, "value must be an integer")

	return ValidationFunc(func(value string) error {
		if len(value) == 0 || !(value[0] == '-' || (value[0] >= '0' && value[0] <= '9')) {
			return err
		}

		for _, c := range value[1:] {
			if c < '0' || c > '9' {
				return err
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

func ValidateKeyValue(delimiter string) func(key string) ValidationFunc {
	return func(key string) ValidationFunc {
		return ValidationFunc(func(value string) error {
			for i := 0; i < len(value)-len(delimiter)+1; i++ {
				if value[i:i+len(delimiter)] == delimiter {
					return nil
				}
			}

			return NewInvalidValue(key, "value does not contain the '"+delimiter+"' delimiter")
		})
	}
}
