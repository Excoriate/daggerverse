# Gotest Module for Dagger

A Dagger module that provides comprehensive Go testing capabilities with full control over test execution, build options, and test configurations.

## Configuration 🛠️

Through the [Dagger CLI](https://docs.dagger.io/cli/465058/install), you can configure:

- ⚙️ `version`: Go version to use (e.g., "1.22.5"). Defaults to `latest`
- ⚙️ `image`: Custom Go image. Defaults to `golang:alpine`
- ⚙️ `envVarsFromHost`: Environment variables to pass from host
- ⚙️ `ctr`: Base container for customization

## Features 🎨

| Feature             | Description                             | Example                                                  |
| ------------------- | --------------------------------------- | -------------------------------------------------------- |
| Test Execution      | Run Go tests with comprehensive options | `dagger call test --source=. --cover=true`               |
| Build Configuration | Control build flags and options         | `dagger call test --race=true --buildTags="integration"` |
| Test Filtering      | Filter and control test execution       | `dagger call test --run="TestSpecific" --short=true`     |
| Profiling           | CPU, memory, and block profiling        | `dagger call test --cpuprofile="cpu.prof"`               |
| Benchmarking        | Run and configure benchmarks            | `dagger call test --benchmark="." --benchmem=true`       |

## Usage Examples 🚀

### Basic Test Execution

```bash
dagger call test --source=. --enableDefaultOptions=true
```

### With Coverage and Profiling

```bash
dagger call test \
  --source=. \
  --cover=true \
  --coverprofile="coverage.out" \
  --cpuprofile="cpu.prof" \
  --verbose=true
```

### Race Detection and Build Tags

```bash
dagger call test \
  --source=. \
  --race=true \
  --buildTags="integration" \
  --ldflags="-X main.version=test"
```

### Benchmark Testing

```bash
dagger call test \
  --source=. \
  --benchmark="." \
  --benchmem=true \
  --benchtime="1s" \
  --testCount=3
```

## Available Options

### Build Options

| Option      | Description               |
| ----------- | ------------------------- |
| `race`      | Enable race detection     |
| `msan`      | Enable memory sanitizer   |
| `asan`      | Enable address sanitizer  |
| `buildTags` | Specify build constraints |
| `ldflags`   | Set linker flags          |
| `gcflags`   | Set Go compiler flags     |
| `asmflags`  | Set assembler flags       |
| `trimpath`  | Remove file system paths  |
| `buildMode` | Set build mode            |
| `compiler`  | Specify compiler          |
| `mod`       | Set module mode           |

### Test Options

| Option             | Description                    |
| ------------------ | ------------------------------ |
| `benchmark`        | Run benchmarks matching regexp |
| `benchmem`         | Report memory allocations      |
| `benchtime`        | Run time for benchmarks        |
| `cover`            | Enable coverage analysis       |
| `coverprofile`     | Write coverage profile         |
| `cpuprofile`       | Write CPU profile              |
| `testCount`        | Run tests multiple times       |
| `failfast`         | Stop on first failure          |
| `enableJsonOutput` | Enable JSON output             |
| `parallel`         | Set parallel test count        |
| `run`              | Run tests matching pattern     |
| `short`            | Run in short mode              |
| `timeout`          | Set test timeout               |
| `verbose`          | Enable verbose output          |

## Environment Variables and Secrets

The module supports:

- Setting environment variables via `envVars`
- Passing secrets securely via `secrets`
- Inheriting environment variables from host via `envVarsFromHost`

Example:

```bash
dagger call test \
  --source=. \
  --envVars='["GO_ENV=test", "DEBUG=true"]' \
  --secrets='["MY_SECRET"]'
```

## Testing 🧪

Run the test suite:

```bash
just test gotest
```

## Developer Experience 🛠️

Development commands:

```bash
# Initialize pre-commit hooks
just run-hooks

# Run linting
just lintall gotest

# Run tests
just test gotest

# Run CI pipeline locally
just ci gotest
```

## API Reference

### Main Functions

#### RunTest

Executes Go tests with full configuration options:

```go
RunTest(source *dagger.Directory, packages []string, envVars []string, secrets []*dagger.Secret, ...) (*dagger.Container, error)
```

#### RunTestCmd

Executes tests and returns command output:

```go
RunTestCmd(source *dagger.Directory, packages []string, envVars []string, secrets []*dagger.Secret, ...) (string, error)
```

For detailed API documentation and more examples, see the [Dagger documentation](https://docs.dagger.io).
