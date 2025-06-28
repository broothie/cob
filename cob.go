package cob

import (
	"bytes"
	"context"
	"errors"
	"os/exec"

	"github.com/broothie/option"
)

// ErrCommandBuild is the error returned when building a command.
var ErrCommandBuild = errors.New("error building command")

// New builds a new command.
func New(ctx context.Context, name string, options ...option.Option[*exec.Cmd]) (*exec.Cmd, error) {
	cmd, err := option.Apply(exec.CommandContext(ctx, name), options...)
	if err != nil {
		return nil, errors.Join(err, ErrCommandBuild)
	}

	return cmd, nil
}

// Start starts the command and returns it.
func Start(ctx context.Context, name string, options ...option.Option[*exec.Cmd]) (*exec.Cmd, error) {
	cmd, err := New(ctx, name, options...)
	if err != nil {
		return nil, err
	}

	return cmd, cmd.Start()
}

// Run runs the command and returns the combined output.
func Run(ctx context.Context, name string, options ...option.Option[*exec.Cmd]) (*exec.Cmd, error) {
	cmd, err := New(ctx, name, options...)
	if err != nil {
		return nil, err
	}

	return cmd, cmd.Run()
}

func Output(ctx context.Context, name string, options ...option.Option[*exec.Cmd]) (*bytes.Buffer, *bytes.Buffer, *exec.Cmd, error) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	cmd, err := Run(ctx, name, append(options, WithStdout(stdout), WithStderr(stderr))...)
	return stdout, stderr, cmd, err
}
