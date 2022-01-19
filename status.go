package mcstatus

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/rand"
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

type rawJavaStatus struct {
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
	MOTD      MOTD          `json:"motd"`
	Favicon   Favicon       `json:"favicon"`
	SRVResult *SRVRecord    `json:"srv_result"`
	Latency   time.Duration `json:"latency"`
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

	r := bufio.NewReader(conn)

	if err = conn.SetDeadline(time.Now().Add(opts.Timeout)); err != nil {
		return nil, err
	}

	// Handshake packet
	// https://wiki.vg/Server_List_Ping#Handshake
	{
		buf := &bytes.Buffer{}

		// Packet ID - varint
		if _, err := writeVarInt(0x00, buf); err != nil {
			return nil, err
		}

		// Protocol version - varint
		if _, err = writeVarInt(int32(opts.ProtocolVersion), buf); err != nil {
			return nil, err
		}

		// Host - string
		if err := writeString(host, buf); err != nil {
			return nil, err
		}

		// Port - uint16
		if err := binary.Write(buf, binary.BigEndian, port); err != nil {
			return nil, err
		}

		// Next state - varint
		if _, err := writeVarInt(1, buf); err != nil {
			return nil, err
		}

		if err := writePacket(buf, conn); err != nil {
			return nil, err
		}
	}

	// Request packet
	// https://wiki.vg/Server_List_Ping#Request
	{
		buf := &bytes.Buffer{}

		// Packet ID - varint
		if _, err := writeVarInt(0x00, buf); err != nil {
			return nil, err
		}

		if err := writePacket(buf, conn); err != nil {
			return nil, err
		}
	}

	var result rawJavaStatus

	// Response packet
	// https://wiki.vg/Server_List_Ping#Response
	{
		// Packet length - varint
		{
			if _, _, err := readVarInt(r); err != nil {
				return nil, err
			}
		}

		// Packet type - varint
		{
			packetType, _, err := readVarInt(r)

			if err != nil {
				return nil, err
			}

			if packetType != 0x00 {
				return nil, ErrUnexpectedResponse
			}
		}

		// Data - string
		{
			data, err := readString(r)

			if err != nil {
				return nil, err
			}

			if err = json.Unmarshal(data, &result); err != nil {
				return nil, err
			}

		}
	}

	payload := rand.Int63()

	// Ping packet
	// https://wiki.vg/Server_List_Ping#Ping
	{
		buf := &bytes.Buffer{}

		// Packet ID - varint
		if _, err := writeVarInt(0x01, buf); err != nil {
			return nil, err
		}

		// Payload - int64
		if err := binary.Write(buf, binary.BigEndian, payload); err != nil {
			return nil, err
		}

		if err := writePacket(buf, conn); err != nil {
			return nil, err
		}
	}

	pingStart := time.Now()

	// Pong packet
	// https://wiki.vg/Server_List_Ping#Pong
	{
		// Packet length - varint
		{
			if _, _, err := readVarInt(r); err != nil {
				return nil, err
			}
		}

		// Packet type - varint
		{
			packetType, _, err := readVarInt(r)

			if err != nil {
				return nil, err
			}

			if packetType != 0x01 {
				return nil, ErrUnexpectedResponse
			}
		}

		// Payload - int64
		{
			var returnPayload int64

			if err := binary.Read(r, binary.BigEndian, &returnPayload); err != nil {
				return nil, err
			}

			if payload != returnPayload {
				return nil, ErrUnexpectedResponse
			}
		}
	}

	motd, err := ParseMOTD(result.Description)

	if err != nil {
		return nil, err
	}

	return &JavaStatusResponse{
		Version:   result.Version,
		Players:   result.Players,
		MOTD:      *motd,
		Favicon:   parseFavicon(result.Favicon),
		SRVResult: srvResult,
		Latency:   time.Since(pingStart),
	}, nil
}

func parseJavaStatusOptions(opts ...JavaStatusOptions) JavaStatusOptions {
	if len(opts) < 1 {
		return defaultJavaStatusOptions
	}

	return opts[0]
}
