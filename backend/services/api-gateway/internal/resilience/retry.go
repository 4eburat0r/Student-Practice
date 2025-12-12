package resilience

import (
	"time"

	"api-gateway/pkg/logger"

	"github.com/avast/retry-go"
)

func RetryWithBackoff(fn func() error) error {
	return retry.Do(
		fn,
		retry.Attempts(3),
		retry.Delay(500*time.Millisecond),
		retry.MaxDelay(5*time.Second),
		retry.DelayType(retry.BackOffDelay),
		retry.OnRetry(func(n uint, err error) {
			logger.Warn("Retrying request",
				logger.Uint("attempt", n+1),
				logger.Err(err),
			)
		}),
	)
}
