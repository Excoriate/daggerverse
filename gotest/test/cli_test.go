package test

import (
	"context"
	"github.com/testcontainers/testcontainers-go/wait"
	"io"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
)

func TestWithDaggerCLI(t *testing.T) {
	daggerModuleSrc, err := filepath.Abs(filepath.Join(".", "../", "dagger"))
	assert.NoErrorf(t, err, "Could not get dagger module source: %s", err)

	ctx := context.Background()

	ctrReq := testcontainers.ContainerRequest{
		Image: "alpine:latest",
		Cmd:   []string{"tail", "-f", "/dev/null"}, // Add this line

		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      daggerModuleSrc,
				ContainerFilePath: "/dagger",
				FileMode:          0755,
			},
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("")).
			WithDeadline(30 * time.Second),
	}

	daggerC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: ctrReq,
		Started:          true,
	})

	assert.NoErrorf(t, err, "Could not create container: %s", err)

	// Just installing Curl
	c, reader, err := daggerC.Exec(ctx, []string{"sh", "-c", "apk add --no-cache curl"})
	assert.NoErrorf(t, err, "Could not install curl: %s", err)
	assert.Equalf(t, c, 0, "Expected 0, got %d", c)
	buf := new(strings.Builder)
	_, err = io.Copy(buf, reader)
	assert.NoErrorf(t, err, "Could not read from reader: %s", err)
	t.Logf("Output: %s", buf.String())

	// Installing the dagger CLI.
	cDaggerInstall, readerDaggerInstall, daggerInstallErr := daggerC.Exec(ctx, []string{"sh", "-c", "curl -L https://dl.dagger.io/dagger/install.sh | sh"})
	assert.NoErrorf(t, daggerInstallErr, "Could not install dagger: %s", daggerInstallErr)
	assert.Equalf(t, cDaggerInstall, 0, "Expected 0, got %d", cDaggerInstall)
	bufDaggerInstall := new(strings.Builder)
	_, err = io.Copy(bufDaggerInstall, readerDaggerInstall)
	assert.NoErrorf(t, err, "Could not read from reader: %s", err)
	t.Logf("Output: %s", bufDaggerInstall.String())

	cDaggerHelp, readerDaggerHelp, daggerHelpErr := daggerC.Exec(ctx, []string{"sh", "-c", "ls -ltrah dagger"})
	assert.NoErrorf(t, daggerHelpErr, "Could not execute command: %s", daggerHelpErr)
	assert.Equalf(t, cDaggerHelp, 0, "Expected 0, got %d", cDaggerHelp)
	bufDaggerHelp := new(strings.Builder)
	_, err = io.Copy(bufDaggerHelp, readerDaggerHelp)
	assert.NoErrorf(t, err, "Could not read from reader: %s", err)
	t.Logf("Output: %s", bufDaggerHelp.String())

	defer func() {
		err := daggerC.Terminate(ctx)
		assert.NoErrorf(t, err, "Could not terminate container: %s", err)
	}()
}
