package eerror

import (
	"errors"
)

var (
	MissingAuthenticationError = errors.New("missing authentication")
	InvalidAuthenticationError = errors.New("invalid authentication")
)