package event

import "github.com/newm4n/mihp/pkg/helper"

type ProbeData map[string]*ProbeStatistic

type ProbeStatistic struct {
	HttpCallDurationStat *Stat
	DownTimeDurationStat *Stat
	DownTimeInterval     *helper.Interval
	ErrorRecord          map[int64]string
}

func NewStat() *Stat {
	return &Stat{
		HourlyDuration:  make([]*Metric, (365*24)/2), // half year
		DailyDuration:   make([]*Metric, 365),        // whole year
		WeeklyDuration:  make([]*Metric, 52),         // whole year
		MonthlyDuration: make([]*Metric, 12),         // whole year
	}
}

type Stat struct {
	HourlyDuration  []*Metric
	DailyDuration   []*Metric
	WeeklyDuration  []*Metric
	MonthlyDuration []*Metric
}

type Metric struct {
	Top    int64
	Bottom int64
	Count  int64
	Total  int64
	Avg    int64
}

func (m *Metric) Add(value int64) {
	m.Count++
	m.Total += value
	if value > m.Top {
		m.Top = value
	} else {
		m.Bottom = value
	}
	m.Avg = m.Total / m.Count
}
