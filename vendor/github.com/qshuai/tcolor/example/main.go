package main

import (
	"fmt"

	"github.com/qshuai/tcolor"
)

func main() {
	// using the convenient package function to format a single string
	fmt.Println(tcolor.WithColor(tcolor.Green, "Hello World"))
	fmt.Println(tcolor.WithColor(tcolor.GreenBold, "Hello World"))
	fmt.Println(tcolor.WithColor(tcolor.GreenUnderline, "Hello World"))
	fmt.Println(tcolor.WithColor(tcolor.GreenBackground, "Hello World"))

	// a blank line
	fmt.Println()

	// to format a passage text
	fmt.Println(tcolor.GetColor(tcolor.Purple) + "I have a dream!")
	fmt.Println("I am happy to join with you today in what will" +
		" go down in history as the greatest demonstration for " +
		"freedom in the history of our nation.")
	fmt.Println("Five score years ago, a great American, in whose" +
		" symbolic shadow we stand today, signed the Emancipation Proclamation.")
	fmt.Println(tcolor.End)

	// the following lines will not be affected
	fmt.Println("A new line with origin color")
}
