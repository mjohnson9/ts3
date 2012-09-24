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
			t.Errorf("Escape(%s) = %s, want %s", raw, escaped, expected)
		}
	}
}

func TestUnescape(t *testing.T) {
	for expected, raw := range escapeTests {
		if escaped := UnescapeTS3String(raw); escaped != expected {
			t.Errorf("Unescape(%s) = %s, want %s", raw, escaped, expected)
		}
	}
}

func BenchmarkEscape(b *testing.B) {
	b.StopTimer()

	oldString := "TeamSpeak ]|[ Server"

	b.SetBytes(int64(len(oldString)))

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		EscapeTS3String(oldString)
	}
}

func BenchmarkUnescape(b *testing.B) {
	b.StopTimer()

	oldString := "TeamSpeak\\s]\\p[\\sServer"

	b.SetBytes(int64(len(oldString)))

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		UnescapeTS3String(oldString)
	}
}
