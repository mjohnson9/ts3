package ts3

import (
	"testing"
)

var escapeTests = map[string]string{
	"TeamSpeak ]|[ Server": "TeamSpeak\\s]\\p[\\sServer",
	"\\":                   "\\\\",
	" ":                    "\\s",
	"|":                    "\\p",
	"\a":                   "\\a",
	"\b":                   "\\b",
	"\f":                   "\\f",
	"\n":                   "\\n",
	"\r":                   "\\r",
	"\t":                   "\\t",
	"\v":                   "\\v",
	" |\a\b\f\n\r\t\v":     "\\s\\p\\a\\b\\f\\n\\r\\t\\v",
}

func TestEscape(t *testing.T) {
	for raw, expected := range escapeTests {
		if escaped := EscapeTS3String(raw); escaped != expected {
			t.Errorf("Escape(%v) = %v, want %v", raw, escaped, expected)
		}
	}
}
