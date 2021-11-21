package mcstatus

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

var (
	defaultRCONOptions = RCONOptions{
		Timeout: time.Second * 5,
	}
)

type RCON struct {
	Conn        *net.Conn
	r           *bufio.Reader
	Messages    chan string
	authSuccess bool
	requestID   int64
}

type RCONOptions struct {
	Timeout time.Duration
}

// NewRCON creates a new RCON client from the options parameter
func NewRCON() *RCON {
	return &RCON{
		Conn:        nil,
		r:           nil,
		Messages:    make(chan string),
		authSuccess: false,
		requestID:   0,
	}
}

func (r *RCON) Dial(host string, port uint16, options ...RCONOptions) error {
	opts := parseRCONOptions(options...)

	conn, err := net.Dial("tcp4", fmt.Sprintf("%s:%d", host, port))

	conn.SetDeadline(time.Now().Add(opts.Timeout))

	r.Conn = &conn
	r.r = bufio.NewReader(conn)

	return err
}

func (r *RCON) Login(password string) error {
	if r.Conn == nil {
		return ErrNotConnected
	}

	if r.authSuccess {
		return ErrAlreadyLoggedIn
	}

	// Login request packet
	// https://wiki.vg/RCON#3:_Login
	{
		loginPacket := NewPacket()

		// Length (int32)
		if err := loginPacket.WriteIntLE(int32(10 + len(password))); err != nil {
			return err
		}

		// Request ID (int32) - 0
		if err := loginPacket.WriteIntLE(0); err != nil {
			return err
		}

		// Type (int32) - 3
		if err := loginPacket.WriteIntLE(3); err != nil {
			return err
		}

		// Payload (null-terminated string) - 3
		if err := loginPacket.WriteBytes(append([]byte(password), 0x00)); err != nil {
			return err
		}

		// Padding (null byte) - 0x00
		if err := loginPacket.WriteByte(0x00); err != nil {
			return err
		}

		n, err := loginPacket.WriteTo(*r.Conn)

		if err != nil {
			return err
		}

		if n != int64(14+len(password)) {
			return ErrUnexpectedResponse
		}
	}

	// Login response packet
	// https://wiki.vg/RCON#3:_Login
	{
		var packetLength uint32

		// Length - int32
		{
			data := make([]byte, 4)

			n, err := (*r.r).Read(data)

			if err != nil {
				return err
			}

			if n < 4 {
				return io.EOF
			}

			packetLength = binary.LittleEndian.Uint32(data)

			if packetLength != 10 {
				return ErrUnexpectedResponse
			}
		}

		// Request ID - int32
		{
			data := make([]byte, 4)

			n, err := (*r.r).Read(data)

			if err != nil {
				return err
			}

			if n < 4 {
				return io.EOF
			}

			requestID := int32(binary.LittleEndian.Uint32(data))

			if requestID == -1 {
				return ErrInvalidPassword
			} else if requestID != 0 {
				return ErrUnexpectedResponse
			}
		}

		// Type - int32
		{
			data := make([]byte, 4)

			n, err := (*r.r).Read(data)

			if err != nil {
				return err
			}

			if n < 4 {
				return io.EOF
			}

			if binary.LittleEndian.Uint32(data) != 2 {
				return ErrUnexpectedResponse
			}
		}

		// Remaining bytes
		{
			data := make([]byte, packetLength-8)

			n, err := (*r.r).Read(data)

			if err != nil {
				return err
			}

			if n < int(packetLength-8) {
				return io.EOF
			}
		}
	}

	r.authSuccess = true

	(*r.Conn).SetReadDeadline(time.Time{})

	go (func() {
		<-time.NewTimer(time.Millisecond * 250).C

		for r.Conn != nil {
			r.readMessage()

			// TODO proper error handling of `r.readMessage()` but ignore 'use of closed network connection' when client is closed
		}
	})()

	return nil
}

func (r *RCON) Run(command string) error {
	if r.Conn == nil {
		return ErrNotConnected
	}

	if !r.authSuccess {
		return ErrNotLoggedIn
	}

	r.requestID++

	// Command packet
	// https://wiki.vg/RCON#2:_Command
	{
		commandPacket := NewPacket()

		// Length (int32)
		if err := commandPacket.WriteIntLE(int32(10 + len(command))); err != nil {
			return err
		}

		// Request ID (int32)
		if err := commandPacket.WriteIntLE(int32(r.requestID)); err != nil {
			return err
		}

		// Type (int32) - 2
		if err := commandPacket.WriteIntLE(2); err != nil {
			return err
		}

		// Payload (null-terminated string)
		if err := commandPacket.WriteBytes(append([]byte(command), 0x00)); err != nil {
			return err
		}

		// Padding (null byte) - 0x00
		if err := commandPacket.WriteByte(0x00); err != nil {
			return err
		}

		n, err := commandPacket.WriteTo(*r.Conn)

		if err != nil {
			return err
		}

		if n != int64(14+len(command)) {
			return ErrUnexpectedResponse
		}
	}

	return nil
}

func (r *RCON) Close() error {
	r.authSuccess = false
	r.requestID = 0

	if r.Conn != nil {
		if err := (*r.Conn).Close(); err != nil {
			return err
		}
	}

	r.Conn = nil

	return nil
}

func (r *RCON) readMessage() error {
	// Command response packet
	// https://wiki.vg/RCON#0:_Command_response
	{
		var packetLength int32

		// Length - int32
		{
			data := make([]byte, 4)

			n, err := (*r.r).Read(data)

			if err != nil {
				return err
			}

			if n < 4 {
				return nil
			}

			packetLength = int32(binary.LittleEndian.Uint16(data))
		}

		// Request ID - int32
		{
			data := make([]byte, 4)

			n, err := (*r.r).Read(data)

			if err != nil {
				return err
			}

			if n < 4 {
				return nil
			}
		}

		// Type - int32
		{
			data := make([]byte, 4)

			n, err := (*r.r).Read(data)

			if err != nil {
				return err
			}

			if n < 4 {
				return nil
			}

			if binary.LittleEndian.Uint32(data) != 0 {
				return ErrUnexpectedResponse
			}
		}

		// Payload - null-terminated string
		{
			data := make([]byte, packetLength-8)

			n, err := (*r.r).Read(data)

			if err != nil {
				return err
			}

			if n < 1 {
				return nil
			}

			r.Messages <- string(data)
		}
	}

	return nil
}

func parseRCONOptions(opts ...RCONOptions) RCONOptions {
	if len(opts) < 1 {
		return defaultRCONOptions
	}

	return opts[0]
}
