package cob

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/broothie/option"
)

// WithArgs adds arguments to the command.
func WithArgs(args ...string) option.Func[*exec.Cmd] {
	return func(cmd *exec.Cmd) (*exec.Cmd, error) {
		cmd.Args = append(cmd.Args, args...)
		return cmd, nil
	}
}

// WithEnv adds an environment variable to the command.
func WithEnv(key, value string) option.Func[*exec.Cmd] {
	return func(cmd *exec.Cmd) (*exec.Cmd, error) {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
		return cmd, nil
	}
}

// WithDir sets the working directory for the command (single value).
func WithDir(dir string) option.Func[*exec.Cmd] {
	return func(cmd *exec.Cmd) (*exec.Cmd, error) {
		cmd.Dir = dir
		return cmd, nil
	}
}

// WithStdin adds standard inputs to the command (using MultiReader).
func WithStdin(stdins ...io.Reader) option.Func[*exec.Cmd] {
	return func(cmd *exec.Cmd) (*exec.Cmd, error) {
		var readers []io.Reader
		if cmd.Stdin != nil {
			readers = append(readers, cmd.Stdin)
		}
		readers = append(readers, stdins...)
		cmd.Stdin = io.MultiReader(readers...)
		return cmd, nil
	}
}

// WithStdout adds standard outputs to the command (using MultiWriter).
func WithStdout(stdouts ...io.Writer) option.Func[*exec.Cmd] {
	return func(cmd *exec.Cmd) (*exec.Cmd, error) {
		var writers []io.Writer
		if cmd.Stdout != nil {
			writers = append(writers, cmd.Stdout)
		}
		writers = append(writers, stdouts...)
		cmd.Stdout = io.MultiWriter(writers...)
		return cmd, nil
	}
}

// WithStderr adds standard errors to the command (using MultiWriter).
func WithStderr(stderrs ...io.Writer) option.Func[*exec.Cmd] {
	return func(cmd *exec.Cmd) (*exec.Cmd, error) {
		var writers []io.Writer
		if cmd.Stderr != nil {
			writers = append(writers, cmd.Stderr)
		}
		writers = append(writers, stderrs...)
		cmd.Stderr = io.MultiWriter(writers...)
		return cmd, nil
	}
}

// WithExtraFiles adds extra files to the command.
func WithExtraFiles(extraFiles ...*os.File) option.Func[*exec.Cmd] {
	return func(cmd *exec.Cmd) (*exec.Cmd, error) {
		cmd.ExtraFiles = append(cmd.ExtraFiles, extraFiles...)
		return cmd, nil
	}
}

// WithSysProcAttr sets the system process attributes for the command (single value).
func WithSysProcAttr(attr *syscall.SysProcAttr) option.Func[*exec.Cmd] {
	return func(cmd *exec.Cmd) (*exec.Cmd, error) {
		cmd.SysProcAttr = attr
		return cmd, nil
	}
}

// WithWaitDelay sets the wait delay for the command (single value).
func WithWaitDelay(delay time.Duration) option.Func[*exec.Cmd] {
	return func(cmd *exec.Cmd) (*exec.Cmd, error) {
		cmd.WaitDelay = delay
		return cmd, nil
	}
}
