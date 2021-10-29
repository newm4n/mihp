package cron

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestStart(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
	Start()
	schedule, err := NewSchedule("*/2 * * * * * *")
	assert.NoError(t, err)
	AddJob("2SecTicker", &Job{
		Cron:     schedule,
		Deadline: 10 * time.Second,
		JobFunc: func(ctx context.Context) {
			t.Logf("%s", time.Now().Format(time.RFC3339))
		},
	})
	time.Sleep(10 * time.Second)
	Stop()
}

func TestNewCronStruct(t *testing.T) {
	cron, err := NewSchedule("* * * * * * *")
	assert.NoError(t, err)
	assert.NotNil(t, cron)

	cron, err = NewSchedule("1,4,6,12-23    1,4,6,12-23 1,4,6,12-23  1,4,6,12-23 1,4,6,12-23 1,4,6,12-23 1,4,6,12-23,40-,-50")
	assert.NoError(t, err)
	assert.NotNil(t, cron)
}

func TestCronStruct_IsIn(t *testing.T) {
	cron, err := NewSchedule("10 20 15 14 * 4 2020")
	assert.NoError(t, err)
	assert.NotNil(t, cron)

	assert.True(t, cron.IsIn(time.Date(2020, time.April, 14, 15, 20, 10, 0, time.Local)))
	assert.False(t, cron.IsIn(time.Date(2020, time.April, 14, 15, 20, 11, 0, time.Local)))
	assert.False(t, cron.IsIn(time.Date(2020, time.April, 14, 15, 20, 9, 0, time.Local)))

	cron, err = NewSchedule("*/2 * * * * * *")
	assert.True(t, cron.IsIn(time.Date(2020, time.April, 14, 15, 20, 10, 0, time.Local)))
	assert.False(t, cron.IsIn(time.Date(2020, time.April, 14, 15, 20, 11, 0, time.Local)))
	assert.True(t, cron.IsIn(time.Date(2020, time.April, 14, 15, 20, 12, 0, time.Local)))
	assert.False(t, cron.IsIn(time.Date(2020, time.April, 14, 15, 20, 13, 0, time.Local)))
	assert.True(t, cron.IsIn(time.Date(2020, time.April, 14, 15, 20, 14, 0, time.Local)))
	assert.False(t, cron.IsIn(time.Date(2020, time.April, 14, 15, 20, 15, 0, time.Local)))
	assert.True(t, cron.IsIn(time.Date(2020, time.April, 14, 15, 20, 16, 0, time.Local)))
	assert.False(t, cron.IsIn(time.Date(2020, time.April, 14, 15, 20, 17, 0, time.Local)))
	assert.True(t, cron.IsIn(time.Date(2020, time.April, 14, 15, 20, 18, 0, time.Local)))
	assert.False(t, cron.IsIn(time.Date(2020, time.April, 14, 15, 20, 19, 0, time.Local)))
}
