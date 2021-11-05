package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetUTCCount(t *testing.T) {
	y, m, d, w, h := GetUTCCount(time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC))
	assert.Zero(t, y)
	assert.Zero(t, m)
	assert.Zero(t, d)
	assert.Zero(t, w)
	assert.Zero(t, h)

	y, m, d, w, h = GetUTCCount(time.Date(1970, time.January, 1, 1, 0, 0, 0, time.UTC))
	assert.Zero(t, y)
	assert.Zero(t, m)
	assert.Zero(t, d)
	assert.Zero(t, w)
	assert.Equal(t, int64(1), h)

	y, m, d, w, h = GetUTCCount(time.Date(1970, time.January, 2, 0, 0, 0, 0, time.UTC))
	assert.Zero(t, y)
	assert.Zero(t, m)
	assert.Equal(t, int64(1), d)
	assert.Zero(t, w)
	assert.Equal(t, int64(24), h)

	y, m, d, w, h = GetUTCCount(time.Date(1970, time.January, 3, 1, 0, 0, 0, time.UTC))
	assert.Zero(t, y)
	assert.Zero(t, m)
	assert.Equal(t, int64(2), d)
	assert.Zero(t, w)
	assert.Equal(t, int64(49), h)

	y, m, d, w, h = GetUTCCount(time.Date(1970, time.February, 1, 0, 0, 0, 0, time.UTC))
	assert.Zero(t, y)
	assert.Equal(t, int64(1), m)
	assert.Equal(t, int64(31), d)
	assert.Equal(t, int64(4), w)
	assert.Equal(t, int64(744), h)

	y, m, d, w, h = GetUTCCount(time.Date(1971, time.January, 1, 0, 0, 0, 0, time.UTC))
	assert.Equal(t, int64(1), y)
	assert.Equal(t, int64(12), m)
	assert.Equal(t, int64(365), d)
	assert.Equal(t, int64(52), w)
	assert.Equal(t, int64(8760), h)
}
