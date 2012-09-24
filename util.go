package ts3

import (
	"strings"
)

var escapeStrings = map[string]string{
	"\\": "\\\\",
	" ":  "\\s",
	"|":  "\\p",
	"\a": "\\a",
	"\b": "\\b",
	"\f": "\\f",
	"\n": "\\n",
	"\r": "\\r",
	"\t": "\\t",
	"\v": "\\v",
}

// Escapes a string as specified in the TeamSpeak 3 ServerQuery manual
func EscapeTS3String(str string) (returnString string) {
	returnString = str

	for oldString, newString := range escapeStrings {
		returnString = strings.Replace(returnString, oldString, newString, -1)
	}

	return
}
