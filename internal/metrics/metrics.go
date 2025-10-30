package metrics

import (
	"fmt"
	"os"
	"slices"
	"sync"

	"golang.org/x/exp/maps"
)

var (
	metricMap   = map[string]int{}
	metricMutex sync.Mutex
)

func Increment(key string, num int) {
	metricMutex.Lock()
	defer metricMutex.Unlock()

	orig, ok := metricMap[key]
	if ok {
		// The key exists, so add to it
		metricMap[key] = orig + num
	} else {
		metricMap[key] = num
	}
}

func Export() map[string]int {
	expMap := map[string]int{}

	// Using a shared mutex means we won't have any keys update while we
	// export
	metricMutex.Lock()
	defer metricMutex.Unlock()

	for key, val := range metricMap {
		expMap[key] = val
	}

	return expMap
}

func Clear() {
	metricMutex.Lock()

	metricMap = map[string]int{}

	metricMutex.Unlock()
}

func WriteToFile(file string) error {
	metricMutex.Lock()
	defer metricMutex.Unlock()

	keys := maps.Keys(metricMap)
	slices.Sort(keys)

	fp, err := os.Create(file)
	if err != nil {
		return err
	}

	defer fp.Close() //nolint:errcheck

	for _, key := range keys {
		line := fmt.Sprintf("%v = %v\n", key, metricMap[key])

		if _, err := fp.WriteString(line); err != nil {
			return err
		}
	}

	return nil
}
