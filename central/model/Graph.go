package model

import "sort"

type Graph map[int64]*GraphUnit

func (g Graph) GetGraph(from, to int64) []*GraphUnit {
	if to < from {
		tmp := from
		from = to
		to = tmp
	}
	// lets limit 20 graph length
	if to-from > 50 {
		to = from + 50
	}
	ret := make([]*GraphUnit, 0)

	for i := from; i <= to; i++ {
		if g, ok := g[i]; ok {
			ret = append(ret, g)
		} else {
			ret = append(ret, &GraphUnit{
				Idx: i,
			})
		}
	}

	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Idx < ret[j].Idx
	})

	return ret
}

func NewProbeStatistic(idx, initVal int64) *GraphUnit {
	return &GraphUnit{
		Idx: idx,
		Max: initVal,
		Min: initVal,
		Avg: initVal,
		Tot: initVal,
		Cnt: 1,
	}
}

type GraphUnit struct {
	Idx int64
	Max int64
	Min int64
	Avg int64
	Tot int64
	Cnt int64
}

func (gu *GraphUnit) Add(num int64) {
	gu.Cnt += 1
	gu.Tot += num
	if num > gu.Max {
		gu.Max = num
	}
	if num < gu.Min {
		gu.Min = num
	}
	gu.Avg = gu.Tot / gu.Cnt
}
