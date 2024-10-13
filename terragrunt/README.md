# Terragrunt Module for Dagger

A powerful [Dagger](https://dagger.io) module for managing Terragrunt, Terraform, and OpenTofu operations in a containerized environment.

## Features 🎨

| Feature                          | Description                                                                |
| -------------------------------- | -------------------------------------------------------------------------- |
| 🛠️ Flexible Base Image           | Built using APKO for a secure and optimized container environment.         |
| 🔧 Multi-Tool Support            | Primarily focused on Terragrunt, but also supports Terraform and OpenTofu. |
| ⚙️ Customizable Configurations   | Extensive options for Terragrunt and Terraform settings.                   |
| 💾 Caching Mechanisms            | Implements caching for Terragrunt and Terraform for improved performance.  |
| ☁️ AWS CLI Integration           | Option to include AWS CLI in the container.                                |
| 🔐 Permissions Management        | Fine-grained control over directory permissions.                           |
| 🌐 Environment Variable Handling | Easy setting and management of environment variables.                      |
| 🔒 Secret Management             | Secure handling of sensitive information like Terraform tokens.            |
| 🚀 Execution Flexibility         | Run Terragrunt, Terraform, or shell commands within the container.         |

### Terragrunt Batteries Included 🔋

| Feature                                         | Description                                                             |
| ----------------------------------------------- | ----------------------------------------------------------------------- |
| 🛠️ Terragrunt, Terraform, and OpenTofu binaries | Pre-installed binaries for Terragrunt, Terraform, and OpenTofu.         |
| 📄 Terragrunt Configuration                     | Best practice configuration files for Terragrunt.                       |
| ⚙️ Terragrunt Options                           | Configurable options for Terragrunt (see `terragrunt_opts.go`).         |
| 🔧 Directory Permissions                        | Manage directory permissions (see `terragrunt_cfg.go`).                 |
| 💾 Caching Configuration                        | Setup caching for Terragrunt and Terraform (see `terragrunt_cfg.go`).   |
| 🌐 Environment Variables                        | Handle environment variables for Terragrunt (see `terragrunt_opts.go`). |
| 🔐 Secret Management                            | Secure handling of sensitive information like Terraform tokens.         |

## Configuration 🛠️

### Base Container Options

- `ctr`: Specify a custom base container.
- `imageURL`: Specify a custom base image URL.
- `tgVersion`: Set the version of Terragrunt (default: `0.68.1`).
- `tfVersion`: Set the version of Terraform (default: `1.9.5`).
- `openTofuVersion`: Set the version of OpenTofu (default: `1.8.2`).
- `enableAWSCLI`: Enable or disable the installation of the AWS CLI (default: `false`).
- `awscliVersion`: Set the version of the AWS CLI to install (default: `2.15.1`).
- `extraPackages`: A list of extra packages to install with APKO, from the Alpine packages repository (default: `[]`).

### IaC Tool Versions

- Terragrunt, Terraform, and OpenTofu versions can be specified or will use defaults.

### Permissions and Caching

- Configure directory permissions and set up caching for Terragrunt and Terraform.

### Environment and Secrets

- Set environment variables and manage secrets securely.

## Usage Examples 🚀

### Basic Terragrunt Execution

```go
	testEnvVars := []string{
		"AWS_ACCESS_KEY_ID=test",
		"AWS_SECRET_ACCESS_KEY=test",
		"AWS_SESSION_TOKEN=test",
	}

	// Initialize the Terragrunt module
	tgModule := dag.
		Terragrunt(dagger.TerragruntOpts{
			EnvVarsFromHost: testEnvVars,
		}).
		WithTerragruntPermissionsOnDirsDefault().
		WithTerragruntLogOptions(
			dagger.TerragruntWithTerragruntLogOptionsOpts{
				TgLogLevel:        "debug",
				TgForwardTfStdout: true,
			},
		)

	// Execute the init command, but don't run it in a container
	tgCtrConfigured := tgModule.
		Exec("init", dagger.TerragruntExecOpts{
			Source: m.getTestDir("").
				Directory("terragrunt"),
		})

	// Evaluate the terragrunt init command.
	tgInitCmdOut, tgInitCmdErr := tgCtrConfigured.
		Stdout(ctx)
```

### Running Terragrunt with Custom Options

```go
	tgModule := dag.
		Terragrunt(dagger.TerragruntOpts{
			EnvVarsFromHost: testEnvVars,
			TfVersion:       "1.7.0",
		}).
		WithTerragruntPermissionsOnDirsDefault().
		WithTerragruntLogOptions(dagger.TerragruntWithTerragruntLogOptionsOpts{
			TgLogDisableFormatting: true,
			TgLogShowAbsPaths:      true,
			TgLogLevel:             "debug",
		}).
		WithTerraformLogOptions(dagger.TerragruntWithTerraformLogOptionsOpts{
			TfLog:     "debug",
			TfLogPath: "/mnt/tflogs", // it's a directory that the terragrunt user owns.
		}).
		// Extra options added for more realism.
		WithTerragruntOptions(dagger.TerragruntWithTerragruntOptionsOpts{
			IgnoreDependencyErrors:     true,
			IgnoreExternalDependencies: true,
			DisableBucketUpdate:        true,
		})

	// Execute the plan command and get the container back.
	tgCtr := tgModule.Exec("plan", dagger.TerragruntExecOpts{
		Source: m.
			getTestDir("").
			Directory("terragrunt"),
		Secrets: []*dagger.Secret{
			dbPasswordSecret,
			apiKeySecret,
			sshKeySecret,
		},
		// Args to output the plan to a file.
		Args: []string{
			"-out=plan.tfplan",
			"-refresh=true",
		},
	})


```

## Testing 🧪

The module includes comprehensive tests covering various aspects of functionality. You can run these tests using:

```bash
just test terragrunt
```

## Developer Experience 🛠️

To contribute or modify the module:

1. Use [Just](https://just.systems) for task automation.
2. Utilize [Nix](https://nixos.org) for managing the development environment.

Common commands:

```bash
just run-hooks           # Initialize pre-commit hooks
just lintall terragrunt  # Run linter
just test terragrunt    # Run tests
just ci terragrunt # Run entire CI tasks locally
```

## APKO Base Image

This module uses [APKO](https://github.com/chainguard-dev/apko) to build its base image, ensuring:

- Enhanced security through minimal attack surface
- Optimized container size and performance
- Reproducible and declarative image builds

For more information on APKO, refer to the [Chainguard APKO documentation](https://github.com/chainguard-dev/apko/tree/main/docs).

---

For detailed API documentation and more examples, please refer to the source code and test files in the `tests/` directory.
