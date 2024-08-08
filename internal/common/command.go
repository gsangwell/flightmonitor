package common

import (
        "log/slog"
        "os/exec"
        "bytes"
        "strings"
	"errors"
)

func RunCommand(command string, args ...string) (int, string, string, error) {
	var stdout, stderr bytes.Buffer
        cmd := exec.Command(command, args...)
        cmd.Stdout = &stdout
        cmd.Stderr = &stderr
        err := cmd.Run()

        if err != nil {
		if Debug {
	                slog.Error("common.RunCommand", "command", command, "args", strings.Join(args, " "), "stderr", stderr.String(), "stdout", stdout.String(), "err", err)
		}
                return -1, stdout.String(), stderr.String(), errors.New("Failed to run command.")
        }

        exit_code := cmd.ProcessState.ExitCode()

	slog.Debug("common.RunCommand", "command", command, "args", strings.Join(args, " "), "exitcode", exit_code, "stderr", stderr.String(), "stdout", stdout.String(), "err", err)
        return exit_code, stdout.String(), stderr.String(), err
}
