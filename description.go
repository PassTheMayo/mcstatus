package mcstatus

import "regexp"

var (
	stripFormattingRegExp = regexp.MustCompile("\u00A7[0-9a-f]")
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
