// Package ratelimit provides a token-bucket rate limiter used to cap
// the number of log lines processed per second during streaming.
//
// Usage:
//
//	limiter := ratelimit.New(100) // 100 lines/sec
//	defer limiter.Stop()
//	for line := range lines {
//		if err := limiter.Wait(ctx); err != nil {
//			break
//		}
//		process(line)
//	}
package ratelimit
