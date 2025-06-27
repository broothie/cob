# `cob` 🌽 COmmand Builder

[![Go Reference](https://pkg.go.dev/badge/github.com/broothie/cob.svg)](https://pkg.go.dev/github.com/broothie/cob)
[![codecov](https://codecov.io/gh/broothie/cob/graph/badge.svg?token=FgyhQS4tMX)](https://codecov.io/gh/broothie/cob)
[![gosec](https://github.com/broothie/cob/actions/workflows/gosec.yml/badge.svg)](https://github.com/broothie/cob/actions/workflows/gosec.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/broothie/cob)](https://goreportcard.com/report/github.com/broothie/cob)

`cob` is a Go package for building and running `*exec.Cmd` objects.

## Installation

```shell
go get github.com/broothie/cob@latest
```

## Documentation

Detailed documentation can be found at [pkg.go.dev](https://pkg.go.dev/github.com/broothie/cob).

## Usage

```go
// Easily build `*exec.Cmd` objects:
cmd, err := cob.New(ctx, "echo",
	cob.AddArgs("Hello", "World"),
	cob.AddEnv("SHELL", "bash"),
	cob.SetStdout(os.Stdout),
)

// Or, run them and easily get `stdout` and `stderr`:
stdout, stderr, cmd, err := cob.Output(ctx, "echo", cob.AddArgs("Hello", "World"))
```

Here is a list of all available options:

```go
cmd, err := cob.New(ctx, "echo",
	cob.SetArgs("Hello", "World"),         // Set args directly
	cob.AddArgs("more", "args"),           // Add to existing args
	cob.SetEnv("SHELL=bash"),              // Set env directly
	cob.AddEnv("SHELL", "bash"),           // Add to existing env (validates key/value)
	cob.SetStdin(os.Stdin),                // Set stdin directly
	cob.AddStdins(someReader),             // Add to existing stdin
	cob.SetStdout(os.Stdout),              // Set stdout directly
	cob.AddStdouts(someWriter),            // Add to existing stdout
	cob.SetStderr(os.Stderr),              // Set stderr directly
	cob.AddStderrs(errorWriter),           // Add to existing stderr
	cob.SetDir("/tmp"),                    // Set working directory
	cob.SetExtraFiles(os.Stdin),           // Set extra files directly
	cob.AddExtraFiles(someFile),           // Add to existing extra files
	cob.SetSysProcAttr(&syscall.SysProcAttr{}),  // Set process attributes
	cob.SetWaitDelay(time.Second),         // Set timeout for Wait()
)
```

## Environment Variable Validation

The `AddEnv()` function validates environment variable keys and values to prevent common errors:

### Key Validation
- Must be non-empty
- Must start with a letter or underscore
- Can only contain letters, digits, and underscores
- Cannot contain `=` or null bytes

### Value Validation
- Cannot contain null bytes

### Examples

```go
// Valid keys
cob.AddEnv("PATH", "/usr/bin")           // ✓ starts with letter
cob.AddEnv("_PRIVATE", "secret")         // ✓ starts with underscore  
cob.AddEnv("VAR_123", "value")           // ✓ contains digits and underscore

// Invalid keys (will return error)
cob.AddEnv("", "value")                  // ✗ empty key
cob.AddEnv("123VAR", "value")            // ✗ starts with digit
cob.AddEnv("VAR-NAME", "value")          // ✗ contains hyphen
cob.AddEnv("KEY=VAL", "value")           // ✗ contains equals sign
cob.AddEnv("KEY\x00", "value")           // ✗ contains null byte

// Invalid values (will return error)
cob.AddEnv("KEY", "value\x00")           // ✗ value contains null byte
```

## Why?

The `*exec.Cmd` API isn't terrible by any means,
I just find building the actual `*exec.Cmd` objects to be somewhat cumbersome.

With this library we go from this:

```go
stdout := new(bytes.Buffer)
stderr := new(bytes.Buffer)

cmd := exec.CommandContext(ctx, "command", "arg1", "arg2")
cmd.Stdin = os.Stdin
cmd.Stdout = io.MultiWriter(stdout, os.Stdout)
cmd.Stderr = stderr
cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))

if err := cmd.Run(); err != nil {
	return err
}
```

to this:

```go
stdout, stderr, cmd, err := cmd.Output(ctx, "command",
	cob.AddArgs("arg1", "arg2"),
	cob.SetStdin(os.Stdin),
	cob.AddStdouts(os.Stdout),
	cob.AddEnv(key, value),
)
if err != nil {
	return err
}
```

which I personally find clearer and more concise.
