package main

import (
	"context"
	"strings"

	"github.com/Excoriate/daggerverse/module-template/tests/internal/dagger"
)

// TestIACWithTerragruntAlpine tests Terragrunt installation on Alpine.
func (m *Tests) TestIACWithTerragruntAlpine(ctx context.Context) error {
	targetModule := dag.ModuleTemplate(dagger.ModuleTemplateOpts{
		Ctr: dag.
			Container().
			From("alpine:latest").
			WithExec([]string{"apk", "add", "curl"}),
	})

	return m.testTerragruntVersions(ctx, targetModule, true)
}

// TestIACWithTerragruntUbuntu tests Terragrunt installation on Ubuntu.
func (m *Tests) TestIACWithTerragruntUbuntu(ctx context.Context) error {
	targetModule := dag.ModuleTemplate(dagger.ModuleTemplateOpts{
		Ctr: dag.
			Container().
			From("ubuntu:latest").
			WithExec([]string{"apt-get", "update"}).
			WithExec([]string{"apt-get", "install", "-y", "curl", "wget", "unzip"}),
	})

	return m.testTerragruntVersions(ctx, targetModule, false)
}

func (m *Tests) testTerragruntVersions(ctx context.Context, targetModule *dagger.ModuleTemplate, isAlpine bool) error {
	versions := map[string]string{
		"0.53.7": "1.8.0",
		"0.66.5": "1.7.0",
		"0.63.5": "1.6.0",
		"0.66.8": "1.9.4",
	}

	for terragruntVersion, tfVersion := range versions {
		if err := m.verifyModule(ctx, targetModule, terragruntVersion, tfVersion, isAlpine); err != nil {
			return err
		}
	}

	return nil
}

func (m *Tests) verifyModule(
	ctx context.Context,
	targetModule *dagger.ModuleTemplate,
	terragruntVersion, tfVersion string,
	isAlpine bool,
) error {
	var opts interface{}
	if isAlpine {
		opts = dagger.ModuleTemplateWithTerragruntAlpineOpts{
			Version:   terragruntVersion,
			TfVersion: tfVersion,
		}

		alpineOpts, ok := opts.(dagger.ModuleTemplateWithTerragruntAlpineOpts)

		if !ok {
			return Errorf("failed to assert type for Alpine options")
		}

		targetModule = targetModule.WithTerragruntAlpine(alpineOpts)
	} else {
		opts = dagger.ModuleTemplateWithTerragruntUbuntuOpts{
			Version:   terragruntVersion,
			TfVersion: tfVersion,
		}

		ubuntuOpts, ok := opts.(dagger.ModuleTemplateWithTerragruntUbuntuOpts)

		if !ok {
			return Errorf("failed to assert type for Ubuntu options")
		}

		targetModule = targetModule.WithTerragruntUbuntu(ubuntuOpts)
	}

	tests := []struct {
		command []string
		output  string
		check   string
	}{
		{[]string{"terraform", "version"}, tfVersion, "Terraform version"},
		{[]string{"which", "terraform"}, "/usr/local/bin/terraform", "Terraform path"},
		{[]string{"terragrunt", "--version"}, terragruntVersion, "Terragrunt version"},
		{[]string{"which", "terragrunt"}, "/usr/local/bin/terragrunt", "Terragrunt path"},
	}

	for _, test := range tests {
		out, err := targetModule.Ctr().WithExec(test.command).Stdout(ctx)

		if err != nil {
			return WrapErrorf(err, "failed to get %s, the output was: %s", test.check, out)
		}

		if !strings.Contains(out, test.output) {
			return Errorf("expected %s to contain %s, got %s", test.check, test.output, out)
		}
	}

	return nil
}
