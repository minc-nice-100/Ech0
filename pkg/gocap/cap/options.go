// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025-2026 lin-snow

package cap

import (
	"crypto/rand"
	"time"

	"github.com/lin-snow/ech0/pkg/gocap/store"
)

// Option mutates engine configuration during construction.
type Option func(*config)

type config struct {
	challengeTTL      time.Duration
	redeemTTL         time.Duration
	gcInterval        time.Duration
	secretPepper      []byte
	customStore       store.Store
	rateLimit         RateLimitConfig
	rateLimitOnRedeem bool
	rateLimitOnVerify bool
	enableCORS        bool
	ipHeader          string
	maxBodyBytes      int64
}

func defaultConfig() config {
	pepper := make([]byte, 32)
	_, _ = rand.Read(pepper)
	return config{
		challengeTTL: 15 * time.Minute,
		redeemTTL:    2 * time.Hour,
		gcInterval:   2 * time.Second,
		secretPepper: pepper,
		rateLimit: RateLimitConfig{
			Max:    30,
			Window: 5 * time.Second,
			Scope:  "cap",
		},
		rateLimitOnRedeem: false,
		rateLimitOnVerify: false,
		maxBodyBytes:      1 << 20,
	}
}

// WithChallengeTTL sets the challenge token validity duration.
func WithChallengeTTL(ttl time.Duration) Option {
	return func(c *config) {
		c.challengeTTL = ttl
	}
}

// WithRedeemTTL sets the redeem token validity duration.
func WithRedeemTTL(ttl time.Duration) Option {
	return func(c *config) {
		c.redeemTTL = ttl
	}
}

// WithGCInterval sets the background cleanup interval for in-memory storage.
func WithGCInterval(interval time.Duration) Option {
	return func(c *config) {
		c.gcInterval = interval
	}
}

// WithSecretPepper sets the pepper used to hash site secrets.
func WithSecretPepper(pepper []byte) Option {
	return func(c *config) {
		if len(pepper) > 0 {
			c.secretPepper = append([]byte(nil), pepper...)
		}
	}
}

// WithStore injects a custom store implementation.
func WithStore(st store.Store) Option {
	return func(c *config) {
		c.customStore = st
	}
}

// WithInMemoryStore selects the default in-memory store.
func WithInMemoryStore() Option {
	return func(c *config) {
		c.customStore = nil
	}
}

// WithRateLimit configures fixed-window rate limiting parameters.
func WithRateLimit(max int, window time.Duration) Option {
	return func(c *config) {
		c.rateLimit.Max = max
		c.rateLimit.Window = window
	}
}

// WithRateLimitScope sets the scope prefix for rate-limit keys.
func WithRateLimitScope(scope string) Option {
	return func(c *config) {
		c.rateLimit.Scope = scope
	}
}

// WithEnableCORS enables permissive CORS handling for the HTTP handler.
func WithEnableCORS(enabled bool) Option {
	return func(c *config) {
		c.enableCORS = enabled
	}
}

// WithIPHeader sets the request header used to extract client IP.
func WithIPHeader(header string) Option {
	return func(c *config) {
		c.ipHeader = header
	}
}

// WithRateLimitOnRedeem toggles rate limiting on redeem requests.
func WithRateLimitOnRedeem(enabled bool) Option {
	return func(c *config) {
		c.rateLimitOnRedeem = enabled
	}
}

// WithRateLimitOnSiteVerify toggles rate limiting on siteverify requests.
func WithRateLimitOnSiteVerify(enabled bool) Option {
	return func(c *config) {
		c.rateLimitOnVerify = enabled
	}
}

// WithMaxBodyBytes sets the per-request JSON body size limit.
func WithMaxBodyBytes(n int64) Option {
	return func(c *config) {
		if n > 0 {
			c.maxBodyBytes = n
		}
	}
}
