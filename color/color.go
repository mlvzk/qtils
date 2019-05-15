package color

const escape = "\x1b"

func Green(str string) string {
	return escape + "[32m" + str + escape + "[0m"
}

func Gold(str string) string {
	return escape + "[33m" + str + escape + "[0m"
}
