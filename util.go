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

var escapeOrder = []string{"\\", " ", "|", "\a", "\b", "\f", "\n", "\r", "\t", "\v"}

// Escapes a string as specified in the TeamSpeak 3 ServerQuery manual
func EscapeTS3String(str string) (returnString string) {
	returnString = str

	for i, l := 0, len(escapeOrder); i < l; i++ {
		oldString := escapeOrder[i]
		if newString := escapeStrings[oldString]; len(newString) > 0 {
			returnString = strings.Replace(returnString, oldString, newString, -1)
		}
	}

	return
}

// Does the opposite of EscapeTS3String
func UnescapeTS3String(str string) (returnString string) {
	returnString = str

	for i := len(escapeOrder) - 1; i >= 0; i-- {
		oldString := escapeOrder[i]
		if newString := escapeStrings[oldString]; len(newString) > 0 {
			returnString = strings.Replace(returnString, newString, oldString, -1)
		}
	}

	return
}
