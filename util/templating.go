package util

import (
	"fmt"
	"strings"
)

// GenerateFilename creates a filename from a template and a data map using simple substitution.
func GenerateFilename(template string, data map[string]interface{}) (string, error) {
	var replacements []string

	for key, value := range data {
		placeholder := "{" + key + "}"
		// Convert value to string for the replacer
		valStr := fmt.Sprintf("%v", value)
		replacements = append(replacements, placeholder, valStr)
	}

	replacer := strings.NewReplacer(replacements...)
	result := replacer.Replace(template)

	return result, nil
}
