package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIncrement(t *testing.T) {
	assert.NotContains(t, metricMap, "abc")

	// Does it set the key and value?
	Increment("abc", 6)
	assert.Contains(t, metricMap, "abc")
	val := metricMap["abc"]
	assert.Equal(t, val, 6)

	// Does add to an existing key?
	Increment("abc", 1)
	val = metricMap["abc"]
	assert.Equal(t, val, 7)

	// Does it add whatever we give it?
	Increment("abc", 3)
	val = metricMap["abc"]
	assert.Equal(t, val, 10)
}

func TestExport(t *testing.T) {
	Increment("fffffff", 111)

	// Does Export() return a cloned map that's equal to metricMap?
	exp := Export()
	assert.Equal(t, exp, metricMap)
	assert.NotSame(t, exp, metricMap)

	// Does the value we set earlier exist in the cloned map?
	val := exp["fffffff"]
	assert.Equal(t, val, 111)
}

func TestClear(t *testing.T) {
	Increment("abc", 1)

	Clear()
	assert.Empty(t, metricMap)
}
