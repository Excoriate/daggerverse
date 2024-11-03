package main

import "github.com/Excoriate/daggerx/pkg/containerx"

const (
	// defaultContainerVersion is the version of the default container image used in the Gotest module.
	// It is set to "1.23.0-alpine3.20" for compatibility with Alpine 3.20.
	defaultContainerVersion = "1.23.0-alpine3.20"

	// defaultContainerImage is the name of the default container image used in the Gotest module.
	// It is set to "golang" to provide a Go environment for the module.
	defaultContainerImage = "golang"
)

// Base sets the base image for the Gotest module and creates the base container.
//
// Parameters:
// - imageURL: The URL of the image to use as the base container.
//
// Returns a pointer to the updated Gotest instance.
func (m *Gotest) Base(imageURL string) *Gotest {
	c := dag.Container().
		From(imageURL)

	m.Ctr = c

	return m
}

// getImageURL gets the image URL and validates it.
//
// If the image URL is not valid, it returns an error.
func (m *Gotest) getImageURL(image, version string) (string, error) {
	imageURL, err := containerx.GetImageURL(&containerx.NewBaseContainerOpts{
		Image:           image,
		Version:         version,
		FallbackImage:   defaultContainerImage,
		FallBackVersion: defaultContainerVersion,
	})

	if err != nil {
		return "", WrapErrorf(err, "failed to get image URL from image %s and version %s", image, version)
	}

	isValid, invalidImageErr := containerx.ValidateImageURL(imageURL)
	if invalidImageErr != nil {
		return "", WrapErrorf(invalidImageErr, "failed to validate image URL %s", imageURL)
	}

	if !isValid {
		return "", Errorf("the image URL %s is not valid", imageURL)
	}

	return imageURL, nil
}
