package container

import "fmt"

func FormImageAddress(image, version string) string {
	if version == "" {
		version = "latest"
	}

	return fmt.Sprintf("%s:%s", image, version)
}
