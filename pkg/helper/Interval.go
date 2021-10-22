package helper

import (
	"fmt"
	"github.com/newm4n/mihp/pkg/errors"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type Range struct {
	From int
	To   int
}

func (r *Range) IsIn(val int) bool {
	return val >= r.From && val <= r.To
}

func (r *Range) Touches(that *Range) bool {
	return r.From == that.To+1 || r.To == that.From-1
}

func (r *Range) Overlaps(that *Range) bool {
	return that.IsIn(r.From) || that.IsIn(r.To) || r.IsIn(that.From) || r.IsIn(that.To) || (that.IsIn(r.From) && that.IsIn(r.To)) || (r.IsIn(that.From) && r.IsIn(that.To))
}

func minint(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxint(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (r *Range) Combine(that *Range) (*Range, error) {
	if r.Touches(that) || r.Overlaps(that) {
		return &Range{
			From: minint(r.From, that.From),
			To:   maxint(r.To, that.To),
		}, nil
	}
	return nil, fmt.Errorf("ranges not touching nor overlaps")
}

type Interval struct {
	Ranges []*Range
}

func (i *Interval) IsIn(val int) bool {
	for _, r := range i.Ranges {
		if r.IsIn(val) {
			return true
		}
	}
	return false
}

func (i *Interval) Add(val int) {
	i.AddRange(val, val)
}

func (i *Interval) AddRange(a, b int) {
	var r *Range
	if a < b {
		r = &Range{
			From: a,
			To:   b,
		}
	} else {
		r = &Range{
			From: b,
			To:   a,
		}
	}
	nRange := make([]*Range, 0)
	for _, er := range i.Ranges {
		merges, err := r.Combine(er)
		if err != nil {
			nRange = append(nRange, er)
		} else {
			r = merges
		}
	}
	nRange = append(nRange, r)
	i.Ranges = nRange
}

func StringToInterval(seg string) (*Interval, error) {
	itrv := &Interval{Ranges: make([]*Range, 0)}
	if strings.ContainsAny(seg, " \t\n\r") {
		return nil, fmt.Errorf("%w : %s", errors.ErrInvalidCronExpression, seg)
	}
	toks := strings.Split(seg, ",")
	numOnly := regexp.MustCompile(`^[0-9]+$`)
	anyBigger := regexp.MustCompile(`^[0-9]+\-$`)
	anySmaller := regexp.MustCompile(`^\-[0-9]+$`)
	rangeIn := regexp.MustCompile(`^[0-9]+\-[0-9]+$`)
	for _, t := range toks {
		if strings.ContainsAny(t, " \t\n\r") {
			return nil, fmt.Errorf("%w : %s", errors.ErrInvalidCronExpression, seg)
		}
		if t == "*" {
			itrv.AddRange(math.MinInt32, math.MaxInt32)
		} else if numOnly.MatchString(t) {
			if num, err := strconv.Atoi(t); err == nil {
				itrv.Add(num)
				continue
			} else {
				return nil, fmt.Errorf("%w : %s", errors.ErrInvalidCronExpression, seg)
			}
		} else if anyBigger.MatchString(t) {
			if num, err := strconv.Atoi(t[:len(t)-2]); err == nil {
				itrv.AddRange(num, math.MaxInt32)
				continue
			} else {
				return nil, fmt.Errorf("%w : %s", errors.ErrInvalidCronExpression, seg)
			}
		} else if anySmaller.MatchString(t) {
			if num, err := strconv.Atoi(t[1:]); err == nil {
				itrv.AddRange(math.MinInt32, num)
				continue
			} else {
				return nil, fmt.Errorf("%w : %s", errors.ErrInvalidCronExpression, seg)
			}
		} else if rangeIn.MatchString(t) {
			idx := strings.Index(t, "-")
			from, ferr := strconv.Atoi(t[:idx])
			to, terr := strconv.Atoi(t[idx+1:])
			if ferr != nil || terr != nil {
				return nil, fmt.Errorf("%w : %s", errors.ErrInvalidCronExpression, seg)
			} else {
				itrv.AddRange(from, to)
			}
		} else {
			return nil, fmt.Errorf("%w : %s", errors.ErrInvalidCronExpression, seg)
		}
	}
	return itrv, nil
}
