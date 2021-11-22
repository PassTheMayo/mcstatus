package mcstatus

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

type VoteOptions struct {
	ServiceName string
	Username    string
	Token       string
	UUID        string
	Timestamp   time.Time
	Timeout     time.Duration
}

type voteMessage struct {
	Payload   string `json:"payload"`
	Signature string `json:"signature"`
}

type votePayload struct {
	ServiceName string `json:"serviceName"`
	Username    string `json:"username"`
	Address     string `json:"address"`
	Timestamp   int64  `json:"timestamp"`
	Challenge   string `json:"challenge"`
	UUID        string `json:"uuid,omitempty"`
}

type voteResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

// SendVote sends a Votifier vote to the specified Minecraft server
func SendVote(host string, port uint16, options VoteOptions) error {
	conn, err := net.Dial("tcp4", fmt.Sprintf("%s:%d", host, port))

	if err != nil {
		return err
	}

	defer conn.Close()

	conn.SetDeadline(time.Now().Add(options.Timeout))

	var challenge string

	// Handshake packet
	// https://github.com/NuVotifier/NuVotifier/wiki/Technical-QA#handshake
	{
		data := make([]byte, 0)

		for {
			byteData := make([]byte, 1)

			n, err := conn.Read(byteData)

			if err != nil {
				return err
			}

			if n < 1 {
				return io.EOF
			}

			if byteData[0] == '\n' {
				break
			}

			data = append(data, byteData...)
		}

		split := strings.Split(string(data), " ")

		if split[1] != "2" {
			return ErrUnknownVersion
		}

		challenge = split[2]
	}

	// Vote packet
	// https://github.com/NuVotifier/NuVotifier/wiki/Technical-QA#protocol-v2
	{
		packet := NewPacket()

		payload := votePayload{
			ServiceName: options.ServiceName,
			Username:    options.Username,
			Address:     fmt.Sprintf("%s:%d", host, port),
			Timestamp:   options.Timestamp.UnixNano() / int64(time.Millisecond),
			Challenge:   challenge,
			UUID:        options.UUID,
		}

		payloadData, err := json.Marshal(payload)

		if err != nil {
			return err
		}

		hash := hmac.New(sha256.New, []byte(options.Token))
		hash.Write(payloadData)

		message := voteMessage{
			Payload:   string(payloadData),
			Signature: base64.StdEncoding.EncodeToString(hash.Sum(nil)),
		}

		messageData, err := json.Marshal(message)

		if err != nil {
			return err
		}

		if err = packet.WriteUInt16BE(0x733A); err != nil {
			return err
		}

		if err = packet.WriteUInt16BE(uint16(len(messageData))); err != nil {
			return err
		}

		if err = packet.WriteBytes(messageData); err != nil {
			return err
		}

		if _, err := packet.WriteTo(conn); err != nil {
			return err
		}
	}

	// Response packet
	// https://github.com/NuVotifier/NuVotifier/wiki/Technical-QA#protocol-v2
	{
		data := make([]byte, 0)

		for {
			byteData := make([]byte, 1)

			n, err := conn.Read(byteData)

			if err != nil {
				return err
			}

			if n < 1 {
				return io.EOF
			}

			if byteData[0] == '\n' {
				break
			}

			data = append(data, byteData...)
		}

		response := voteResponse{}

		if err = json.Unmarshal(data, &response); err != nil {
			return err
		}

		switch response.Status {
		case "ok":
			{
				return nil
			}
		case "error":
			{
				return fmt.Errorf("server returned error: %s", response.Error)
			}
		default:
			{
				return ErrUnexpectedResponse
			}
		}
	}
}
