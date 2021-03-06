package mcstatus

import (
	"net"
)

type SRVRecord struct {
	Host string `json:"host"`
	Port uint16 `json:"port"`
}

func lookupSRV(host string, port uint16) (*net.SRV, error) {
	_, addrs, err := net.LookupSRV("minecraft", "tcp", host)

	if err != nil {
		return nil, err
	}

	if len(addrs) < 1 {
		return nil, nil
	}

	return addrs[0], nil
}
