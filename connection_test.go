package ts3

import (
	"net"
	"testing"
)

func TestConnection(t *testing.T) {
	addr, err := net.ResolveTCPAddr("tcp", "ts3.nightexcessive.us:10011")

	if err != nil {
		t.Errorf("failed to resolve TCP address: %s", err)
		return
	}

	connection, err := Dial(addr)
	defer connection.Close()

	if err != nil {
		t.Errorf("failed to connect to TeamSpeak server: %s", err)
		return
	}

	res, err := connection.SendCommand(&Command{
		Name: "use",
		Parameters: map[string]string{
			"port": "9987",
		},
	})

	if err != nil {
		t.Errorf("failed to send use command to TeamSpeak server: %s", err)
		return
	}

	t.Logf("received use results: %#v", res)

	command := &Command{
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
