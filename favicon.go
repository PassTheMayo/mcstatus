package mcstatus

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"os"
	"strings"
)

// Favicon contains helper functions for reading and writing the favicon
type Favicon struct {
	raw    string
	exists bool
}

func (f Favicon) Exists() bool {
	return f.exists
}

func (f Favicon) String() string {
	if !f.exists {
		return "<nil>"
	}

	return f.raw
}

func (f Favicon) Image() (image.Image, error) {
	data, err := base64.StdEncoding.DecodeString(strings.Replace(f.raw, "data:image/png;base64,", "", 1))

	if err != nil {
		return nil, err
	}

	img, err := png.Decode(bytes.NewBuffer(data))

	if err != nil {
		return nil, err
	}

	return img, nil
}

func (f Favicon) SaveToFile(path string) error {
	img, err := f.Image()

	if err != nil {
		return err
	}

	file, err := os.Create(path)

	if err != nil {
		return err
	}

	defer file.Close()

	return png.Encode(file, img)
}

func parseFavicon(raw interface{}) Favicon {
	if v, ok := raw.(string); ok {
		return Favicon{
			raw:    v,
			exists: true,
		}
	}

	return Favicon{
		raw:    "",
		exists: false,
	}
}
