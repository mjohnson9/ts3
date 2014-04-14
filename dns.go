package ts3

import (
	"net"
	"strconv"
)

// This looks up a TeamSpeak host, immitating the manner in which the TeamSpeak
// client looks it up. hostCombo should either be a host:port combination or a
// lone hostname. If a port is not specified, the default port of 9987 will be
// used.
func ResolveHost(hostCombo string) ([]string, error) {
	hostname, port, err := net.SplitHostPort(hostCombo)
	if addrErr, ok := err.(*net.AddrError); ok && addrErr.Err == "missing port in address" {
		hostname = hostCombo
		port = "9987"
	} else if err != nil {
		return nil, err
	}

	_, records, err := net.LookupSRV("ts3", "udp", hostname)
	if err != nil {
		return nil, err
	}

	if len(records) > 0 {
		hosts := make([]string, len(records))
		for i, record := range records {
			hosts[i] = net.JoinHostPort(record.Target, strconv.FormatInt(int64(record.Port), 10))
		}

		return hosts, nil
	}

	addrs, err := net.LookupIP(hostname)
	if err != nil {
		return nil, err
	}

	hosts := make([]string, len(addrs))
	for i, addr := range addrs {
		hosts[i] = net.JoinHostPort(addr.String(), port)
	}

	return hosts, nil
}
