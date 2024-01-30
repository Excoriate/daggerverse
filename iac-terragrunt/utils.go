package main

import "fmt"

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
