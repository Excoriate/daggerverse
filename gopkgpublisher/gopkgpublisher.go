package main

import (
	"context"

	"github.com/Excoriate/daggerverse/gopkgpublisher/internal/dagger"
)

// Publish represents the configuration for publishing a package using GopkgPublisher.
type Publish struct {
	// GopkgPublisher is the Gopkgpublisher module that handles the publishing process.
	// +private
	GopkgPublisher *Gopkgpublisher
	// PublishOpts contains the options for the Publish command, allowing customization of the publishing behavior.
	// +private
	PublishOpts *PublishOpts
}

// PublishOpts holds the options for the Publish command.
type PublishOpts struct {
	// src is the source directory to publish.
	// +private
	Src *dagger.Directory
	// SkipTests indicates whether to skip running tests during the publish process.
	// +private
	SkipTests bool
	// DryRun indicates whether to perform a dry run without making any actual changes.
	// +private
	DryRun bool
	// GitPushTags indicates whether to push the tags to the remote repository after publishing.
	// +private
	GitPushTags bool
	// PackageVersion is the version of the package to publish.
	// +private
	PackageVersion string
}

// PublishOption is a function type that defines a configuration option for the PublishOpts structure.
// It takes a pointer to PublishOpts and returns an error if the option cannot be applied.
// This allows for flexible and composable configuration of publishing options when using the GopkgPublisher.
type PublishOptionFn func(*PublishOpts) error

// WithSrc is a method of the Gopkgpublisher type that returns a PublishOption function.
// This function allows the user to specify the source directory for the publishing process.
//
// The source directory is represented by a pointer to a dagger.Directory instance, which
// is assigned to the Src field of the PublishOpts structure. This enables the GopkgPublisher
// to know where to find the package that needs to be published.
//
// Parameters:
//   - src: A pointer to a dagger.Directory that represents the source directory to be used
//     during the publishing process. This parameter must not be nil; otherwise, the behavior
//     of the publishing process may be undefined.
//
// Returns:
//   - A PublishOption function that takes a pointer to PublishOpts and returns an error.
//     The returned function sets the Src field of the PublishOpts to the provided src value.
//
// Example usage:
//
//	publishOption := gopkgPublisher.WithSrc(mySourceDirectory)
//	err := publishOption(&myPublishOpts)
func (p *Publish) WithSrc(src *dagger.Directory) PublishOptionFn {
	return func(opts *PublishOpts) error {
		if opts == nil {
			return NewError("publish options cannot be nil")
		}

		// check if src has a file called go.mod
		fileID, err := src.File("go.mod").ID(context.Background())
		if err != nil {
			return WrapError(err, "failed to check for go.mod file")
		}

		if fileID == "" {
			return NewError("go.mod file not found in source directory")
		}

		opts.Src = src

		return nil
	}
}

func (p *Publish) CompileOpts(opts ...PublishOptionFn) error {
	for _, opt := range opts {
		if err := opt(p.PublishOpts); err != nil {
			return WrapError(err, "failed to compile publish options")
		}
	}

	return nil
}

// validateGoVersion validates and normalizes a Go version string.
// It ensures the version follows Go's versioning format (v1.X or v1.X.Y).
// Returns the normalized version string and any validation error.
// func validateGoVersion(version string) (string, error) {
// 	// Remove 'v' prefix if present for consistent handling
// 	version = strings.TrimPrefix(version, "v")

// 	// Basic semver validation
// 	parts := strings.Split(version, ".")
// 	if len(parts) < 2 || len(parts) > 3 {
// 		return "", Errorf("invalid version format: v%s. Must be in format v1.X or v1.X.Y", version)
// 	}

// 	// Validate major version is 1 for Go
// 	if parts[0] != "1" {
// 		return "", Errorf("invalid Go major version: %s. Must be 1", parts[0])
// 	}

// 	// Validate each part is a valid number
// 	for i, part := range parts {
// 		num, err := strconv.Atoi(part)
// 		if err != nil {
// 			return "", Errorf("invalid version number in %s: %s is not a number", version, part)
// 		}

// 		// Additional validation for each part
// 		if i == 0 && num != 1 {
// 			return "", Errorf("Go major version must be 1, got: %d", num)
// 		}
// 		if num < 0 {
// 			return "", Errorf("version numbers cannot be negative: %d", num)
// 		}
// 	}

// 	// Normalize to semver format with v prefix
// 	if len(parts) == 2 {
// 		return "v" + version + ".0", nil
// 	}
// 	return "v" + version, nil
// }

// func (m *Gopkgpublisher) isValid(opts *PublishOpts) error {
// 	if opts == nil {
// 		return NewError("publish options cannot be nil")
// 	}

// 	// Validate source directory
// 	if opts.Src == nil {
// 		return NewError("source directory is required")
// 	}

// 	// Validate Go version format if specified
// 	if opts.GoVersion != "" {
// 		normalizedVersion, err := validateGoVersion(opts.GoVersion)
// 		if err != nil {
// 			return err
// 		}
// 		opts.GoVersion = normalizedVersion
// 	}

// 	// Validate mutually exclusive options
// 	if opts.DryRun && opts.GitPushTags {
// 		return NewError("cannot use --git-push-tags with --dry-run")
// 	}

// 	// Validate container configuration
// 	if m.Ctr == nil {
// 		return NewError("container configuration is not initialized")
// 	}

// 	// If git push tags is enabled, verify git repository is valid
// 	if opts.GitPushTags {
// 		// Check if directory is a git repository
// 		_, err := git.PlainOpen(opts.Src.Path())
// 		if err != nil {
// 			if err == git.ErrRepositoryNotExists {
// 				return NewError("git repository not found when --git-push-tags is enabled")
// 			}
// 			return WrapError(err, "failed to open git repository")
// 		}
// 	}

// 	// Validate go.mod exists in source directory
// 	exists, err := opts.Src.File("go.mod").ID(context.Background())
// 	if err != nil {
// 		return WrapError(err, "failed to check for go.mod file")
// 	}
// 	if exists == "" {
// 		return NewError("go.mod file not found in source directory")
// 	}

// 	return nil
// }
