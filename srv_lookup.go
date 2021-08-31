package mcstatus

import (
	"fmt"
	"net"
)

func lookupSRV(host string, port uint16) (*net.SRV, error) {
	_, addrs, err := net.LookupSRV("minecraft", "tcp", fmt.Sprintf("%s:%d", host, port))

	if err != nil {
		return nil, err
	}

	if len(addrs) < 1 {
		return nil, nil
	}

	return addrs[0], nil
}
