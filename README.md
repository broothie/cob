#  `cob` 🌽 COmmand Builder

[![Go Reference](https://pkg.go.dev/badge/github.com/broothie/cob.svg)](https://pkg.go.dev/github.com/broothie/cob)

`cob` is a Go package for building shell commands.

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
