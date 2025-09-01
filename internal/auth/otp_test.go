package auth

import (
	"context"
	"sync"
	"testing"
	"time"
)

type casesTest struct {
	want            bool
	timeAlive       time.Duration
	cleanupInterval time.Duration
	checkAfter      time.Duration
}

func TestNewOTPMap(t *testing.T) {
	t.Run("Creates OTPs and correctly removes it after they expire", func(t *testing.T) {
		cases := []casesTest{
			{
				want:            true,
				timeAlive:       10 * time.Millisecond,
				cleanupInterval: 2 * time.Millisecond,
				checkAfter:      5 * time.Millisecond,
			},
			{
				want:            false,
				timeAlive:       10 * time.Millisecond,
				cleanupInterval: 2 * time.Millisecond,
				checkAfter:      12 * time.Millisecond,
			},
		}

		var wg sync.WaitGroup

		for _, c := range cases {
			otpc := NewOTPCollection(context.Background(), c.cleanupInterval)
			otp := otpc.NewOTP(c.timeAlive)

			wg.Add(1)
			go func() {
				ticker := time.NewTicker(c.checkAfter)
				<-ticker.C
				wg.Done()

				otpc.mu.RLock()
				if _, ok := otpc.m[otp.Key]; ok != c.want {
					t.Errorf("Got %v, want %v", ok, c.want)
				}
				otpc.mu.RUnlock()
			}()
		}
		wg.Wait()
	})
}
