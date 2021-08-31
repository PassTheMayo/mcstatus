package mcstatus

import (
	"io"
)

func readString(r io.Reader) (string, error) {
	strLength, _, err := readVarInt(r)

	if err != nil {
		return "", err
	}

	result := make([]byte, 0)

	for i := 0; i < int(strLength); {
		data := make([]byte, 4096)

		n, err := r.Read(data)

		if err != nil {
			return "", err
		}

		result = append(result, data[:n]...)

		i += n
	}

	return string(result), err
}

func writeString(val string, w io.Writer) error {
	_, err := writeVarInt(int32(len(val)), w)

	if err != nil {
		return err
	}

	_, err = w.Write([]byte(val))

	return err
}
