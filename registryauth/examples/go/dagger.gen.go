// Code generated by dagger. DO NOT EDIT.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/Excoriate/daggerverse/registryauth/examples/go/internal/dagger"
	"github.com/Excoriate/daggerverse/registryauth/examples/go/internal/querybuilder"
	"github.com/Excoriate/daggerverse/registryauth/examples/go/internal/telemetry"
)

var dag = dagger.Connect()

func Tracer() trace.Tracer {
	return otel.Tracer("dagger.io/sdk.go")
}

// used for local MarshalJSON implementations
var marshalCtx = context.Background()

// called by main()
func setMarshalContext(ctx context.Context) {
	marshalCtx = ctx
	dagger.SetMarshalContext(ctx)
}

type DaggerObject = querybuilder.GraphQLMarshaller

type ExecError = dagger.ExecError

// ptr returns a pointer to the given value.
func ptr[T any](v T) *T {
	return &v
}

// convertSlice converts a slice of one type to a slice of another type using a
// converter function
func convertSlice[I any, O any](in []I, f func(I) O) []O {
	out := make([]O, len(in))
	for i, v := range in {
		out[i] = f(v)
	}
	return out
}

func (r Go) MarshalJSON() ([]byte, error) {
	var concrete struct {
		TestDir *dagger.Directory
	}
	concrete.TestDir = r.TestDir
	return json.Marshal(&concrete)
}

func (r *Go) UnmarshalJSON(bs []byte) error {
	var concrete struct {
		TestDir *dagger.Directory
	}
	err := json.Unmarshal(bs, &concrete)
	if err != nil {
		return err
	}
	r.TestDir = concrete.TestDir
	return nil
}

func main() {
	ctx := context.Background()

	// Direct slog to the new stderr. This is only for dev time debugging, and
	// runtime errors/warnings.
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	})))

	if err := dispatch(ctx); err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}
}

func dispatch(ctx context.Context) error {
	ctx = telemetry.InitEmbedded(ctx, resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("dagger-go-sdk"),
		// TODO version?
	))
	defer telemetry.Close()

	// A lot of the "work" actually happens when we're marshalling the return
	// value, which entails getting object IDs, which happens in MarshalJSON,
	// which has no ctx argument, so we use this lovely global variable.
	setMarshalContext(ctx)

	fnCall := dag.CurrentFunctionCall()
	parentName, err := fnCall.ParentName(ctx)
	if err != nil {
		return fmt.Errorf("get parent name: %w", err)
	}
	fnName, err := fnCall.Name(ctx)
	if err != nil {
		return fmt.Errorf("get fn name: %w", err)
	}
	parentJson, err := fnCall.Parent(ctx)
	if err != nil {
		return fmt.Errorf("get fn parent: %w", err)
	}
	fnArgs, err := fnCall.InputArgs(ctx)
	if err != nil {
		return fmt.Errorf("get fn args: %w", err)
	}

	inputArgs := map[string][]byte{}
	for _, fnArg := range fnArgs {
		argName, err := fnArg.Name(ctx)
		if err != nil {
			return fmt.Errorf("get fn arg name: %w", err)
		}
		argValue, err := fnArg.Value(ctx)
		if err != nil {
			return fmt.Errorf("get fn arg value: %w", err)
		}
		inputArgs[argName] = []byte(argValue)
	}

	result, err := invoke(ctx, []byte(parentJson), parentName, fnName, inputArgs)
	if err != nil {
		return fmt.Errorf("invoke: %w", err)
	}
	resultBytes, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	if err = fnCall.ReturnValue(ctx, dagger.JSON(resultBytes)); err != nil {
		return fmt.Errorf("store return value: %w", err)
	}
	return nil
}
func invoke(ctx context.Context, parentJSON []byte, parentName string, fnName string, inputArgs map[string][]byte) (_ any, err error) {
	_ = inputArgs
	switch parentName {
	case "Go":
		switch fnName {
		case "AllRecipes":
			var parent Go
			err = json.Unmarshal(parentJSON, &parent)
			if err != nil {
				panic(fmt.Errorf("%s: %w", "failed to unmarshal parent object", err))
			}
			return nil, (*Go).AllRecipes(&parent, ctx)
		case "BuiltInRecipes":
			var parent Go
			err = json.Unmarshal(parentJSON, &parent)
			if err != nil {
				panic(fmt.Errorf("%s: %w", "failed to unmarshal parent object", err))
			}
			return nil, (*Go).BuiltInRecipes(&parent, ctx)
		case "Registryauth_PassedEnvVars":
			var parent Go
			err = json.Unmarshal(parentJSON, &parent)
			if err != nil {
				panic(fmt.Errorf("%s: %w", "failed to unmarshal parent object", err))
			}
			return nil, (*Go).Registryauth_PassedEnvVars(&parent, ctx)
		case "Registryauth_OpenTerminal":
			var parent Go
			err = json.Unmarshal(parentJSON, &parent)
			if err != nil {
				panic(fmt.Errorf("%s: %w", "failed to unmarshal parent object", err))
			}
			return (*Go).Registryauth_OpenTerminal(&parent), nil
		case "Registryauth_CreateNetRcFileForGithub":
			var parent Go
			err = json.Unmarshal(parentJSON, &parent)
			if err != nil {
				panic(fmt.Errorf("%s: %w", "failed to unmarshal parent object", err))
			}
			return (*Go).Registryauth_CreateNetRcFileForGithub(&parent, ctx)
		case "Registryauth_RunArbitraryCommand":
			var parent Go
			err = json.Unmarshal(parentJSON, &parent)
			if err != nil {
				panic(fmt.Errorf("%s: %w", "failed to unmarshal parent object", err))
			}
			return (*Go).Registryauth_RunArbitraryCommand(&parent, ctx)
		case "Registryauth_CreateContainer":
			var parent Go
			err = json.Unmarshal(parentJSON, &parent)
			if err != nil {
				panic(fmt.Errorf("%s: %w", "failed to unmarshal parent object", err))
			}
			return (*Go).Registryauth_CreateContainer(&parent, ctx)
		case "":
			var parent Go
			err = json.Unmarshal(parentJSON, &parent)
			if err != nil {
				panic(fmt.Errorf("%s: %w", "failed to unmarshal parent object", err))
			}
			return New(), nil
		default:
			return nil, fmt.Errorf("unknown function %s", fnName)
		}
	case "":
		return dag.Module().
			WithDescription("A generated module for Go example functions\n").
			WithObject(
				dag.TypeDef().WithObject("Go", dagger.TypeDefWithObjectOpts{Description: "Go is a Dagger module that exemplifies the usage of the Registryauth module.\n\nThis module is used to create and manage containers."}).
					WithFunction(
						dag.Function("AllRecipes",
							dag.TypeDef().WithKind(dagger.VoidKind).WithOptional(true)).
							WithDescription("AllRecipes executes all tests.\n\nAllRecipes is a helper method for tests, executing the built-in recipes and\nother specific functionalities of the Registryauth module.")).
					WithFunction(
						dag.Function("BuiltInRecipes",
							dag.TypeDef().WithKind(dagger.VoidKind).WithOptional(true)).
							WithDescription("BuiltInRecipes demonstrates how to run built-in recipes\n\nThis method showcases the use of various built-in recipes provided by the Registryauth\nmodule, including creating a container, running an arbitrary command, and creating a .netrc\nfile for GitHub authentication.\n\nParameters:\n  - ctx: The context for controlling the function's timeout and cancellation.\n\nReturns:\n  - An error if any of the internal methods fail, or nil otherwise.")).
					WithFunction(
						dag.Function("Registryauth_PassedEnvVars",
							dag.TypeDef().WithKind(dagger.VoidKind).WithOptional(true)).
							WithDescription("Registryauth_PassedEnvVars demonstrates how to pass environment variables to the Registryauth module.\n\nThis method configures a Registryauth module to use specific environment variables from the host.")).
					WithFunction(
						dag.Function("Registryauth_OpenTerminal",
							dag.TypeDef().WithObject("Container")).
							WithDescription("Registryauth_OpenTerminal demonstrates how to open an interactive terminal session\nwithin a Registryauth module container.\n\nThis function showcases the initialization and configuration of a\nRegistryauth module container using various options like enabling Cgo,\nutilizing build cache, and including a GCC compiler.\n\nParameters:\n  - None\n\nReturns:\n  - *dagger.Container: A configured Dagger container with an open terminal.\n\nUsage:\n\n\tThis function can be used to interactively debug or inspect the\n\tcontainer environment during test execution.")).
					WithFunction(
						dag.Function("Registryauth_CreateNetRcFileForGithub",
							dag.TypeDef().WithObject("Container")).
							WithDescription("Registryauth_CreateNetRcFileForGithub creates and configures a .netrc file for GitHub authentication.\n\nThis method exemplifies the creation of a .netrc file with credentials for accessing GitHub,\nand demonstrates how to pass this file as a secret to the Registryauth module.\n\nParameters:\n  - ctx: The context for controlling the function's timeout and cancellation.\n\nReturns:\n  - A configured Container with the .netrc file set as a secret, or an error if the process fails.\n\nSteps Involved:\n 1. Define GitHub password as a secret.\n 2. Configure the Registryauth module to use the .netrc file with the provided credentials.\n 3. Run a command inside the container to verify the .netrc file's contents.")).
					WithFunction(
						dag.Function("Registryauth_RunArbitraryCommand",
							dag.TypeDef().WithKind(dagger.StringKind)).
							WithDescription("Registryauth_RunArbitraryCommand runs an arbitrary shell command in the test container.\n\nThis function demonstrates how to execute a shell command within the container\nusing the Registryauth module.\n\nParameters:\n\n\tctx - context for controlling the function lifetime.\n\nReturns:\n\n\tA string containing the output of the executed command, or an error if the command fails or if the output is empty.")).
					WithFunction(
						dag.Function("Registryauth_CreateContainer",
							dag.TypeDef().WithObject("Container")).
							WithDescription("Registryauth_CreateContainer initializes and returns a configured Dagger container.\n\nThis method exemplifies the setup of a container within the Registryauth module using the source directory.\n\nParameters:\n  - ctx: The context for controlling the function's timeout and cancellation.\n\nReturns:\n  - A configured Dagger container if successful, or an error if the process fails.\n\nSteps Involved:\n 1. Configure the Registryauth module with the source directory.\n 2. Run a command inside the container to check the OS information.")).
					WithField("TestDir", dag.TypeDef().WithObject("Directory")).
					WithConstructor(
						dag.Function("New",
							dag.TypeDef().WithObject("Go")).
							WithDescription("New creates a new Tests instance.\n\nIt's the initial constructor for the Tests struct."))), nil
	default:
		return nil, fmt.Errorf("unknown object %s", parentName)
	}
}
