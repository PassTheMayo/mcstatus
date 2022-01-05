package mcstatus

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
	"unicode/utf16"
)

var (
	defaultJavaStatusLegacyOptions = JavaStatusLegacyOptions{
		EnableSRV: true,
		Timeout:   time.Second * 5,
	}
)

type JavaStatusLegacyResponse struct {
	Version   *JavaStatusLegacyVersion `json:"version"`
	Players   JavaStatusLegacyPlayers  `json:"players"`
	MOTD      MOTD                     `json:"motd"`
	SRVResult *SRVRecord               `json:"srv_result"`
}

type JavaStatusLegacyVersion struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}

type JavaStatusLegacyPlayers struct {
	Online int `json:"online"`
	Max    int `json:"max"`
}

type JavaStatusLegacyOptions struct {
	EnableSRV       bool
	Timeout         time.Duration
	ProtocolVersion int
}

func StatusLegacy(host string, port uint16, options ...JavaStatusLegacyOptions) (*JavaStatusLegacyResponse, error) {
	opts := parseJavaStatusLegacyOptions(options...)

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

	// Client to server packet
	// https://wiki.vg/Server_List_Ping#Client_to_server
	{
		if _, err = conn.Write([]byte{0xFE, 0x01}); err != nil {
			return nil, err
		}
	}

	// Server to client packet
	// https://wiki.vg/Server_List_Ping#Server_to_client
	{
		packetType, err := r.ReadByte()

		if err != nil {
			return nil, err
		}

		if packetType != 0xFF {
			return nil, fmt.Errorf("unexpected packet type returned from server: 0x%X", packetType)
		}

		var length uint16

		if err = binary.Read(r, binary.BigEndian, &length); err != nil {
			return nil, err
		}

		data := make([]byte, length*2)

		if _, err = r.Read(data); err != nil {
			return nil, err
		}

		byteData := make([]uint16, length)

		for i, l := 0, len(data); i < l; i += 2 {
			byteData[i/2] = (uint16(data[i]) << 8) | uint16(data[i+1])
		}

		response := string(utf16.Decode(byteData))

		if byteData[0] == 0x00A7 && byteData[1] == 0x0031 {
			// 1.4+ server

			split := strings.Split(response, "\x00")

			protocolVersion, err := strconv.ParseInt(split[1], 10, 32)

			if err != nil {
				return nil, err
			}

			onlinePlayers, err := strconv.ParseInt(split[4], 10, 32)

			if err != nil {
				return nil, err
			}

			maxPlayers, err := strconv.ParseInt(split[5], 10, 32)

			if err != nil {
				return nil, err
			}

			motd, err := ParseMOTD(split[3])

			if err != nil {
				return nil, err
			}

			return &JavaStatusLegacyResponse{
				Version: &JavaStatusLegacyVersion{
					Name:     split[2],
					Protocol: int(protocolVersion),
				},
				Players: JavaStatusLegacyPlayers{
					Online: int(onlinePlayers),
					Max:    int(maxPlayers),
				},
				MOTD:      *motd,
				SRVResult: srvResult,
			}, nil
		} else {
			// < 1.4 server

			split := strings.Split(response, "\u00A7")

			onlinePlayers, err := strconv.ParseInt(split[1], 10, 32)

			if err != nil {
				return nil, err
			}

			maxPlayers, err := strconv.ParseInt(split[2], 10, 32)

			if err != nil {
				return nil, err
			}

			motd, err := ParseMOTD(split[0])

			if err != nil {
				return nil, err
			}

			return &JavaStatusLegacyResponse{
				Version: nil,
				Players: JavaStatusLegacyPlayers{
					Online: int(onlinePlayers),
					Max:    int(maxPlayers),
				},
				MOTD:      *motd,
				SRVResult: srvResult,
			}, nil
		}
	}
}

func parseJavaStatusLegacyOptions(opts ...JavaStatusLegacyOptions) JavaStatusLegacyOptions {
	if len(opts) < 1 {
		return defaultJavaStatusLegacyOptions
	}

	return opts[0]
}
