package hw09structvalidator

import "strings"

type TagData struct {
	Name        string
	Args        []string
	OriginalTag string
}

// ParseTag parses given tag to TagData slice.
// If tag is empty it returns nil.
func ParseTag(tag string) []TagData {
	if tag == "" {
		return nil
	}
	validators := strings.Split(tag, "|")
	var data []TagData
	for _, validator := range validators {
		parts := strings.SplitN(validator, ":", 2)
		var args []string
		if len(parts) > 1 {
			args = strings.Split(parts[1], ",")
		}
		data = append(data, TagData{Name: parts[0], Args: args, OriginalTag: validator})
	}
	return data
}
