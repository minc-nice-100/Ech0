package memstore

import "time"

func (s *Store) startGC(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				s.gcOnce(time.Now())
			case <-s.stopCh:
				return
			}
		}
	}()
}

func (s *Store) gcOnce(now time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for sig, exp := range s.usedChallengeSig {
		if !exp.After(now) {
			delete(s.usedChallengeSig, sig)
		}
	}

	for k, entry := range s.redeemTokens {
		if !entry.expiresAt.After(now) {
			delete(s.redeemTokens, k)
		}
	}

	for k, w := range s.rateWindows {
		if !w.expiresAt.After(now) {
			delete(s.rateWindows, k)
		}
	}
}
