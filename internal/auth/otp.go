package auth

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type OTP struct {
	Key       string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type OTPCollection struct {
	mu sync.RWMutex
	m  map[string]OTP
}

func NewOTPCollection(ctx context.Context, cleanupInterval time.Duration) *OTPCollection {
	rm := OTPCollection{
		m: make(map[string]OTP),
	}

	go rm.OTPRemover(cleanupInterval)
	return &rm
}

func (otpc *OTPCollection) NewOTP(aliveTime time.Duration) OTP {
	otp := OTP{
		Key:       uuid.NewString(),
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(aliveTime),
	}
	otpc.mu.Lock()
	defer otpc.mu.Unlock()

	otpc.m[otp.Key] = otp
	return otp
}

func (otpc *OTPCollection) OTPRemover(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		<-ticker.C
		otpc.mu.Lock()
		for key, otp := range otpc.m {
			if otp.ExpiresAt.Before(time.Now()) {
				delete(otpc.m, key)
			}
		}
		otpc.mu.Unlock()
	}
}
