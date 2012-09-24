package ts3

// A command structure repesents a command sent to the server
type Command struct {
	Name       string            // name of this command
	Parameters map[string]string // the key-value parameters of this command
	Options    []string          // the options (-name) of this command
}

func (cmd *Command) Encode() (encoded string) {
	encoded += cmd.Name

	for key, value := range cmd.Parameters {
		encoded += " " + EscapeTS3String(key) + "=" + EscapeTS3String(value)
	}

	for _, opt := range cmd.Options {
		encoded += " -" + EscapeTS3String(opt)
	}

	return
}
