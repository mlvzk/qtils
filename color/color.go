package color

const escape = "\x1b"

var Important = Red
var Section = Olive
var Info = Teal

func Red(str string) string {
	return escape + "[31m" + str + escape + "[0m"
}

func Green(str string) string {
	return escape + "[32m" + str + escape + "[0m"
}

func Olive(str string) string {
	return escape + "[33m" + str + escape + "[0m"
}

func Teal(str string) string {
	return escape + "[36m" + str + escape + "[0m"
}
