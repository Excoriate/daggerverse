package main

import (
	"fmt"
	"strings"
)

func getContainerImage(image, version string) string {
	if version == "" {
		version = "latest"
	}

	return fmt.Sprintf("%s:%s", image, version)
}

// NewOptional creates a new Optional of any type.
func toDaggerOptional[T any](value T) Optional[T] {
	return Optional[T]{value: value, isSet: true}
}

func convertSliceToMap(input []string) (map[string]string, error) {
	result := make(map[string]string)

	for _, item := range input {
		pairs := strings.Split(item, ",")
		for _, pair := range pairs {
			kv := strings.SplitN(pair, "=", 2)
			if len(kv) != 2 {
				return nil, fmt.Errorf("invalid format for pair: %s", pair)
			}
			key := strings.TrimSpace(unescape(kv[0]))
			value := strings.TrimSpace(unescape(kv[1]))
			result[key] = value
		}
	}

	return result, nil
}

func unescape(input string) string {
	result := strings.ReplaceAll(input, `\\`, `\`)
	result = strings.ReplaceAll(result, `\,`, `,`)
	result = strings.ReplaceAll(result, `\=`, `=`)
	return strings.TrimSpace(result)
}
