package mcstatus

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

var (
	defaultBedrockStatusOptions = BedrockStatusOptions{
		EnableSRV:  true,
		Timeout:    time.Second * 5,
		ClientGUID: 0,
	}
	bedrockMagic = []byte{0x00, 0xFF, 0xFF, 0x00, 0xFE, 0xFE, 0xFE, 0xFE, 0xFD, 0xFD, 0xFD, 0xFD, 0x12, 0x34, 0x56, 0x78}
)

type BedrockStatusResponse struct {
	ServerGUID      int64       `json:"server_guid"`
	Edition         string      `json:"edition"`
	MOTDLine1       Description `json:"motd_line_1"`
	MOTDLine2       Description `json:"motd_line_2"`
	ProtocolVersion int64       `json:"protocol_version"`
	Version         string      `json:"version"`
	OnlinePlayers   int64       `json:"online_players"`
	MaxPlayers      int64       `json:"max_players"`
	ServerID        uint64      `json:"server_id"`
	Gamemode        string      `json:"gamemode"`
	GamemodeID      int64       `json:"gamemode_id"`
	PortIPv4        uint16      `json:"port_ipv4"`
	PortIPv6        uint16      `json:"port_ipv6"`
	SRVResult       *SRVRecord  `json:"srv_result"`
}

type BedrockStatusOptions struct {
	EnableSRV  bool
	Timeout    time.Duration
	ClientGUID int64
}

// StatusBedrock retrieves the status of a Bedrock Minecraft server
func StatusBedrock(host string, port uint16, options ...BedrockStatusOptions) (*BedrockStatusResponse, error) {
	opts := parseBedrockStatusOptions(options...)

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

	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", host, port))

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	r := bufio.NewReader(conn)

	conn.SetDeadline(time.Now().Add(opts.Timeout))

	// Unconnected ping packet
	// https://wiki.vg/Raknet_Protocol#Unconnected_Ping
	{
		packet := NewPacket()

		// Packet ID (byte) - 0x01
		if err = packet.WriteByte(0x01); err != nil {
			return nil, err
		}

		// Time (int64)
		if err = packet.WriteInt64BE(time.Now().UnixNano() / int64(time.Millisecond)); err != nil {
			return nil, err
		}

		// Magic (bytes)
		if err = packet.WriteBytes(bedrockMagic); err != nil {
			return nil, err
		}

		// Client GUID (int64)
		if err = packet.WriteInt64BE(opts.ClientGUID); err != nil {
			return nil, err
		}

		if _, err = packet.WriteTo(conn); err != nil {
			return nil, err
		}
	}

	var serverGUID int64
	var response string

	// Unconnected pong packet
	// https://wiki.vg/Raknet_Protocol#Unconnected_Pong
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

			if data[0] != 0x1C {
				return nil, ErrUnexpectedResponse
			}
		}

		// Time - int64
		{
			data := make([]byte, 8)

			n, err := r.Read(data)

			if err != nil {
				return nil, err
			}

			if n < 8 {
				return nil, io.EOF
			}
		}

		// Server GUID - int64
		{
			data := make([]byte, 8)

			n, err := r.Read(data)

			if err != nil {
				return nil, err
			}

			if n < 8 {
				return nil, io.EOF
			}

			serverGUID = int64(binary.BigEndian.Uint64(data))
		}

		// Magic - bytes
		{
			data := make([]byte, 16)

			n, err := r.Read(data)

			if err != nil {
				return nil, err
			}

			if n < 16 {
				return nil, io.EOF
			}
		}

		// Server ID - string
		{
			lengthData := make([]byte, 2)

			n, err := r.Read(lengthData)

			if err != nil {
				return nil, err
			}

			if n < 2 {
				return nil, io.EOF
			}

			length := binary.BigEndian.Uint16(lengthData)

			data := make([]byte, length)

			n, err = r.Read(data)

			if err != nil {
				return nil, err
			}

			if n < int(length) {
				return nil, io.EOF
			}

			response = string(data)
		}
	}

	splitResponse := strings.Split(response, ";")

	if len(splitResponse) < 12 {
		return nil, ErrUnexpectedResponse
	}

	protocolVersion, err := strconv.ParseInt(splitResponse[2], 10, 64)

	if err != nil {
		return nil, err
	}

	onlinePlayers, err := strconv.ParseInt(splitResponse[4], 10, 64)

	if err != nil {
		return nil, err
	}

	maxPlayers, err := strconv.ParseInt(splitResponse[5], 10, 64)

	if err != nil {
		return nil, err
	}

	serverID, err := strconv.ParseUint(splitResponse[6], 10, 64)

	if err != nil {
		return nil, err
	}

	gamemodeID, err := strconv.ParseInt(splitResponse[9], 10, 64)

	if err != nil {
		return nil, err
	}

	portIPv4, err := strconv.ParseInt(splitResponse[10], 10, 64)

	if err != nil {
		return nil, err
	}

	portIPv6, err := strconv.ParseInt(splitResponse[11], 10, 64)

	if err != nil {
		return nil, err
	}

	return &BedrockStatusResponse{
		ServerGUID:      serverGUID,
		Edition:         splitResponse[0],
		MOTDLine1:       parseDescription(splitResponse[1]),
		MOTDLine2:       parseDescription(splitResponse[7]),
		ProtocolVersion: protocolVersion,
		Version:         splitResponse[3],
		OnlinePlayers:   onlinePlayers,
		MaxPlayers:      maxPlayers,
		ServerID:        serverID,
		Gamemode:        splitResponse[8],
		GamemodeID:      gamemodeID,
		PortIPv4:        uint16(portIPv4),
		PortIPv6:        uint16(portIPv6),
		SRVResult:       srvResult,
	}, nil
}

func parseBedrockStatusOptions(opts ...BedrockStatusOptions) BedrockStatusOptions {
	if len(opts) < 1 {
		options := BedrockStatusOptions(defaultBedrockStatusOptions)

		options.ClientGUID = rand.Int63()

		return options
	}

	return opts[0]
}
