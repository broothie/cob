package cob

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
	"unicode"

	"github.com/broothie/option"
)

// SetArgs sets the arguments for the command.
func SetArgs(args ...string) option.Func[*exec.Cmd] {
	return func(cmd *exec.Cmd) (*exec.Cmd, error) {
		cmd.Args = args
		return cmd, nil
	}
}

// AddArgs adds arguments to the command.
func AddArgs(args ...string) option.Func[*exec.Cmd] {
	return func(cmd *exec.Cmd) (*exec.Cmd, error) {
		return SetArgs(append(cmd.Args, args...)...).Apply(cmd)
	}
}

// SetEnv sets the environment variables for the command.
func SetEnv(env ...string) option.Func[*exec.Cmd] {
	return func(cmd *exec.Cmd) (*exec.Cmd, error) {
		cmd.Env = env
		return cmd, nil
	}
}

// validateEnvKey validates environment variable key according to POSIX standards.
func validateEnvKey(key string) error {
	if key == "" {
		return fmt.Errorf("environment variable key cannot be empty")
	}

	if strings.ContainsRune(key, '=') {
		return fmt.Errorf("environment variable key cannot contain '=' character")
	}

	if strings.ContainsRune(key, '\x00') {
		return fmt.Errorf("environment variable key cannot contain null bytes")
	}

	// POSIX compliant: must start with letter or underscore
	if !unicode.IsLetter(rune(key[0])) && key[0] != '_' {
		return fmt.Errorf("environment variable key must start with letter or underscore, got: %q", key[0])
	}

	// POSIX compliant: can only contain letters, digits, and underscores
	for i, r := range key {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return fmt.Errorf("environment variable key can only contain letters, digits, and underscores, found invalid character %q at position %d", r, i)
		}
	}

	return nil
}

// validateEnvValue validates environment variable value.
func validateEnvValue(value string) error {
	if strings.ContainsRune(value, '\x00') {
		return fmt.Errorf("environment variable value cannot contain null bytes")
	}

	return nil
}

// AddEnv adds an environment variable to the command.
func AddEnv(key, value string) option.Func[*exec.Cmd] {
	return func(cmd *exec.Cmd) (*exec.Cmd, error) {
		if err := validateEnvKey(key); err != nil {
			return nil, fmt.Errorf("invalid environment variable key: %w", err)
		}

		if err := validateEnvValue(value); err != nil {
			return nil, fmt.Errorf("invalid environment variable value: %w", err)
		}

		return SetEnv(append(cmd.Env, fmt.Sprintf("%s=%s", key, value))...).Apply(cmd)
	}
}

// SetDir sets the working directory for the command.
func SetDir(dir string) option.Func[*exec.Cmd] {
	return func(cmd *exec.Cmd) (*exec.Cmd, error) {
		cmd.Dir = dir
		return cmd, nil
	}
}

// SetStdin sets the standard input for the command.
func SetStdin(stdin io.Reader) option.Func[*exec.Cmd] {
	return func(cmd *exec.Cmd) (*exec.Cmd, error) {
		cmd.Stdin = stdin
		return cmd, nil
	}
}

// AddStdins adds standard inputs to the command.
func AddStdins(stdins ...io.Reader) option.Func[*exec.Cmd] {
	return func(cmd *exec.Cmd) (*exec.Cmd, error) {
		if cmd.Stdin != nil {
			stdins = append([]io.Reader{cmd.Stdin}, stdins...)
		}

		return SetStdin(io.MultiReader(stdins...)).Apply(cmd)
	}
}

// SetStdout sets the standard output for the command.
func SetStdout(stdout io.Writer) option.Func[*exec.Cmd] {
	return func(cmd *exec.Cmd) (*exec.Cmd, error) {
		cmd.Stdout = stdout
		return cmd, nil
	}
}

// AddStdouts adds standard outputs to the command.
func AddStdouts(stdouts ...io.Writer) option.Func[*exec.Cmd] {
	return func(cmd *exec.Cmd) (*exec.Cmd, error) {
		if cmd.Stdout != nil {
			stdouts = append([]io.Writer{cmd.Stdout}, stdouts...)
		}

		return SetStdout(io.MultiWriter(stdouts...)).Apply(cmd)
	}
}

// SetStderr sets the standard error for the command.
func SetStderr(stderr io.Writer) option.Func[*exec.Cmd] {
	return func(cmd *exec.Cmd) (*exec.Cmd, error) {
		cmd.Stderr = stderr
		return cmd, nil
	}
}

// AddStderrs adds standard errors to the command.
func AddStderrs(stderrs ...io.Writer) option.Func[*exec.Cmd] {
	return func(cmd *exec.Cmd) (*exec.Cmd, error) {
		if cmd.Stderr != nil {
			stderrs = append([]io.Writer{cmd.Stderr}, stderrs...)
		}

		return SetStderr(io.MultiWriter(stderrs...)).Apply(cmd)
	}
}

// SetExtraFiles sets the extra files for the command.
func SetExtraFiles(extraFiles ...*os.File) option.Func[*exec.Cmd] {
	return func(cmd *exec.Cmd) (*exec.Cmd, error) {
		cmd.ExtraFiles = extraFiles
		return cmd, nil
	}
}

// AddExtraFiles adds extra files to the command.
func AddExtraFiles(extraFiles ...*os.File) option.Func[*exec.Cmd] {
	return func(cmd *exec.Cmd) (*exec.Cmd, error) {
		return SetExtraFiles(append(cmd.ExtraFiles, extraFiles...)...).Apply(cmd)
	}
}

// SetSysProcAttr sets the system process attributes for the command.
func SetSysProcAttr(attr *syscall.SysProcAttr) option.Func[*exec.Cmd] {
	return func(cmd *exec.Cmd) (*exec.Cmd, error) {
		cmd.SysProcAttr = attr
		return cmd, nil
	}
}

// SetWaitDelay sets the wait delay for the command.
func SetWaitDelay(delay time.Duration) option.Func[*exec.Cmd] {
	return func(cmd *exec.Cmd) (*exec.Cmd, error) {
		cmd.WaitDelay = delay
		return cmd, nil
	}
}
