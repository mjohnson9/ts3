package ts3

import (
	"net"
	"testing"
)

func TestConnection(t *testing.T) {
	addr, err := net.ResolveTCPAddr("tcp", "public-ts3.cz:10011")
	t.Logf("connection address: %s", addr)

	if err != nil {
		t.Errorf("failed to resolve TCP address: %s", err)
		return
	}

	connection, err := Dial(addr)
	if err != nil {
		t.Errorf("failed to connect to TeamSpeak server: %s", err)
		return
	}
	defer connection.Close()

	command := &Command{
		Name: "use",
		Parameters: map[string]string{
			"port": "9987",
		},
	}
	res, err := connection.SendCommand(command)

	if err != nil {
		t.Errorf("failed to send use command to TeamSpeak server: %s", err)
		return
	}

	t.Logf("received use results (%s): %#v", command, res)

	command = &Command{
		Name: "clientlist",
		Options: []string{
			"away",
			"times",
		},
	}

	res, err = connection.SendCommand(command)

	if err != nil {
		t.Errorf("failed to send clientlist command to TeamSpeak server: %s", err)
		return
	}

	t.Logf("received client list (%s): %#v", command, res)
}
