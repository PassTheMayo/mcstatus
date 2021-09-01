package mcstatus

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
)

var (
	stripFormattingRegExp = regexp.MustCompile("\u00A7[0-9a-f]")
	colorLookupTable      = map[string]string{
		"black":        "0",
		"dark_blue":    "1",
		"dark_green":   "2",
		"dark_aqua":    "3",
		"dark_red":     "4",
		"dark_purple":  "5",
		"gold":         "6",
		"gray":         "7",
		"dark_gray":    "8",
		"blue":         "9",
		"green":        "a",
		"aqua":         "b",
		"red":          "c",
		"light_purple": "d",
		"yellow":       "e",
		"white":        "f",
		"0":            "0",
		"1":            "1",
		"2":            "2",
		"3":            "3",
		"4":            "4",
		"5":            "5",
		"6":            "6",
		"7":            "7",
		"8":            "8",
		"9":            "9",
		"a":            "a",
		"b":            "b",
		"c":            "c",
		"d":            "d",
		"e":            "e",
		"f":            "f",
	}
	htmlColorLookupTable = map[rune]string{
		'0': "#000000",
		'1': "#0000aa",
		'2': "#00aa00",
		'3': "#00aaaa",
		'4': "#aa0000",
		'5': "#aa00aa",
		'6': "#ffaa00",
		'7': "#aaaaaa",
		'8': "#555555",
		'9': "#5555ff",
		'a': "#55ff55",
		'b': "#55ffff",
		'c': "#ff5555",
		'd': "#ff55ff",
		'e': "#ffff55",
		'f': "#ffffff",
	}
)

// Description contains helper functions for reading and writing the description
type Description struct {
	raw string
}

// String converts the description into raw text without formatting
func (d Description) String() string {
	return stripFormattingRegExp.ReplaceAllString(d.raw, "")
}

// Raw returns the raw description with formatting
func (d Description) Raw() string {
	return d.raw
}

// HTML returns the description with HTML formatting
func (d Description) HTML() (string, error) {
	result := "<span>"
	buf := bytes.NewBufferString(d.Raw())

	tagsOpen := 1
	bold := false
	italics := false
	underline := false
	strikethrough := false
	color := "r"

	for {
		c, n, err := buf.ReadRune()

		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return "", err
		}

		if n < 1 {
			break
		}

		if c == '\u00A7' {
			charCode, _, err := buf.ReadRune()

			if err != nil {
				return "", err
			}

			if v, ok := htmlColorLookupTable[charCode]; ok {
				if color == v {
					continue
				}

				color = v

				result += fmt.Sprintf("</span><span style=\"color: %s;\">", v)
			}

			if charCode == 'l' && !bold {
				result += "<span style=\"font-weight: bold;\">"
				bold = true
				tagsOpen++
			}

			if charCode == 'm' && !strikethrough {
				result += "<span style=\"text-decoration: line-through;\">"
				strikethrough = true
				tagsOpen++
			}

			if charCode == 'n' && !underline {
				result += "<span style=\"text-decoration: underline;\">"
				underline = true
				tagsOpen++
			}

			if charCode == 'o' && !italics {
				result += "<span style=\"font-style: italic;\">"
				italics = true
				tagsOpen++
			}

			if charCode == 'r' {
				if bold {
					bold = false

					result += "</span>"
				}

				if strikethrough {
					strikethrough = false

					result += "</span>"
				}

				if underline {
					underline = false

					result += "</span>"
				}

				if italics {
					italics = false

					result += "</span>"
				}
			}
		} else {
			result += string(c)
		}
	}

	for i := 0; i < tagsOpen; i++ {
		result += "</span>"
	}

	return result, nil
}

func parseChatObject(m map[string]interface{}) string {
	result := ""

	if v, ok := m["bold"].(string); ok && v == "true" {
		result += "\u00A7l"
	}

	if v, ok := m["italic"].(string); ok && v == "true" {
		result += "\u00A7l"
	}

	if v, ok := m["underlined"].(string); ok && v == "true" {
		result += "\u00A7l"
	}

	if v, ok := m["strikethrough"].(string); ok && v == "true" {
		result += "\u00A7l"
	}

	if v, ok := m["obfuscated"].(string); ok && v == "true" {
		result += "\u00A7l"
	}

	if v, ok := m["color"].(string); ok {
		if c, ok := colorLookupTable[v]; ok {
			result += "\u00A7" + c
		}
	}

	if v, ok := m["text"].(string); ok {
		result += v
	}

	if e, ok := m["extra"].([]map[string]interface{}); ok {
		for _, v := range e {
			result += parseChatObject(v)
		}
	}

	return result
}

func parseDescription(raw interface{}) Description {
	if v, ok := raw.(string); ok {
		return Description{v}
	}

	if m, ok := raw.(map[string]interface{}); ok {
		return Description{parseChatObject(m)}
	}

	return Description{}
}
