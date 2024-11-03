package main

import (
	"strings"
)

// GoBuildOptions holds the options for building Go packages.
//
// This struct contains flags that can be passed to the Go build command.
type GoBuildOptions struct {
	// Flags are the flags to pass to the Go build command.
	// +private
	Flags []string
}

// NewGoBuildOptions creates a new GoBuildOptions instance.
//
// This function initializes a new GoBuildOptions struct with default values.
// It returns a pointer to the newly created struct.
func NewGoBuildOptions() *GoBuildOptions {
	return &GoBuildOptions{
		Flags: []string{},
	}
}

// WithRace enables the race detector.
//
// This function adds the -race flag to the build options.
// Equivalent to: go build -race.
func (o *GoBuildOptions) WithRace() *GoBuildOptions {
	o.Flags = append(o.Flags, "-race")

	return o
}

// WithMSan enables interoperation with memory sanitizer.
//
// This function adds the -msan flag to the build options.
// Equivalent to: go build -msan.
func (o *GoBuildOptions) WithMSan() *GoBuildOptions {
	o.Flags = append(o.Flags, "-msan")

	return o
}

// WithASan enables interoperation with address sanitizer.
//
// This function adds the -asan flag to the build options.
// Equivalent to: go build -asan.
func (o *GoBuildOptions) WithASan() *GoBuildOptions {
	o.Flags = append(o.Flags, "-asan")

	return o
}

// WithTags sets build tags.
//
// This function adds the -tags flag with the specified tags to the build options.
// Multiple tags should be comma-separated.
// Equivalent to: go build -tags=<tags>.
func (o *GoBuildOptions) WithTags(tags string) *GoBuildOptions {
	tags = strings.TrimSpace(tags)
	if tags != "" {
		o.Flags = append(o.Flags, "-tags", tags)
	}

	return o
}

// WithLDFlags sets linker flags.
//
// This function adds the -ldflags flag with the specified flags to the build options.
func (o *GoBuildOptions) WithLDFlags(flags string) *GoBuildOptions {
	flags = strings.TrimSpace(flags)

	if flags != "" {
		o.Flags = append(o.Flags, "-ldflags", flags)
	}

	return o
}

// WithGCFlags sets garbage collection flags.
//
// This function adds the -gcflags flag with the specified flags to the build options.
func (o *GoBuildOptions) WithGCFlags(flags string) *GoBuildOptions {
	flags = strings.TrimSpace(flags)

	if flags != "" {
		o.Flags = append(o.Flags, "-gcflags", flags)
	}

	return o
}

// WithAsmFlags sets assembler flags.
//
// This function adds the -asmflags flag with the specified flags to the build options.
// Equivalent to: go build -asmflags=<flags>.
func (o *GoBuildOptions) WithAsmFlags(flags string) *GoBuildOptions {
	flags = strings.TrimSpace(flags)
	if flags != "" {
		o.Flags = append(o.Flags, "-asmflags", flags)
	}

	return o
}

// WithTrimPath removes file system paths from resulting executable.
//
// This function adds the -trimpath flag to the build options.
// Equivalent to: go build -trimpath.
func (o *GoBuildOptions) WithTrimPath() *GoBuildOptions {
	o.Flags = append(o.Flags, "-trimpath")

	return o
}

// WithWork prints the name of the temporary work directory and does not delete it when exiting.
//
// This function adds the -work flag to the build options.
// Equivalent to: go build -work.
func (o *GoBuildOptions) WithWork() *GoBuildOptions {
	o.Flags = append(o.Flags, "-work")

	return o
}

// WithBuildMode sets the build mode.
//
// This function adds the -buildmode flag with the specified mode to the build options.
// Valid modes are: archive, c-archive, c-shared, default, shared, exe, pie, plugin.
// Equivalent to: go build -buildmode=<mode>.
func (o *GoBuildOptions) WithBuildMode(mode string) *GoBuildOptions {
	mode = strings.TrimSpace(mode)
	validModes := map[string]bool{
		"archive": true, "c-archive": true, "c-shared": true,
		"default": true, "shared": true, "exe": true,
		"pie": true, "plugin": true,
	}

	if validModes[mode] {
		o.Flags = append(o.Flags, "-buildmode", mode)
	}

	return o
}

// WithCompiler sets the compiler to use.
//
// This function adds the -compiler flag with the specified compiler to the build options.
// Valid compilers are: gc, gccgo.
// Equivalent to: go build -compiler=<compiler>.
func (o *GoBuildOptions) WithCompiler(compiler string) *GoBuildOptions {
	compiler = strings.TrimSpace(compiler)
	if compiler == "gc" || compiler == "gccgo" {
		o.Flags = append(o.Flags, "-compiler", compiler)
	}

	return o
}

// WithGCCGOFlags sets gccgo flags.
//
// This function adds the -gccgoflags flag with the specified flags to the build options.
// Equivalent to: go build -gccgoflags=<flags>.
func (o *GoBuildOptions) WithGCCGOFlags(flags string) *GoBuildOptions {
	flags = strings.TrimSpace(flags)
	if flags != "" {
		o.Flags = append(o.Flags, "-gccgoflags", flags)
	}

	return o
}

// WithMod sets module download mode.
//
// This function adds the -mod flag with the specified mode to the build options.
// Valid modes are: readonly, vendor, mod.
// Equivalent to: go build -mod=<mode>.
func (o *GoBuildOptions) WithMod(mode string) *GoBuildOptions {
	mode = strings.TrimSpace(mode)
	validModes := map[string]bool{"readonly": true, "vendor": true, "mod": true}

	if validModes[mode] {
		o.Flags = append(o.Flags, "-mod", mode)
	}

	return o
}

// Validate checks if the build options are valid.
//
// This function performs validation on the build options to ensure they are correctly formed.
// It returns an error if any issues are found, otherwise it returns nil.
func (o *GoBuildOptions) Validate() error {
	for i := 0; i < len(o.Flags); i++ {
		flag := o.Flags[i]
		if strings.HasPrefix(flag, "-") && len(flag) > 1 {
			if i+1 < len(o.Flags) && !strings.HasPrefix(o.Flags[i+1], "-") {
				// Skip the next item as it's a value for this flag
				i++
			}
		} else {
			return Errorf("invalid build flag: %s", flag)
		}
	}

	return nil
}
