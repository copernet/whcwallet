package tcolor

const End = "\033[0m"

var isPanic bool

// SetPanic set the package variable isPanic = true. That means
// the program using this library will panic if the specified
// color type not found.
func SetPanic() {
	isPanic = true
}

// WithColor pack the input string with the specified decorated
// color type and return the assembled string.
func WithColor(colorType color, str string) string {
	if int(colorType) > len(colorSet)-1 {
		if isPanic {
			panic("not found the specified color type")
		}

		return str
	}

	return colorSet[int(colorType)] + str + End
}

// GetColor return the origin defined color label for your custom
// format using. For example, you can using the returned color label
// to get a passage text instead of the single string.
func GetColor(colorType color) string {
	if int(colorType) > len(colorSet)-1 {
		if isPanic {
			panic("not found the specified color type")
		}

		return ""
	}

	return colorSet[int(colorType)]
}
