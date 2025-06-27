package cob_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/broothie/option"
	"github.com/broothie/test"

	"github.com/broothie/cob"
)

func TestStart(t *testing.T) {
	output := new(bytes.Buffer)
	cmd, err := cob.Start(context.TODO(), "echo", cob.AddArgs("Hello, World"), cob.SetStdout(output))
	test.MustNoError(t, err)
	test.NoError(t, cmd.Wait())
	test.Equal(t, output.String(), "Hello, World\n")
}

func TestOutput(t *testing.T) {
	type Args struct {
		ctx     context.Context
		name    string
		options option.Options[*exec.Cmd]
	}

	type TestCase struct {
		args           Args
		expectedStdout string
		expectedError  string
	}

	testCases := map[string]TestCase{
		"Hello, World": {
			args: Args{
				ctx:     context.TODO(),
				name:    "echo",
				options: option.NewOptions(cob.AddArgs("Hello, World")),
			},
			expectedStdout: "Hello, World\n",
		},
		"env vars work": {
			args: Args{
				ctx:  context.TODO(),
				name: "bash",
				options: option.NewOptions(
					cob.AddArgs("-c", `echo '$NAME' is "$NAME"`),
					cob.AddEnv("NAME", "alice"),
				),
			},
			expectedStdout: "$NAME is alice\n",
		},
		"kitchen sink, happy path": {
			args: Args{
				ctx:  context.TODO(),
				name: "tr",
				options: option.NewOptions(
					cob.AddArgs("[:upper:]"),
					cob.AddArgs("[:lower:]"),
					cob.SetDir("."),
					cob.AddEnv("KEY", "value"),
					cob.AddStdins(strings.NewReader("hello, world")),
					cob.AddStdouts(os.Stdout),
					cob.AddStderrs(os.Stderr),
					cob.AddExtraFiles(os.Stdout),
					cob.SetSysProcAttr(&syscall.SysProcAttr{Setsid: true}),
					cob.SetWaitDelay(time.Second),
				),
			},
			expectedStdout: "hello, world",
		},
		"erroring option": {
			args: Args{
				ctx:  context.TODO(),
				name: "echo",
				options: option.NewOptions(
					option.Func[*exec.Cmd](func(cmd *exec.Cmd) (*exec.Cmd, error) {
						return nil, fmt.Errorf("some error")
					}),
				),
			},
			expectedError: "failed to apply option 0: some error\nerror building command",
		},
		"AddEnv with empty key": {
			args: Args{
				ctx:  context.TODO(),
				name: "echo",
				options: option.NewOptions(
					cob.AddEnv("", "value"),
				),
			},
			expectedError: "failed to apply option 0: invalid environment variable key: environment variable key cannot be empty\nerror building command",
		},
		"AddEnv with key containing equals": {
			args: Args{
				ctx:  context.TODO(),
				name: "echo",
				options: option.NewOptions(
					cob.AddEnv("KEY=BAD", "value"),
				),
			},
			expectedError: "failed to apply option 0: invalid environment variable key: environment variable key cannot contain '=' character\nerror building command",
		},
		"AddEnv with key containing null byte": {
			args: Args{
				ctx:  context.TODO(),
				name: "echo",
				options: option.NewOptions(
					cob.AddEnv("KEY\x00", "value"),
				),
			},
			expectedError: "failed to apply option 0: invalid environment variable key: environment variable key cannot contain null bytes\nerror building command",
		},
		"AddEnv with key starting with digit": {
			args: Args{
				ctx:  context.TODO(),
				name: "echo",
				options: option.NewOptions(
					cob.AddEnv("9KEY", "value"),
				),
			},
			expectedError: "failed to apply option 0: invalid environment variable key: environment variable key must start with letter or underscore, got: '9'\nerror building command",
		},
		"AddEnv with key containing invalid character": {
			args: Args{
				ctx:  context.TODO(),
				name: "echo",
				options: option.NewOptions(
					cob.AddEnv("KEY-BAD", "value"),
				),
			},
			expectedError: "failed to apply option 0: invalid environment variable key: environment variable key can only contain letters, digits, and underscores, found invalid character '-' at position 3\nerror building command",
		},
		"AddEnv with value containing null byte": {
			args: Args{
				ctx:  context.TODO(),
				name: "echo",
				options: option.NewOptions(
					cob.AddEnv("KEY", "value\x00"),
				),
			},
			expectedError: "failed to apply option 0: invalid environment variable value: environment variable value cannot contain null bytes\nerror building command",
		},
		"AddEnv with valid key starting with underscore": {
			args: Args{
				ctx:     context.TODO(),
				name:    "bash",
				options: option.NewOptions(cob.AddArgs("-c", "echo $_TEST"), cob.AddEnv("_TEST", "underscore")),
			},
			expectedStdout: "underscore\n",
		},
		"AddEnv with valid key containing digits": {
			args: Args{
				ctx:     context.TODO(),
				name:    "bash",
				options: option.NewOptions(cob.AddArgs("-c", "echo $TEST123"), cob.AddEnv("TEST123", "digits")),
			},
			expectedStdout: "digits\n",
		},
		"AddEnv with valid key mixed case": {
			args: Args{
				ctx:     context.TODO(),
				name:    "bash",
				options: option.NewOptions(cob.AddArgs("-c", "echo $MyVar_123"), cob.AddEnv("MyVar_123", "mixed")),
			},
			expectedStdout: "mixed\n",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			stdout, _, _, err := cob.Output(testCase.args.ctx, testCase.args.name, testCase.args.options...)
			if testCase.expectedError == "" {
				test.MustNoError(t, err)
				test.Equal(t, stdout.String(), testCase.expectedStdout)
			} else {
				test.Equal(t, err.Error(), testCase.expectedError)
			}
		})
	}
}

func ExampleRun() {
	// Output:
	// Hello, World

	cob.Run(context.TODO(), "echo",
		cob.AddArgs("Hello,"),
		cob.AddArgs("World"),
		cob.AddEnv("SHELL", "bash"),
		cob.SetStdin(os.Stdin),
		cob.SetStdout(os.Stdout),
	)
}

func ExamplePipe() {
	// Output:
	// hello, world

	pipeOutput, pipeInput := io.Pipe()

	go func() {
		cob.Run(context.TODO(), "echo",
			cob.AddArgs("Hello, World"),
			cob.SetStdout(pipeInput),
		)

		pipeInput.Close()
	}()

	cob.Run(context.TODO(), "tr",
		cob.AddArgs("[:upper:]"),
		cob.AddArgs("[:lower:]"),
		cob.SetStdin(pipeOutput),
		cob.SetStdout(os.Stdout),
	)
}
