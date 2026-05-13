package cap

import "time"

// SiteRegistration defines the site-level configuration used at registration time.
type SiteRegistration struct {
	SiteKey        string
	Secret         string
	Difficulty     int
	ChallengeCount int
	SaltSize       int
}

// RateLimitConfig controls fixed-window rate limiting for incoming requests.
type RateLimitConfig struct {
	Max    int
	Window time.Duration
	Scope  string
}
