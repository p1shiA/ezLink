package domain

import "errors"

var (
	// User errors
	ErrUserNotFound = errors.New("user not found")
	ErrUserBanned   = errors.New("user is banned")

	// Plan errors
	ErrPlanNotFound         = errors.New("plan not found")
	ErrExtraPackageNotFound = errors.New("extra traffic package not found")

	// Subscription errors
	ErrNoActiveSubscription    = errors.New("no active subscription found")
	ErrSubscriptionExpired     = errors.New("subscription has expired")
	ErrInsufficientQuota       = errors.New("insufficient traffic quota")
	ErrOnlyPaidCanAddTraffic   = errors.New("extra traffic can only be added to an active paid subscription")

	// Transaction errors
	ErrTransactionNotFound         = errors.New("transaction not found")
	ErrTransactionAlreadyProcessed = errors.New("transaction already completed or failed")
	ErrPaymentVerificationFailed   = errors.New("payment verification failed")

	// Admin errors
	ErrAdminNotFound     = errors.New("admin not found")
	ErrAdminUnauthorized = errors.New("unauthorized: admin privileges required")
)
