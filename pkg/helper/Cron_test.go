package helper

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewCronStruct(t *testing.T) {
	cron, err := NewCronStruct("* * * * * * *")
	assert.NoError(t, err)
	assert.NotNil(t, cron)

	cron, err = NewCronStruct("1,4,6,12-23    1,4,6,12-23 1,4,6,12-23  1,4,6,12-23 1,4,6,12-23 1,4,6,12-23 1,4,6,12-23,40-,-50")
	assert.NoError(t, err)
	assert.NotNil(t, cron)
}

func TestCronStruct_IsIn(t *testing.T) {
	cron, err := NewCronStruct("10 20 15 14 * 4 2020")
	assert.NoError(t, err)
	assert.NotNil(t, cron)

	assert.True(t, cron.IsIn(time.Date(2020, time.April, 14, 15, 20, 10, 0, time.Local)))
	assert.False(t, cron.IsIn(time.Date(2020, time.April, 14, 15, 20, 11, 0, time.Local)))
	assert.False(t, cron.IsIn(time.Date(2020, time.April, 14, 15, 20, 9, 0, time.Local)))
}
