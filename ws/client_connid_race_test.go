package ws

import (
	"fmt"
	"sync"
	"testing"
)

// Fails under -race without the connIDMu guard.
func TestFmtLogConnIDNoRace(t *testing.T) {
	c := &Client{}

	const goroutines = 4
	const iterations = 2000
	var wg sync.WaitGroup

	for w := 0; w < goroutines; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				c.setConnID(fmt.Sprintf("conn-%d-%d", w, i))
				c.setConnID("")
			}
		}(w)
	}

	for r := 0; r < goroutines; r++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				if out := c.fmtLog("probe %d", i); len(out) != 1 && len(out) != 2 {
					t.Errorf("fmtLog returned %d parts, want 1 or 2", len(out))
					return
				}
			}
		}()
	}

	wg.Wait()
}
