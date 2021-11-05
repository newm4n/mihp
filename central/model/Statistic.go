package model

import (
	"github.com/newm4n/mihp/pkg/helper"
	"time"
)

type ProbeStatistic struct {
	HourlyResponseTime  Graph
	DailyResponseTime   Graph
	WeeklyResponseTime  Graph
	MonthlyResponseTime Graph
	YearlyResponseTime  Graph
	DailyDownInterval   helper.Interval
}

func GetUTCCount(t time.Time) (y int64, m int64, d int64, w int64, h int64) {
	UTC := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
	dur := t.Sub(UTC)

	y = int64(t.Year() - 1970)
	m = (y * 12) + int64(t.Month()-1)
	d = int64(dur / (time.Hour * 24))
	h = int64(dur / time.Hour)
	w = int64(dur / (time.Hour * 24 * 7))

	return y, m, d, w, h
}

func (os *ProbeStatistic) Add(t time.Time, num int64) {
	y, m, d, w, h := GetUTCCount(t)

	if g, ok := os.YearlyResponseTime[y]; !ok {
		ngu := NewProbeStatistic(y, num)
		os.YearlyResponseTime[y] = ngu
	} else {
		g.Add(num)
	}

	if g, ok := os.MonthlyResponseTime[m]; !ok {
		ngu := NewProbeStatistic(m, num)
		os.MonthlyResponseTime[m] = ngu
	} else {
		g.Add(num)
	}

	if g, ok := os.DailyResponseTime[d]; !ok {
		ngu := NewProbeStatistic(d, num)
		os.DailyResponseTime[d] = ngu
	} else {
		g.Add(num)
	}

	if g, ok := os.WeeklyResponseTime[w]; !ok {
		ngu := NewProbeStatistic(w, num)
		os.WeeklyResponseTime[w] = ngu
	} else {
		g.Add(num)
	}

	if g, ok := os.HourlyResponseTime[h]; !ok {
		ngu := NewProbeStatistic(h, num)
		os.HourlyResponseTime[h] = ngu
	} else {
		g.Add(num)
	}

}
