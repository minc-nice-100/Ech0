package memstore

import (
	"testing"
	"time"
)

func TestConsumeRedeemTokenOnce(t *testing.T) {
	st := New(Options{GCInterval: time.Hour})
	defer func() {
		_ = st.Close()
	}()

	now := time.Now()
	if err := st.StoreRedeemToken("site", "tok", now.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}

	found, expired, err := st.ConsumeRedeemToken("site", "tok", now)
	if err != nil {
		t.Fatal(err)
	}
	if !found || expired {
		t.Fatalf("expected found and not expired")
	}

	found, expired, err = st.ConsumeRedeemToken("site", "tok", now)
	if err != nil {
		t.Fatal(err)
	}
	if found || expired {
		t.Fatalf("expected token to be consumed")
	}
}

func TestRateLimitWindow(t *testing.T) {
	st := New(Options{GCInterval: time.Hour})
	defer func() {
		_ = st.Close()
	}()

	now := time.Unix(1000, 0)
	ok, _, err := st.AllowRateLimit("cap", "1.1.1.1", 2, time.Second, now)
	if err != nil || !ok {
		t.Fatalf("first request should pass")
	}
	ok, _, err = st.AllowRateLimit("cap", "1.1.1.1", 2, time.Second, now)
	if err != nil || !ok {
		t.Fatalf("second request should pass")
	}
	ok, _, err = st.AllowRateLimit("cap", "1.1.1.1", 2, time.Second, now)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatalf("third request should be limited")
	}
}
