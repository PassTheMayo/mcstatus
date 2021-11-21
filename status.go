package mcstatus

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

var (
	defaultStatusOptions = StatusOptions{
		EnableSRV:       true,
		Timeout:         time.Second * 5,
		ProtocolVersion: 47,
	}
)

type rawStatusResponse struct {
	Version struct {
		Name     string `json:"name"`
		Protocol int    `json:"protocol"`
	} `json:"version"`
	Players struct {
		Max    int `json:"max"`
		Online int `json:"online"`
		Sample []struct {
			Name string `json:"name"`
			ID   string `json:"id"`
		} `json:"sample"`
	} `json:"players"`
	Description interface{} `json:"description"`
	Favicon     interface{} `json:"favicon"`
}

type JavaStatusResponse struct {
	Version struct {
		Name     string `json:"name"`
		Protocol int    `json:"protocol"`
	} `json:"version"`
	Players struct {
		Max    int `json:"max"`
		Online int `json:"online"`
		Sample []struct {
			Name string `json:"name"`
			ID   string `json:"id"`
		} `json:"sample"`
	} `json:"players"`
	Description Description `json:"description"`
	Favicon     Favicon     `json:"favicon"`
	SRVResult   *SRVRecord  `json:"srv_result"`
}

type SRVRecord struct {
	Host string `json:"host"`
	Port uint16 `json:"port"`
}

type StatusOptions struct {
	EnableSRV       bool
	Timeout         time.Duration
	ProtocolVersion int
}

// Status retrieves the status of any Minecraft server
func Status(host string, port uint16, options ...StatusOptions) (*JavaStatusResponse, error) {
	opts := parseStatusOptions(options...)

	var srvResult *SRVRecord = nil

	if opts.EnableSRV {
		record, err := lookupSRV(host, port)

		if err == nil && record != nil {
			host = record.Target
			port = record.Port

			srvResult = &SRVRecord{
				Host: record.Target,
				Port: record.Port,
			}
		}
	}

	conn, err := net.Dial("tcp4", fmt.Sprintf("%s:%d", host, port))

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	conn.SetDeadline(time.Now().Add(opts.Timeout))

	// Handshake packet
	// https://wiki.vg/Server_List_Ping#Handshake
	{
		handshakePacket := NewPacket()

		// Packet ID (varint) - 0x00
		if err = handshakePacket.WriteVarInt(0); err != nil {
			return nil, err
		}

		// Protocol version (varint)
		if err = handshakePacket.WriteVarInt(int32(opts.ProtocolVersion)); err != nil {
			return nil, err
		}

		// Host (string)
		if err = handshakePacket.WriteString(host); err != nil {
			return nil, err
		}

		// Port (uint16)
		if err = handshakePacket.WriteUnsignedShort(port); err != nil {
			return nil, err
		}

		// Next state (varint)
		if err = handshakePacket.WriteVarInt(1); err != nil {
			return nil, err
		}

		if err = handshakePacket.WriteLength(); err != nil {
			return nil, err
		}

		if _, err = handshakePacket.WriteTo(conn); err != nil {
			return nil, err
		}
	}

	// Request packet
	// https://wiki.vg/Server_List_Ping#Request
	{
		handshakePacket := NewPacket()

		// Packet ID (varint) - 0x00
		if err = handshakePacket.WriteVarInt(0); err != nil {
			return nil, err
		}

		if err = handshakePacket.WriteLength(); err != nil {
			return nil, err
		}

		if _, err = handshakePacket.WriteTo(conn); err != nil {
			return nil, err
		}
	}

	// Response packet
	// https://wiki.vg/Server_List_Ping#Response
	{
		if _, _, err := readVarInt(conn); err != nil {
			return nil, err
		}

		packetType, _, err := readVarInt(conn)

		if err != nil {
			return nil, err
		}

		if packetType != 0 {
			return nil, ErrUnexpectedResponse
		}

		data, err := readString(conn)

		if err != nil {
			return nil, err
		}

		result := &rawStatusResponse{}

		if err = json.Unmarshal([]byte(data), result); err != nil {
			return nil, err
		}

		return &JavaStatusResponse{
			Version:     result.Version,
			Players:     result.Players,
			Description: parseDescription(result.Description),
			Favicon:     parseFavicon(result.Favicon),
			SRVResult:   srvResult,
		}, nil
	}
}

func parseStatusOptions(opts ...StatusOptions) StatusOptions {
	if len(opts) < 1 {
		return defaultStatusOptions
	}

	return opts[0]
}
