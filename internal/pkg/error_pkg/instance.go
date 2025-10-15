package error_pkg

import "fmt"

var ErrNotFoundInstance = fmt.Errorf("resource was not found error")

var ErrBadRequestInstance = fmt.Errorf("invalid request error")

var ErrConflicInstance = fmt.Errorf("improcesable request error")

var ErrForbiddenInstance = fmt.Errorf("forbidden resource error")

var ErrUnauthenticatedInstance = fmt.Errorf("unanthenticated error")

var ErrAbortedInstance = fmt.Errorf("aborted error")

var ErrRateLimitInstance = fmt.Errorf("rate limit error")

var ErrExternalServiceInstance = fmt.Errorf("external service error")

var ErrInternalInstance = fmt.Errorf("internal error")

var ErrUnknownInstance = fmt.Errorf("unknown error")
