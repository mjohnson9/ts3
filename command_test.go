package ts3

import (
	"testing"
)

func TestEncoding(t *testing.T) {
	var (
		command  *Command
		expected string
	)

	command = &Command{
		Name: "use",
		Parameters: map[string]string{
			"port": "9987",
		},
	}

	expected = "use port=9987"

	if encoded := command.Encode(); encoded != expected {
		t.Errorf("Escape(%#v) = %v, want %v", command, encoded, expected)
	}

	command = &Command{
		Name: "use",
		Parameters: map[string]string{
			"port": "9987",
		},
	}

	expected = "use port=9987"

	if encoded := command.Encode(); encoded != expected {
		t.Errorf("Escape(%#v) = %v, want %v", command, encoded, expected)
	}
}
