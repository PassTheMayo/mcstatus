package mcstatus

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

var (
	defaultJavaStatusOptions = JavaStatusOptions{
		EnableSRV:       true,
		Timeout:         time.Second * 5,
		ProtocolVersion: 47,
	}
)

type rawStatus struct {
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

type JavaStatusOptions struct {
	EnableSRV       bool
	Timeout         time.Duration
	ProtocolVersion int
}

// Status retrieves the status of any Minecraft server
func Status(host string, port uint16, options ...JavaStatusOptions) (*JavaStatusResponse, error) {
	opts := parseJavaStatusOptions(options...)

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
		packet := NewPacket()

		// Packet ID (varint) - 0x00
		if err = packet.WriteVarInt(0); err != nil {
			return nil, err
		}

		// Protocol version (varint)
		if err = packet.WriteVarInt(int32(opts.ProtocolVersion)); err != nil {
			return nil, err
		}

		// Host (string)
		if err = packet.WriteString(host); err != nil {
			return nil, err
		}

		// Port (uint16)
		if err = packet.WriteUInt16BE(port); err != nil {
			return nil, err
		}

		// Next state (varint)
		if err = packet.WriteVarInt(1); err != nil {
			return nil, err
		}

		if err = packet.WriteLength(); err != nil {
			return nil, err
		}

		if _, err = packet.WriteTo(conn); err != nil {
			return nil, err
		}
	}

	// Request packet
	// https://wiki.vg/Server_List_Ping#Request
	{
		packet := NewPacket()

		// Packet ID (varint) - 0x00
		if err = packet.WriteVarInt(0); err != nil {
			return nil, err
		}

		if err = packet.WriteLength(); err != nil {
			return nil, err
		}

		if _, err = packet.WriteTo(conn); err != nil {
			return nil, err
		}
	}

	// Response packet
	// https://wiki.vg/Server_List_Ping#Response
	{
		// Packet length - varint
		{
			if _, _, err := readVarInt(conn); err != nil {
				return nil, err
			}
		}

		// Packet type - varint
		{
			packetType, _, err := readVarInt(conn)

			if err != nil {
				return nil, err
			}

			if packetType != 0 {
				return nil, ErrUnexpectedResponse
			}
		}

		// Data - string
		{
			data, err := readString(conn)

			if err != nil {
				return nil, err
			}

			result := &rawStatus{}

			if err = json.Unmarshal([]byte(data), result); err != nil {
				return nil, err
			}

			description, err := NewDescription(result.Description)

			if err != nil {
				return nil, err
			}

			return &JavaStatusResponse{
				Version:     result.Version,
				Players:     result.Players,
				Description: *description,
				Favicon:     parseFavicon(result.Favicon),
				SRVResult:   srvResult,
			}, nil
		}
	}
}

func parseJavaStatusOptions(opts ...JavaStatusOptions) JavaStatusOptions {
	if len(opts) < 1 {
		return defaultJavaStatusOptions
	}

	return opts[0]
}
