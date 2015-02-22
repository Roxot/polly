package polly

import (
	"fmt"
	"os"
	"strings"
)

const (
	POLLY_HOME = "POLLY_HOME"
)

func GetPollyHome() (string, error) {
	env := os.Getenv(POLLY_HOME)

	if env == "" {
		err := fmt.Errorf("$%s not set.", POLLY_HOME)
		return env, err
	}

	if !strings.HasSuffix(env, "/") {
		env = env + "/"
	}

	return env, nil
}
