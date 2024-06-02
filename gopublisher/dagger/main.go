package main

import (
	"fmt"

	"github.com/Excoriate/daggerx/pkg/containerx"
)

type Gopublisher struct {
	// Src is the directory that contains all the source code, including the module directory.
	// It represents the source code directory.
	Src *Directory
	// Ctr is the container to use as a base container for gopublisher, if it's passed, it's used as the base container.
	Ctr *Container
}

const (
	defaultGoVersion = "1.22.3-alpine3.19"
	pkgRegistryURL   = "https://pkg.go.dev"
	goProxyURL       = "proxy.golang.org"
)

func New(
	// version is the version of Go that the publisher module will use, e.g., "1.22.0".
	// +optional
	version string,
) (*Gopublisher, error) {
	m := &Gopublisher{}

	if version == "" {
		version = defaultGoVersion
	} else {
		// by design, this module is going to use the alpine3.19 version of Go by default.
		version = fmt.Sprintf("%s-alpine3.19", version)
	}

	imageURL, err := containerx.GetImageURL(&containerx.NewBaseContainerOpts{
		Image:   "golang",
		Version: version,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get image URL: %w", err)
	}

	m.Ctr = m.Base(imageURL).Ctr
	return m, nil
}

// Base sets the base container for gopublisher.
//
// It sets the base container for the gopublisher module. It's meant to be used as a base container for the module.
// Arguments:
// - imageURL: The URL of the image to use as the base container.
// Returns:
// - *Gopublisher: The Gopublisher module with the base container set.
func (m *Gopublisher) Base(imageURL string) *Gopublisher {
	c := dag.Container().From(imageURL)

	m.Ctr = c

	return m.WithCGODisabled().
		WithGit().
		WithCURL()
}
