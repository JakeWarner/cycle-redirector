package log

import (
	"fmt"
)

var colors = map[string]int{
	// purple
	"darkpurple": 53, "purple": 57, "magenta": 129, "lightpurple": 63,

	// blue
	"darkblue": 18, "blue": 21, "skyblue": 26, "teal": 36, "lightblue": 39,

	// green
	"darkgreen": 22, "green": 28, "kelly": 29, "olive": 70, "lightgreen": 40, "lime": 46,

	// yellow
	"darkyellow": 220, "yellow": 226, "lightyellow": 228,

	// orange
	"darkorange": 166, "orange": 202, "lightorange": 208,

	// red
	"darkred": 88, "red": 160, "pink": 207, "lightred": 196,

	// black/grayscale
	"black": 232, "darkgray": 233, "mediumgray": 237, "gray": 243, "lightgray": 249, "white": 256,
}

func TextColor(contents, color string) string {
	if code, ok := colors[color]; ok {
		contents = fmt.Sprintf("\033[38;5;%dm%v\033[0m", code, contents)
	}

	return contents
}

func BgColor(contents, color string) string {
	if code, ok := colors[color]; ok {
		contents = fmt.Sprintf("\033[48;5;%dm%v\033[0m", code, contents)
	}
	return contents
}

func Bold(contents string) string {
	return fmt.Sprintf("\033[1m%s\033[0m", contents)
}

func Underline(contents string) string {
	return fmt.Sprintf("\033[4m%s\033[0m", contents)
}
