package mcstatus

func decodeASCII(input []byte) string {
	data := make([]rune, len(input))

	for i, b := range input {
		data[i] = rune(b)
	}

	return string(data)
}
