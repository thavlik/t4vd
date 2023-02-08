package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

// INPUT_URL=https://www.youtube.com/@weshammer runuser -pptruser -- node -e "$(cat /scripts/query-channel.js)"

func nodeQuery(
	ctx context.Context,
	scriptPath string,
	input string,
	dest interface{},
) error {
	command := fmt.Sprintf(`INPUT_URL="%s" runuser -u pptruser -- node -e "$(cat %s)"`,
		input,
		scriptPath,
	)
	cmd := exec.Command("bash", "-c", command)
	var stdout bytes.Buffer
	cmd.Stdin = os.Stdin
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr
	done := make(chan error, 1)
	go func() {
		done <- cmd.Run()
	}()
	select {
	case <-ctx.Done():
		if err := cmd.Process.Kill(); err != nil {
			return errors.Wrap(err, "kill")
		}
		return errors.Wrap(ctx.Err(), "context")
	case err := <-done:
		if err != nil {
			return errors.Wrap(err, "run")
		}
	}
	if err := json.Unmarshal(stdout.Bytes(), &dest); err != nil {
		return errors.Wrap(err, "unmarshal")
	}
	return nil
}
