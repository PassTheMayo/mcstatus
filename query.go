package mcstatus

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
)

var (
	sessionID           int32 = 0
	defaultQueryOptions       = QueryOptions{
		Timeout:   time.Second * 5,
		SessionID: 0,
	}
	magic = []byte{0xFE, 0xFD}
)

type QueryOptions struct {
	Timeout   time.Duration
	SessionID int32
}

type BasicQueryResponse struct {
	MOTD          MOTD
	GameType      string
	Map           string
	OnlinePlayers uint64
	MaxPlayers    uint64
	HostPort      uint16
	HostIP        string
}

type FullQueryResponse struct {
	Data    map[string]string
	Players []string
}

// BasicQuery runs a query on the server and returns basic information
func BasicQuery(host string, port uint16, options ...QueryOptions) (*BasicQueryResponse, error) {
	opts := parseQueryOptions(options...)

	conn, err := net.DialTimeout("udp", fmt.Sprintf("%s:%d", host, port), opts.Timeout)

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	r := bufio.NewReader(conn)

	if err = conn.SetDeadline(time.Now().Add(opts.Timeout)); err != nil {
		return nil, err
	}

	// Handshake request packet
	// https://wiki.vg/Query#Request
	{
		buf := &bytes.Buffer{}

		// Magic - uint16
		if _, err := buf.Write(magic); err != nil {
			return nil, err
		}

		// Type - byte
		if err := buf.WriteByte(0x09); err != nil {
			return nil, err
		}

		// Session ID - int32
		if err := binary.Write(buf, binary.BigEndian, opts.SessionID&0x0F0F0F0F); err != nil {
			return nil, err
		}

		if _, err := io.Copy(conn, buf); err != nil {
			return nil, err
		}
	}

	var challengeToken int32

	// Handshake response packet
	// https://wiki.vg/Query#Response
	{
		// Type - byte
		{
			v, err := r.ReadByte()

			if err != nil {
				return nil, err
			}

			if v != 0x09 {
				return nil, ErrUnexpectedResponse
			}
		}

		// Session ID - int32
		{
			var sessionID int32

			if err := binary.Read(r, binary.BigEndian, &sessionID); err != nil {
				return nil, err
			}

			if sessionID != opts.SessionID {
				return nil, ErrUnexpectedResponse
			}
		}

		// Challenge Token - string
		{
			data, err := r.ReadBytes(0x00)

			if err != nil {
				return nil, err
			}

			v, err := strconv.ParseInt(string(data[:len(data)-1]), 10, 32)

			if err != nil {
				return nil, err
			}

			challengeToken = int32(v)
		}
	}

	// Basic stat request packet
	// https://wiki.vg/Query#Request_2
	{
		buf := &bytes.Buffer{}

		// Magic - uint16
		if _, err := buf.Write(magic); err != nil {
			return nil, err
		}

		// Type - byte
		if err := buf.WriteByte(0x00); err != nil {
			return nil, err
		}

		// Session ID - int32
		if err := binary.Write(buf, binary.BigEndian, opts.SessionID&0x0F0F0F0F); err != nil {
			return nil, err
		}

		// Challenge Token - int32
		if err := binary.Write(buf, binary.BigEndian, challengeToken); err != nil {
			return nil, err
		}

		if _, err := io.Copy(conn, buf); err != nil {
			return nil, err
		}
	}

	response := BasicQueryResponse{}

	// Basic stat response packet
	// https://wiki.vg/Query#Response_2
	{
		// Type - byte
		{
			v, err := r.ReadByte()

			if err != nil {
				return nil, err
			}

			if v != 0x00 {
				return nil, ErrUnexpectedResponse
			}
		}

		// Session ID - int32
		{
			var sessionID int32

			if err := binary.Read(r, binary.BigEndian, &sessionID); err != nil {
				return nil, err
			}

			if sessionID != opts.SessionID {
				return nil, ErrUnexpectedResponse
			}
		}

		// MOTD - null-terminated string
		{
			data, err := r.ReadBytes(0x00)

			if err != nil {
				return nil, err
			}

			description, err := ParseMOTD(decodeASCII(data[:len(data)-1]))

			if err != nil {
				return nil, err
			}

			response.MOTD = *description
		}

		// Game Type - null-terminated string
		{
			data, err := r.ReadBytes(0x00)

			if err != nil {
				return nil, err
			}

			response.GameType = string(data[:len(data)-1])
		}

		// Map - null-terminated string
		{
			data, err := r.ReadBytes(0x00)

			if err != nil {
				return nil, err
			}

			response.Map = string(data[:len(data)-1])
		}

		// Online Players - null-terminated string
		{
			data, err := r.ReadBytes(0x00)

			if err != nil {
				return nil, err
			}

			onlinePlayers, err := strconv.ParseUint(string(data[:len(data)-1]), 10, 64)

			if err != nil {
				return nil, err
			}

			response.OnlinePlayers = onlinePlayers
		}

		// Max Players - null-terminated string
		{
			data, err := r.ReadBytes(0x00)

			if err != nil {
				return nil, err
			}

			maxPlayers, err := strconv.ParseUint(string(data[:len(data)-1]), 10, 64)

			if err != nil {
				return nil, err
			}

			response.MaxPlayers = maxPlayers
		}

		// Host Port - uint16
		{
			var hostPort uint16

			if err := binary.Read(r, binary.LittleEndian, &hostPort); err != nil {
				return nil, err
			}

			response.HostPort = hostPort
		}

		// Host IP - null-terminated string
		{
			data, err := r.ReadBytes(0x00)

			if err != nil {
				return nil, err
			}

			response.HostIP = string(data[:len(data)-1])
		}
	}

	return &response, nil
}

// FullQuery runs a query on the server and returns the full information
func FullQuery(host string, port uint16, options ...QueryOptions) (*FullQueryResponse, error) {
	opts := parseQueryOptions(options...)

	conn, err := net.DialTimeout("udp", fmt.Sprintf("%s:%d", host, port), opts.Timeout)

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	r := bufio.NewReader(conn)

	if err = conn.SetDeadline(time.Now().Add(opts.Timeout)); err != nil {
		return nil, err
	}

	// Handshake request packet
	// https://wiki.vg/Query#Request
	{
		buf := &bytes.Buffer{}

		// Magic - uint16
		if _, err := buf.Write(magic); err != nil {
			return nil, err
		}

		// Type - byte
		if err := buf.WriteByte(0x09); err != nil {
			return nil, err
		}

		// Session ID - int32
		if err := binary.Write(buf, binary.BigEndian, opts.SessionID&0x0F0F0F0F); err != nil {
			return nil, err
		}

		if _, err := io.Copy(conn, buf); err != nil {
			return nil, err
		}
	}

	var challengeToken int32

	// Handshake response packet
	// https://wiki.vg/Query#Response
	{
		// Type - byte
		{
			v, err := r.ReadByte()

			if err != nil {
				return nil, err
			}

			if v != 0x09 {
				return nil, ErrUnexpectedResponse
			}
		}

		// Session ID - int32
		{
			var sessionID int32

			if err := binary.Read(r, binary.BigEndian, &sessionID); err != nil {
				return nil, err
			}

			if sessionID != opts.SessionID {
				return nil, ErrUnexpectedResponse
			}
		}

		// Challenge Token - null-terminated string
		{
			data, err := r.ReadBytes(0x00)

			if err != nil {
				return nil, err
			}

			v, err := strconv.ParseInt(string(data[:len(data)-1]), 10, 32)

			if err != nil {
				return nil, err
			}

			challengeToken = int32(v)
		}
	}

	// Full stat request packet
	// https://wiki.vg/Query#Request_2
	{
		buf := &bytes.Buffer{}

		// Magic - uint16
		if _, err := buf.Write(magic); err != nil {
			return nil, err
		}

		// Type - byte
		if err := buf.WriteByte(0x00); err != nil {
			return nil, err
		}

		// Session ID - int32
		if err := binary.Write(buf, binary.BigEndian, opts.SessionID&0x0F0F0F0F); err != nil {
			return nil, err
		}

		// Challenge Token - int32
		if err := binary.Write(buf, binary.BigEndian, challengeToken); err != nil {
			return nil, err
		}

		// Padding - bytes
		if _, err := buf.Write([]byte{0x00, 0x00, 0x00, 0x00}); err != nil {
			return nil, err
		}

		if _, err := io.Copy(conn, buf); err != nil {
			return nil, err
		}
	}

	response := FullQueryResponse{
		Data:    make(map[string]string),
		Players: make([]string, 0),
	}

	// Full stat response packet
	// https://wiki.vg/Query#Response_3
	{
		// Type - byte
		{
			v, err := r.ReadByte()

			if err != nil {
				return nil, err
			}

			if v != 0x00 {
				return nil, ErrUnexpectedResponse
			}
		}

		// Session ID - int16
		{
			var sessionID int32

			if err := binary.Read(r, binary.BigEndian, &sessionID); err != nil {
				return nil, err
			}

			if sessionID != opts.SessionID {
				return nil, ErrUnexpectedResponse
			}
		}

		// Padding - [11]byte
		{
			data := make([]byte, 11)

			if _, err := r.Read(data); err != nil {
				return nil, err
			}
		}

		// K, V section - null-terminated key,pair pair string
		{
			for {
				data, err := r.ReadBytes(0x00)

				if err != nil {
					return nil, err
				}

				if len(data) < 2 {
					break
				}

				key := decodeASCII(data[:len(data)-1])

				data, err = r.ReadBytes(0x00)

				if err != nil {
					return nil, err
				}

				value := decodeASCII(data[:len(data)-1])

				response.Data[key] = value
			}
		}

		// Padding - [10]byte
		{
			data := make([]byte, 10)

			if _, err := r.Read(data); err != nil {
				return nil, err
			}
		}

		// Players section - null-terminated key,value pair string
		{
			for {
				data, err := r.ReadBytes(0x00)

				if err != nil {
					return nil, err
				}

				if len(data) < 2 {
					break
				}

				response.Players = append(response.Players, string(data[:len(data)-1]))
			}
		}
	}

	return &response, nil
}

func parseQueryOptions(opts ...QueryOptions) QueryOptions {
	if len(opts) < 1 {
		options := QueryOptions(defaultQueryOptions)

		sessionID += 1

		options.SessionID = sessionID

		return options
	}

	return opts[0]
}
