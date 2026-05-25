package multiurl_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/yourorg/logslice/internal/multiurl"
)

func BenchmarkFan100Sources(b *testing.B) {
	urls := make([]string, 100)
	for i := range urls {
		urls[i] = fmt.Sprintf("http://host%d", i)
	}

	streamer := func(ctx context.Context, url string, out chan<- multiurl.Line) error {
		for i := 0; i < 10; i++ {
			out <- multiurl.Line{Source: url, Data: fmt.Sprintf(`{"msg":"line%d"}`, i)}
		}
		return nil
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ctx := context.Background()
		lineCh, _ := multiurl.Fan(ctx, urls, streamer)
		count := 0
		for range lineCh {
			count++
		}
		if count != 1000 {
			b.Fatalf("expected 1000 lines, got %d", count)
		}
	}
}
