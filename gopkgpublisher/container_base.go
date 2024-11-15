package main

import (
	"path/filepath"

	"github.com/Excoriate/daggerx/pkg/apkox"
	"github.com/Excoriate/daggerx/pkg/fixtures"
)

// ApkoKeyRingInfo represents the keyring information for Apko.
type ApkoKeyRingInfo apkox.KeyringInfo

const (
	defaultApkoImage   = "cgr.dev/chainguard/apko"
	defaultApkoTarball = "image.tar"
)

// BaseApko sets the base image to an Apko image and creates the base container.
//
// Returns a pointer to the Gopkgpublisher instance.
func (m *Gopkgpublisher) BaseApko() (*Gopkgpublisher, error) {
	apkoCfgFilePath := "config/presets/base-alpine.yaml"
	apkoCfgFile := dag.CurrentModule().
		Source().
		File(apkoCfgFilePath)

	apkoCfgFilePathMounted := filepath.Join(fixtures.MntPrefix, apkoCfgFilePath)

	apkoCtr := dag.Container().
		From(defaultApkoImage).
		WithMountedFile(apkoCfgFilePathMounted, apkoCfgFile)

	apkoBuildCmd := []string{
		"apko",
		"build",
		apkoCfgFilePathMounted,
		"latest",
		defaultApkoTarball,
		"--cache-dir",
		"/var/cache/apko",
	}

	for _, pkg := range m.ApkoPackages {
		apkoBuildCmd = append(apkoBuildCmd, "--package-append", pkg)
	}

	apkoCtr = apkoCtr.
		WithExec(apkoBuildCmd)

	outputTar := apkoCtr.
		File(defaultApkoTarball)

	m.Ctr = dag.
		Container().
		Import(outputTar)

	return m, nil
}
