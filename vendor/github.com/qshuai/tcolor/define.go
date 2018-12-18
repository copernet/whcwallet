package tcolor

type color uint8

const (
	Black color = iota
	Red
	Green
	Yellow
	Blue
	Purple
	Cyan
	White

	BlackBold
	RedBold
	GreenBold
	YellowBold
	BlueBold
	PurpleBold
	CyanBold
	WhiteBold

	BlackUnderline
	RedUnderline
	GreenUnderline
	YellowUnderline
	BlueUnderline
	PurpleUnderline
	CyanUnderline
	WhiteUnderline

	BlackBackground
	RedBackground
	GreenBackground
	YellowBackground
	BlueBackground
	PurpleBackground
	CyanBackground
	WhiteBackground
)

var colorSet = [...]string{
	// regular style
	"\033[0;30m",
	"\033[0;31m",
	"\033[0;32m",
	"\033[0;33m",
	"\033[0;34m",
	"\033[0;35m",
	"\033[0;36m",
	"\033[0;37m",
	// bold style
	"\033[1;30m",
	"\033[1;31m",
	"\033[1;32m",
	"\033[1;33m",
	"\033[1;34m",
	"\033[1;35m",
	"\033[1;36m",
	"\033[1;37m",
	// with underline
	"\033[4;30m",
	"\033[4;31m",
	"\033[4;32m",
	"\033[4;33m",
	"\033[4;34m",
	"\033[4;35m",
	"\033[4;36m",
	"\033[4;37m",
	// with background
	"\033[40m",
	"\033[41m",
	"\033[42m",
	"\033[43m",
	"\033[44m",
	"\033[45m",
	"\033[46m",
	"\033[47m",
}
