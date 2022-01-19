package mcstatus

import (
	"bytes"
	"io"
)

func decodeASCII(input []byte) string {
	data := make([]rune, len(input))

	for i, b := range input {
		data[i] = rune(b)
	}

	return string(data)
}

func writePacket(data *bytes.Buffer, w io.Writer) error {
	if _, err := writeVarInt(int32(data.Len()), w); err != nil {
		return err
	}

	_, err := io.Copy(w, data)

	return err
}
