package mcstatus

import (
	"fmt"
	"reflect"
	"strings"
)

var (
	formattingColorCodeLookupTable = map[rune]string{
		'0': "black",
		'1': "dark_blue",
		'2': "dark_green",
		'3': "dark_aqua",
		'4': "dark_red",
		'5': "dark_purple",
		'6': "gold",
		'7': "gray",
		'8': "dark_gray",
		'9': "blue",
		'a': "green",
		'b': "aqua",
		'c': "red",
		'd': "light_purple",
		'e': "yellow",
		'f': "white",
		'g': "minecoin_gold",
	}
	colorNameLookupTable = map[string]rune{
		"black":         '0',
		"dark_blue":     '1',
		"dark_green":    '2',
		"dark_aqua":     '3',
		"dark_red":      '4',
		"dark_purple":   '5',
		"gold":          '6',
		"gray":          '7',
		"dark_gray":     '8',
		"blue":          '9',
		"green":         'a',
		"aqua":          'b',
		"red":           'c',
		"light_purple":  'd',
		"yellow":        'e',
		"white":         'f',
		"minecoin_gold": 'g',
	}
	htmlColorLookupTable = map[string]string{
		"black":         "#000000",
		"dark_blue":     "#0000aa",
		"dark_green":    "#00aa00",
		"dark_aqua":     "#00aaaa",
		"dark_red":      "#aa0000",
		"dark_purple":   "#aa00aa",
		"gold":          "#ffaa00",
		"gray":          "#aaaaaa",
		"dark_gray":     "#555555",
		"blue":          "#5555ff",
		"green":         "#55ff55",
		"aqua":          "#55ffff",
		"red":           "#ff5555",
		"light_purple":  "#ff55ff",
		"yellow":        "#ffff55",
		"white":         "#ffffff",
		"minecoin_gold": "#ddd605",
	}
)

// FormatItem is a formatting item parsed from the MOTD for easy use
type FormatItem struct {
	Text          string `json:"text"`
	Color         string `json:"color"`
	Obfuscated    bool   `json:"obfuscated"`
	Bold          bool   `json:"bold"`
	Strikethrough bool   `json:"strikethrough"`
	Underline     bool   `json:"underline"`
	Italic        bool   `json:"italic"`
}

// MOTD contains helper functions for reading and writing the MOTD from a server
type MOTD struct {
	Tree []FormatItem `json:"-"`
}

func parseMOTD(desc interface{}) (*MOTD, error) {
	if v, ok := desc.(string); ok {
		tree, err := parseString(v)

		if err != nil {
			return nil, err
		}

		return &MOTD{
			Tree: tree,
		}, nil
	}

	if m, ok := desc.(map[string]interface{}); ok {
		str := parseChatObject(m)

		tree, err := parseString(str)

		if err != nil {
			return nil, err
		}

		return &MOTD{
			Tree: tree,
		}, nil
	}

	return nil, fmt.Errorf("unknown description type: %s", reflect.TypeOf(desc))
}

// String returns the description with formatting
func (m MOTD) String() string {
	result := ""

	for _, v := range m.Tree {
		if v.Color != "white" {
			colorCode, ok := colorNameLookupTable[v.Color]

			if ok {
				result += "\u00A7" + string(colorCode)
			}
		}

		if v.Obfuscated {
			result += "\u00A7k"
		}

		if v.Bold {
			result += "\u00A7l"
		}

		if v.Strikethrough {
			result += "\u00A7m"
		}

		if v.Underline {
			result += "\u00A7n"
		}

		if v.Italic {
			result += "\u00A7o"
		}

		result += v.Text
	}

	return result
}

// Raw returns the raw description with formatting
func (m MOTD) Raw() string {
	return m.String()
}

// Clean returns the description with no formatting
func (m MOTD) Clean() string {
	result := ""

	for _, v := range m.Tree {
		result += v.Text
	}

	return result
}

// HTML returns the description with HTML formatting
func (m MOTD) HTML() string {
	result := "<span>"

	for _, v := range m.Tree {
		classes := make([]string, 0)
		styles := make(map[string]string)

		color, ok := htmlColorLookupTable[v.Color]

		if ok {
			styles["color"] = color
		}

		if v.Obfuscated {
			classes = append(classes, "minecraft-format-obfuscated")
		}

		if v.Bold {
			styles["font-weight"] = "bold"
		}

		if v.Strikethrough {
			if _, ok = styles["text-decoration"]; ok {
				styles["text-decoration"] += " "
			}

			styles["text-decoration"] += "line-through"
		}

		if v.Underline {
			if _, ok = styles["text-decoration"]; ok {
				styles["text-decoration"] += " "
			}

			styles["text-decoration"] += "underline"
		}

		if v.Italic {
			styles["font-style"] = "italic"
		}

		result += "<span"

		if len(classes) > 0 {
			result += " class=\""

			for _, v := range classes {
				result += v
			}

			result += "\""
		}

		if len(styles) > 0 {
			result += " style=\""

			keys := make([]string, 0, len(styles))

			for k := range styles {
				keys = append(keys, k)
			}

			for i, l := 0, len(keys); i < l; i++ {
				key := keys[i]
				value := styles[key]

				result += fmt.Sprintf("%s: %s;", key, value)

				if i+1 != l {
					result += " "
				}
			}

			result += "\""
		}

		result += fmt.Sprintf(">%s</span>", v.Text)
	}

	return result + "</span>"
}

func parseChatObject(m map[string]interface{}) string {
	result := ""

	color, ok := m["color"].(string)

	if ok {
		code, ok := colorNameLookupTable[color]

		if ok {
			result += "\u00A7" + string(code)
		}
	}

	bold, ok := m["bold"].(string)

	if ok && bold == "true" {
		result += "\u00A7l"
	}

	italic, ok := m["italic"].(string)

	if ok && italic == "true" {
		result += "\u00A7o"
	}

	underline, ok := m["underlined"].(string)

	if ok && underline == "true" {
		result += "\u00A7n"
	}

	strikethrough, ok := m["strikethrough"].(string)

	if ok && strikethrough == "true" {
		result += "\u00A7m"
	}

	obfuscated, ok := m["obfuscated"].(string)

	if ok && obfuscated == "true" {
		result += "\u00A7k"
	}

	text, ok := m["text"].(string)

	if ok {
		result += text
	}

	extra, ok := m["extra"].([]map[string]interface{})

	if ok {
		for _, v := range extra {
			result += parseChatObject(v)
		}
	}

	return result
}

func parseString(s string) ([]FormatItem, error) {
	tree := make([]FormatItem, 0)

	item := FormatItem{
		Text:  "",
		Color: "white",
	}

	r := strings.NewReader(s)

	for r.Len() > 0 {
		char, n, err := r.ReadRune()

		if err != nil {
			return nil, err
		}

		if n < 1 {
			break
		}

		if char != '\u00A7' {
			item.Text += string(char)

			continue
		}

		code, n, err := r.ReadRune()

		if err != nil {
			return nil, err
		}

		if n < 1 {
			break
		}

		// Color code
		{
			name, ok := formattingColorCodeLookupTable[code]

			if ok {
				if item.Obfuscated || item.Bold || item.Strikethrough || item.Underline || item.Italic || name != item.Color {
					if len(item.Text) > 0 {
						tree = append(tree, item)
					}

					item = FormatItem{
						Text:  "",
						Color: name,
					}
				} else {
					item.Color = name
				}

				continue
			}
		}

		// Formatting code
		{
			switch code {
			case 'k':
				{
					if len(item.Text) > 0 {
						tree = append(tree, item)
					}

					item.Text = ""
					item.Obfuscated = true
				}
			case 'l':
				{
					if len(item.Text) > 0 {
						tree = append(tree, item)
					}

					item.Text = ""
					item.Bold = true
				}
			case 'm':
				{
					if len(item.Text) > 0 {
						tree = append(tree, item)
					}

					item.Text = ""
					item.Strikethrough = true
				}
			case 'n':
				{
					if len(item.Text) > 0 {
						tree = append(tree, item)
					}

					item.Text = ""
					item.Underline = true
				}
			case 'o':
				{
					if len(item.Text) > 0 {
						tree = append(tree, item)
					}

					item.Text = ""
					item.Italic = true
				}
			case 'r':
				{
					if len(item.Text) > 0 {
						tree = append(tree, item)
					}

					item = FormatItem{
						Text:  "",
						Color: "white",
					}
				}
			}
		}
	}

	tree = append(tree, item)

	return tree, nil
}
