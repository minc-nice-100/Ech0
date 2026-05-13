// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025-2026 lin-snow

package core

import "net/http"

// ErrorCode is a machine-readable error category used in API responses.
type ErrorCode string

const (
	// ErrCodeBadRequest indicates request payload/shape problems.
	ErrCodeBadRequest ErrorCode = "bad_request"
	// ErrCodeForbidden indicates authentication/authorization or semantic rejection.
	ErrCodeForbidden ErrorCode = "forbidden"
	// ErrCodeNotFound indicates missing resources.
	ErrCodeNotFound ErrorCode = "not_found"
	// ErrCodeRateLimit indicates request throttling.
	ErrCodeRateLimit ErrorCode = "rate_limit"
	// ErrCodeInternal indicates unexpected server-side failures.
	ErrCodeInternal ErrorCode = "internal"
)

// Error is the domain-level error used by handlers for status mapping.
type Error struct {
	Code       ErrorCode
	Message    string
	StatusCode int
}

func (e *Error) Error() string {
	return e.Message
}

// NewBadRequest creates a bad-request domain error.
func NewBadRequest(msg string) *Error {
	return &Error{Code: ErrCodeBadRequest, Message: msg, StatusCode: http.StatusBadRequest}
}

// NewForbidden creates a forbidden domain error.
func NewForbidden(msg string) *Error {
	return &Error{Code: ErrCodeForbidden, Message: msg, StatusCode: http.StatusForbidden}
}

// NewNotFound creates a not-found domain error.
func NewNotFound(msg string) *Error {
	return &Error{Code: ErrCodeNotFound, Message: msg, StatusCode: http.StatusNotFound}
}

// NewRateLimit creates a rate-limit domain error.
func NewRateLimit(msg string) *Error {
	return &Error{Code: ErrCodeRateLimit, Message: msg, StatusCode: http.StatusTooManyRequests}
}

// NewInternal creates an internal domain error.
func NewInternal(msg string) *Error {
	return &Error{Code: ErrCodeInternal, Message: msg, StatusCode: http.StatusInternalServerError}
}
