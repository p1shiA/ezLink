package domain_test

import (
	"testing"
	"time"

	"github.com/p1shiA/ezLink/internal/core/domain"
)

func TestSubscription_RemainingTrafficBytes_Free(t *testing.T) {
	sub := &domain.Subscription{
		PlanType:          domain.PlanTypeFree,
		TrafficLimitBytes: 2 * 1024 * 1024 * 1024,
		TrafficUsedBytes:  500 * 1024 * 1024,
	}

	expected := int64((2*1024 - 500) * 1024 * 1024)
	if sub.RemainingTrafficBytes() != expected {
		t.Fatalf("expected %d, got %d", expected, sub.RemainingTrafficBytes())
	}

	// Over limit returns 0
	sub.TrafficUsedBytes = 3 * 1024 * 1024 * 1024
	if sub.RemainingTrafficBytes() != 0 {
		t.Fatalf("expected 0, got %d", sub.RemainingTrafficBytes())
	}
}

func TestSubscription_RemainingTrafficBytes_PaidWithExtra(t *testing.T) {
	sub := &domain.Subscription{
		PlanType:               domain.PlanTypePaid,
		TrafficLimitBytes:      10 * 1024 * 1024 * 1024,
		ExtraTrafficLimitBytes: 5 * 1024 * 1024 * 1024,
		TrafficUsedBytes:       4 * 1024 * 1024 * 1024,
	}

	// (10 + 5) - 4 = 11 GB
	expected := int64(11 * 1024 * 1024 * 1024)
	if sub.RemainingTrafficBytes() != expected {
		t.Fatalf("expected %d, got %d", expected, sub.RemainingTrafficBytes())
	}
}

func TestSubscription_IsExpired(t *testing.T) {
	subPast := &domain.Subscription{
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
	if !subPast.IsExpired() {
		t.Fatal("expected expired subscription to be expired")
	}

	subFuture := &domain.Subscription{
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}
	if subFuture.IsExpired() {
		t.Fatal("expected future subscription not to be expired")
	}
}
