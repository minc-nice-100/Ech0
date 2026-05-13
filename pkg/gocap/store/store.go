package store

import "time"

// Site holds all runtime config required for one site key.
type Site struct {
	SiteKey          string
	SecretHash       []byte
	JWTSecret        []byte
	Difficulty       int
	ChallengeCount   int
	SaltSize         int
	BlockOnRateLimit bool
}

// Store abstracts runtime state operations required by core logic.
type Store interface {
	// UpsertSite creates or updates a site configuration.
	UpsertSite(site Site) error
	// GetSite returns a site configuration by key.
	GetSite(siteKey string) (Site, bool)
	// DeleteSite removes one site configuration by key.
	DeleteSite(siteKey string) error

	// TryMarkChallengeSigUsed performs an atomic check-and-set.
	// Returns true when mark is written, false when sig already active.
	TryMarkChallengeSigUsed(sig string, expiresAt time.Time, now time.Time) (bool, error)

	// StoreRedeemToken stores one redeem token with expiration.
	StoreRedeemToken(siteKey, token string, expiresAt time.Time) error

	// ConsumeRedeemToken atomically reads and deletes a redeem token.
	// found=false when token does not exist.
	// expired=true when token exists but is expired at consume time.
	ConsumeRedeemToken(siteKey, token string, now time.Time) (found bool, expired bool, err error)

	// AllowRateLimit performs fixed-window counting for key within scope.
	// Returns allowed and remaining count in current window.
	AllowRateLimit(scope, key string, max int, window time.Duration, now time.Time) (bool, int, error)

	// Close releases resources associated with the store.
	Close() error
}
