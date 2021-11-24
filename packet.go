package mcstatus

import (
	"bytes"
)

func writePacketLength(buf *bytes.Buffer) (*bytes.Buffer, error) {
	result := &bytes.Buffer{}
	data := buf.Bytes()

	if _, err := writeVarInt(int32(len(data)), result); err != nil {
		return nil, err
	}

	if _, err := result.Write(data); err != nil {
		return nil, err
	}

	return result, nil
}
