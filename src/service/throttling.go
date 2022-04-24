package service

import (
	"context"
	"github.com/dmalykh/refurbedsender/sender"
	"sync"
	"time"
)

// The throttlingMiddleware
func (s *Sender) throttlingMiddleware(ctx context.Context, f sender.ProceedFunc, m sender.Message) error {
	var errChan = make(chan error)
	defer close(errChan)
	var wg = new(sync.WaitGroup)
	var once sync.Once
	for {
		select {
		case <-ctx.Done():
			// Waiting for correct errChan closing
			wg.Wait()

			return nil
		case err := <-errChan:
			return err
		default:
			func() {
				s.lock.Lock()
				defer s.lock.Unlock()
				if s.tokens > 0 {
					once.Do(func() {
						s.tokens--
						go func() {
							wg.Add(1)
							defer wg.Done()
							errChan <- f(ctx, m)
						}()
					})
				}
			}()
		}
	}
}

func (s *Sender) throttle(ctx context.Context, max uint64, refill uint64, d time.Duration) {
	s.throttleOnce.Do(func() {
		s.tokens = max
		go func() {
			var ticker = time.NewTicker(d)
			for {
				select {
				case <-ticker.C:
					s.lock.Lock()
					// Add tokens
					var t = s.tokens + refill
					if t > max {
						t = max
					}
					s.tokens = t
					s.lock.Unlock()
				case <-ctx.Done():
					ticker.Stop()

					return
				}
			}
		}()
	})
}
