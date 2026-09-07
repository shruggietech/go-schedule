package commandexample

import (
	"context"
	"os/exec"
)

func runSuggestion(ctx context.Context, suggestion Suggestion) (string, error) {
	command := exec.CommandContext(ctx, suggestion.Program, suggestion.Args...)
	configureTestCommand(command)
	output, err := command.CombinedOutput()
	return string(output), err
}
