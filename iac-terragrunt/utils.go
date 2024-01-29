package main

import "fmt"

func getContainerImage(image, version string) string {
	if version == "" {
		version = "latest"
	}

	return fmt.Sprintf("%s:%s", image, version)
}
