package gaad

import (
	"sync"
	"testing"
)

// TestWindowGroupingRaceCondition tests that a data race does not occur
// when run with the -race flag.
func TestWindowGroupingRaceCondition(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			info := &ics_info{}
			info.Window_sequence = ONLY_LONG_SEQUENCE
			window_grouping(info, 0, 960)

		}()
	}
	wg.Wait()
}

func TestRaceConditionEightShort(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			info := &ics_info{}
			info.Window_sequence = EIGHT_SHORT_SEQUENCE
			window_grouping(info, 0, 960)

		}()
	}
	wg.Wait()
}
