package core

// ChallengeParams carries server-issued challenge parameters.
type ChallengeParams struct {
	C int `json:"c"`
	S int `json:"s"`
	D int `json:"d"`
}

// ChallengeResponse is returned by the challenge endpoint.
type ChallengeResponse struct {
	Challenge ChallengeParams `json:"challenge"`
	Token     string          `json:"token"`
	Expires   int64           `json:"expires"`
}

// RedeemRequest is the request payload for redeeming a challenge token.
type RedeemRequest struct {
	Token        string `json:"token"`
	Solutions    []int  `json:"solutions"`
	Instr        any    `json:"instr,omitempty"`
	InstrTimeout bool   `json:"instr_timeout,omitempty"`
	InstrBlocked bool   `json:"instr_blocked,omitempty"`
}

// RedeemResponse is returned after a successful redeem request.
type RedeemResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
	Expires int64  `json:"expires"`
}

// SiteVerifyRequest is the request payload for server-side verification.
type SiteVerifyRequest struct {
	Secret   string `json:"secret"`
	Response string `json:"response"`
}

// SiteVerifyResponse is returned by the siteverify endpoint.
type SiteVerifyResponse struct {
	Success bool `json:"success"`
}

// ChallengeClaims represents the signed challenge token payload.
type ChallengeClaims struct {
	SiteKey        string `json:"sk"`
	Nonce          string `json:"n"`
	ChallengeCount int    `json:"c"`
	SaltSize       int    `json:"s"`
	Difficulty     int    `json:"d"`
	ExpiresAtMS    int64  `json:"exp"`
	IssuedAtMS     int64  `json:"iat"`
}
