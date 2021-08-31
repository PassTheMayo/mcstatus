package mcstatus

import (
	"errors"
	"io"
)

var (
	ErrVarIntTooBig = errors.New("size of VarInt exceeds maximum data size")
	ErrVarIntNoData = errors.New("failed to read any bytes while reading varint")
)

func readVarInt(r io.Reader) (int32, int, error) {
	var numRead int = 0
	var result int32 = 0

	for {
		data := make([]byte, 1)

		n, err := r.Read(data)

		if err != nil {
			return 0, numRead, err
		}

		if n < 1 {
			return 0, numRead, ErrVarIntNoData
		}

		value := (data[0] & 0b01111111)
		result |= int32(value) << (7 * numRead)

		numRead++

		if numRead > 5 {
			return 0, numRead, ErrVarIntTooBig
		}

		if (data[0] & 0b10000000) == 0 {
			break
		}
	}

	return result, numRead, nil
}

func readVarLong(r io.Reader) (int64, int, error) {
	var numRead int = 0
	var result int64 = 0

	for {
		data := make([]byte, 1)

		n, err := io.ReadAtLeast(r, data, 1)

		if err != nil {
			return 0, numRead, err
		}

		if n < 1 {
			return 0, numRead, ErrVarIntNoData
		}

		value := (data[0] & 0b01111111)
		result |= int64(value) << (7 * numRead)

		numRead++

		if numRead > 10 {
			return 0, numRead, ErrVarIntTooBig
		}

		if (data[0] & 0b10000000) == 0 {
			break
		}
	}

	return result, numRead, nil
}

func writeVarInt(val int32, w io.Writer) (int, error) {
	var numWritten int = 0

	for {
		if (uint32(val) & 0xFFFFFF80) == 0 {
			_, err := w.Write([]byte{byte(val)})

			numWritten++

			return numWritten, err
		}

		_, err := w.Write([]byte{byte(val&0x7F | 0x80)})

		if err != nil {
			return numWritten, err
		}

		val = int32(uint32(val) >> 7)
	}
}

func writeVarLong(val int64, w io.Writer) (int, error) {
	var numWritten int = 0

	for {
		if (uint64(val) & 0xFFFFFFFFFFFFFF80) == 0 {
			_, err := w.Write([]byte{byte(val)})

			numWritten++

			return numWritten, err
		}

		_, err := w.Write([]byte{byte(val&0x7F | 0x80)})

		if err != nil {
			return numWritten, err
		}

		val = int64(uint64(val) >> 7)
	}
}
