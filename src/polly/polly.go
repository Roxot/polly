package polly

import (
	"fmt"
	"os"
	"strings"
)

const (
	POLLY_HOME_ENV = "POLLY_HOME"
)

func GetPollyHome() (string, error) {
	var env string
	env = os.Getenv(POLLY_HOME_ENV)

	if env == "" {
		err := fmt.Errorf("$%s not set.", POLLY_HOME_ENV)
		return env, err
	}

	if !strings.HasSuffix(env, "/") {
		env = env + "/"
	}

	return env, nil
}
