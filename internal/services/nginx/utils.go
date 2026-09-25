package nginx

import (
	"context"
	"errors"
	"os/exec"

	"github.com/tuanta7/uon/pkg/command"
)

func commandSucceedsOrHasStateExit(ctx context.Context, name string, stateExitCode int, args ...string) (bool, error) {
	err := command.Run(ctx, name, args...)
	if hasExitCode(err, stateExitCode) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func hasExitCode(err error, code int) bool {
	if err == nil {
		return false
	}

	var exitErr *exec.ExitError
	return errors.As(err, &exitErr) && exitErr.ExitCode() == code
}
