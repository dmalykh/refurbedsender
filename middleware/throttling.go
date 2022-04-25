package middleware

import (
	"context"
	"github.com/dmalykh/refurbedsender/sender"
	"sync"
	"time"
)

func WithThrottlingMiddleware(max, refill uint64, every time.Duration) sender.Option {
	return func(f sender.ProceedFunc) sender.ProceedFunc {
		return ThrottlingMiddleware(f, max, refill, every)
	}
}

func ThrottlingMiddleware(f sender.ProceedFunc, max, refill uint64, every time.Duration) sender.ProceedFunc {
	var throttleOnce sync.Once
	var tokens = max
	var m sync.Mutex

	return func(ctx context.Context, message sender.Message) error {
		throttleOnce.Do(func() {
			var ticker = time.NewTicker(every)
			go func() {
				defer ticker.Stop()

				for {
					select {
					case <-ctx.Done():
						return

					case <-ticker.C:
						m.Lock()
						t := tokens + refill
						if t > max {
							t = max
						}
						tokens = t
						m.Unlock()
					}
				}
			}()
		})

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				if func() bool {
					m.Lock()
					defer m.Unlock()
					if tokens > 0 {
						tokens--

						return true
					}

					return false
				}() {
					return f(ctx, message)
				}
			}
		}
	}
}
