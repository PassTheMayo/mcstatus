package mcstatus

import (
	"bufio"
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
	MOTD          Description
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

	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", host, port))

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	r := bufio.NewReader(conn)

	conn.SetDeadline(time.Now().Add(opts.Timeout))

	// Handshake request packet
	// https://wiki.vg/Query#Request
	{
		packet := NewPacket()

		// Magic (uint16) - 0xFEFD
		if err = packet.WriteBytes(magic); err != nil {
			return nil, err
		}

		// Type (byte) - 0x09
		if err = packet.WriteByte(0x09); err != nil {
			return nil, err
		}

		// Session ID (int32)
		if err = packet.WriteInt32BE(opts.SessionID & 0x0F0F0F0F); err != nil {
			return nil, err
		}

		n, err := packet.WriteTo(conn)

		if err != nil {
			return nil, err
		}

		if n != 7 {
			return nil, ErrUnexpectedResponse
		}
	}

	var challengeToken int32

	// Handshake response packet
	// https://wiki.vg/Query#Response
	{
		// Type - byte
		{
			data := make([]byte, 1)

			n, err := r.Read(data)

			if err != nil {
				return nil, err
			}

			if n < 1 {
				return nil, io.EOF
			}

			if data[0] != 0x09 {
				return nil, ErrUnexpectedResponse
			}
		}

		// Session ID - int32
		{
			data := make([]byte, 4)

			n, err := r.Read(data)

			if err != nil {
				return nil, err
			}

			if n < 4 {
				return nil, io.EOF
			}

			if int32(binary.BigEndian.Uint32(data)) != opts.SessionID {
				return nil, ErrUnexpectedResponse
			}
		}

		// Challenge Token - string
		{
			data := make([]byte, 1)

			var challengeTokenString string

			for {
				n, err := r.Read(data)

				if err != nil {
					return nil, err
				}

				if n < 1 {
					return nil, io.EOF
				}

				if data[0] == 0x00 {
					break
				}

				challengeTokenString += string(data[0])
			}

			v, err := strconv.ParseInt(challengeTokenString, 10, 32)

			if err != nil {
				return nil, err
			}

			challengeToken = int32(v)
		}
	}

	// Basic stat request packet
	// https://wiki.vg/Query#Request_2
	{
		packet := NewPacket()

		// Magic (uint16) - 0xFEFD
		if err = packet.WriteBytes(magic); err != nil {
			return nil, err
		}

		// Type (byte) - 0x00
		if err = packet.WriteByte(0x00); err != nil {
			return nil, err
		}

		// Session ID (int32)
		if err = packet.WriteInt32BE(opts.SessionID & 0x0F0F0F0F); err != nil {
			return nil, err
		}

		// Challenge Token (int32)
		if err = packet.WriteInt32BE(challengeToken); err != nil {
			return nil, err
		}

		n, err := packet.WriteTo(conn)

		if err != nil {
			return nil, err
		}

		if n != 11 {
			return nil, ErrUnexpectedResponse
		}
	}

	var (
		motd          string
		gameType      string
		mapName       string
		onlinePlayers uint64
		maxPlayers    uint64
		hostPort      uint16
		hostIP        string
	)

	// Basic stat response packet
	// https://wiki.vg/Query#Response_2
	{
		// Type - byte
		{
			data := make([]byte, 1)

			n, err := r.Read(data)

			if err != nil {
				return nil, err
			}

			if n < 1 {
				return nil, io.EOF
			}

			if data[0] != 0x00 {
				return nil, ErrUnexpectedResponse
			}
		}

		// Session ID - int32
		{
			data := make([]byte, 4)

			n, err := r.Read(data)

			if err != nil {
				return nil, err
			}

			if n < 4 {
				return nil, io.EOF
			}

			if int32(binary.BigEndian.Uint32(data)) != opts.SessionID {
				return nil, ErrUnexpectedResponse
			}
		}

		// MOTD - null-terminated string
		{
			data := make([]byte, 1)

			for {
				n, err := r.Read(data)

				if err != nil {
					return nil, err
				}

				if n < 1 {
					return nil, io.EOF
				}

				if data[0] == 0x00 {
					break
				}

				motd += string(data[0])
			}
		}

		// Game Type - null-terminated string
		{
			data := make([]byte, 1)

			for {
				n, err := r.Read(data)

				if err != nil {
					return nil, err
				}

				if n < 1 {
					return nil, io.EOF
				}

				if data[0] == 0x00 {
					break
				}

				gameType += string(data[0])
			}
		}

		// Map - null-terminated string
		{
			data := make([]byte, 1)

			for {
				n, err := r.Read(data)

				if err != nil {
					return nil, err
				}

				if n < 1 {
					return nil, io.EOF
				}

				if data[0] == 0x00 {
					break
				}

				mapName += string(data[0])
			}
		}

		// Online Players - null-terminated string
		{
			data := make([]byte, 1)

			var onlinePlayersString string

			for {
				n, err := r.Read(data)

				if err != nil {
					return nil, err
				}

				if n < 1 {
					return nil, io.EOF
				}

				if data[0] == 0x00 {
					break
				}

				onlinePlayersString += string(data[0])
			}

			onlinePlayers, err = strconv.ParseUint(onlinePlayersString, 10, 64)

			if err != nil {
				return nil, err
			}
		}

		// Max Players - null-terminated string
		{
			data := make([]byte, 1)

			var maxPlayersString string

			for {
				n, err := r.Read(data)

				if err != nil {
					return nil, err
				}

				if n < 1 {
					return nil, io.EOF
				}

				if data[0] == 0x00 {
					break
				}

				maxPlayersString += string(data[0])
			}

			maxPlayers, err = strconv.ParseUint(maxPlayersString, 10, 64)

			if err != nil {
				return nil, err
			}
		}

		// Host Port - uint16
		{
			data := make([]byte, 2)

			n, err := r.Read(data)

			if err != nil {
				return nil, err
			}

			if n != 2 {
				return nil, io.EOF
			}

			hostPort = binary.LittleEndian.Uint16(data)
		}

		// Host IP - null-terminated string
		{
			data := make([]byte, 1)

			for {
				n, err := r.Read(data)

				if err != nil {
					return nil, err
				}

				if n < 1 {
					return nil, io.EOF
				}

				if data[0] == 0x00 {
					break
				}

				hostIP += string(data[0])
			}
		}
	}

	description, err := NewDescription(motd)

	if err != nil {
		return nil, err
	}

	return &BasicQueryResponse{
		MOTD:          *description,
		GameType:      gameType,
		Map:           mapName,
		OnlinePlayers: onlinePlayers,
		MaxPlayers:    maxPlayers,
		HostPort:      hostPort,
		HostIP:        hostIP,
	}, nil
}

// FullQuery runs a query on the server and returns the full information
func FullQuery(host string, port uint16, options ...QueryOptions) (*FullQueryResponse, error) {
	opts := parseQueryOptions(options...)

	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", host, port))

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	r := bufio.NewReader(conn)

	conn.SetDeadline(time.Now().Add(opts.Timeout))

	// Handshake request packet
	// https://wiki.vg/Query#Request
	{
		handshakePacket := NewPacket()

		// Magic (uint16) - 0xFEFD
		if err = handshakePacket.WriteBytes(magic); err != nil {
			return nil, err
		}

		// Type (byte) - 0x09
		if err = handshakePacket.WriteByte(0x09); err != nil {
			return nil, err
		}

		// Session ID (int32)
		if err = handshakePacket.WriteInt32BE(opts.SessionID & 0x0F0F0F0F); err != nil {
			return nil, err
		}

		n, err := handshakePacket.WriteTo(conn)

		if err != nil {
			return nil, err
		}

		if n != 7 {
			return nil, ErrUnexpectedResponse
		}
	}

	var challengeToken int32

	// Handshake response packet
	// https://wiki.vg/Query#Response
	{
		// Type - byte
		{
			data := make([]byte, 1)

			n, err := r.Read(data)

			if err != nil {
				return nil, err
			}

			if n < 1 {
				return nil, io.EOF
			}

			if data[0] != 0x09 {
				return nil, ErrUnexpectedResponse
			}
		}

		// Session ID - uint16
		{
			data := make([]byte, 4)

			n, err := r.Read(data)

			if err != nil {
				return nil, err
			}

			if n < 4 {
				return nil, io.EOF
			}

			if int32(binary.BigEndian.Uint32(data)) != opts.SessionID {
				return nil, ErrUnexpectedResponse
			}
		}

		// Challenge Token - null-terminated string
		{
			data := make([]byte, 1)

			var challengeTokenString string

			for {
				n, err := r.Read(data)

				if err != nil {
					return nil, err
				}

				if n < 1 {
					return nil, io.EOF
				}

				if data[0] == 0x00 {
					break
				}

				challengeTokenString += string(data[0])
			}

			v, err := strconv.ParseInt(challengeTokenString, 10, 32)

			if err != nil {
				return nil, err
			}

			challengeToken = int32(v)
		}
	}

	// Full stat request packet
	// https://wiki.vg/Query#Request_2
	{
		requestPacket := NewPacket()

		// Magic (uint16) - 0xFEFD
		if err = requestPacket.WriteBytes(magic); err != nil {
			return nil, err
		}

		// Type (byte) - 0x00
		if err = requestPacket.WriteByte(0x00); err != nil {
			return nil, err
		}

		// Session ID (int32)
		if err = requestPacket.WriteInt32BE(opts.SessionID & 0x0F0F0F0F); err != nil {
			return nil, err
		}

		// Challenge Token (int32)
		if err = requestPacket.WriteInt32BE(challengeToken); err != nil {
			return nil, err
		}

		// Padding ([4]byte)
		if err = requestPacket.WriteBytes([]byte{0x00, 0x00, 0x00, 0x00}); err != nil {
			return nil, err
		}

		n, err := requestPacket.WriteTo(conn)

		if err != nil {
			return nil, err
		}

		if n != 15 {
			return nil, ErrUnexpectedResponse
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
			data := make([]byte, 1)

			n, err := r.Read(data)

			if err != nil {
				return nil, err
			}

			if n < 1 {
				return nil, io.EOF
			}

			if data[0] != 0x00 {
				return nil, ErrUnexpectedResponse
			}
		}

		// Session ID - int16
		{
			data := make([]byte, 4)

			n, err := r.Read(data)

			if err != nil {
				return nil, err
			}

			if n < 4 {
				return nil, io.EOF
			}

			if int32(binary.BigEndian.Uint32(data)) != opts.SessionID {
				return nil, ErrUnexpectedResponse
			}
		}

		// Padding - [11]byte
		{
			data := make([]byte, 11)

			n, err := r.Read(data)

			if err != nil {
				return nil, err
			}

			if n < 11 {
				return nil, io.EOF
			}
		}

		// K, V section - null-terminated key,pair pair string
		{
		dataLoop:
			for {
				var key string
				var value string
				data := make([]byte, 1)

			keyLoop:
				for {
					n, err := r.Read(data)

					if err != nil {
						return nil, err
					}

					if n < 1 {
						return nil, io.EOF
					}

					if data[0] == 0x00 {
						if len(key) < 1 {
							break dataLoop
						}

						break keyLoop
					}

					key += string(data[0])
				}

			valueLoop:
				for {
					n, err := r.Read(data)

					if err != nil {
						return nil, err
					}

					if n < 1 {
						return nil, io.EOF
					}

					if data[0] == 0x00 {
						break valueLoop
					}

					value += string(data[0])
				}

				response.Data[key] = value
			}
		}

		// Padding - [10]byte
		{
			data := make([]byte, 10)

			n, err := r.Read(data)

			if err != nil {
				return nil, err
			}

			if n < 10 {
				return nil, io.EOF
			}
		}

		// Players section - null-terminated key,value pair string
		{
			var username string
			data := make([]byte, 1)

			for {
				n, err := r.Read(data)

				if err != nil {
					return nil, err
				}

				if n < 1 {
					return nil, io.EOF
				}

				if data[0] == 0x00 {
					if len(username) < 1 {
						break
					} else {
						response.Players = append(response.Players, username)

						username = ""
					}
				} else {
					username += string(data[0])
				}
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
