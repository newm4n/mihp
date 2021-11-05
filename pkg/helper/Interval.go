package helper

import (
	"fmt"
	"github.com/newm4n/mihp/pkg/errors"
	"math"
	"regexp"
	"strconv"
	"strings"
)

func NewRange(from, to int64) *Range {
	if from > to {
		return &Range{From: to, To: from}
	}
	return &Range{From: from, To: to}
}

type Range struct {
	From int64
	To   int64
}

func (r *Range) String() string {
	if r.From != r.To {
		return fmt.Sprintf("%d:%d", r.From, r.To)
	}
	return fmt.Sprintf("%d", r.From)
}

func (r *Range) IsIn(val int64) bool {
	return val >= r.From && val <= r.To
}

func (r *Range) Touches(that *Range) bool {
	return r.From == that.To+1 || r.To == that.From-1
}

func (r *Range) Overlaps(that *Range) bool {
	return that.IsIn(r.From) || that.IsIn(r.To) || r.IsIn(that.From) || r.IsIn(that.To) || (that.IsIn(r.From) && that.IsIn(r.To)) || (r.IsIn(that.From) && r.IsIn(that.To))
}

func minint(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func maxint(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func (r *Range) Combine(that *Range) (*Range, error) {
	if r.Touches(that) || r.Overlaps(that) {
		return NewRange(minint(r.From, that.From), maxint(r.To, that.To)), nil
	}
	return nil, fmt.Errorf("ranges not touching nor overlaps")
}

type Interval struct {
	Ranges []*Range
	Steps  map[int64]bool
}

func (i *Interval) String() string {
	rStrings := make([]string, len(i.Ranges))
	for idx, r := range i.Ranges {
		rStrings[idx] = r.String()
	}
	rSteps := make([]string, 0)
	for idx := range i.Steps {
		rSteps = append(rSteps, fmt.Sprintf("%d", idx))
	}
	return fmt.Sprintf("Ranges:%s Steps:%s", strings.Join(rStrings, ","), strings.Join(rSteps, ","))
}

func (i *Interval) IsIn(val int64) bool {
	if i.Steps != nil {
		for r := range i.Steps {
			if val == 0 || r == 0 {
				return true
			}
			if val%r == 0 {
				return true
			}
		}
	}
	for _, r := range i.Ranges {
		if r.IsIn(val) {
			return true
		}
	}
	return false
}

func (i *Interval) Add(val int64) {
	i.AddRange(val, val)
}

func (i *Interval) AddRange(a, b int64) {
	r := NewRange(a, b)
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
	stepNum := regexp.MustCompile(`^\*/[1-9][0-9]*$`)
	for _, t := range toks {
		if strings.ContainsAny(t, " \t\n\r") {
			return nil, fmt.Errorf("%w : %s", errors.ErrInvalidCronExpression, seg)
		}
		if t == "*" {
			itrv.AddRange(math.MinInt32, math.MaxInt32)
		} else if numOnly.MatchString(t) {
			if num, err := strconv.Atoi(t); err == nil {
				itrv.Add(int64(num))
				continue
			} else {
				return nil, fmt.Errorf("%w : %s", errors.ErrInvalidCronExpression, seg)
			}
		} else if anyBigger.MatchString(t) {
			if num, err := strconv.Atoi(t[:len(t)-2]); err == nil {
				itrv.AddRange(int64(num), math.MaxInt32)
				continue
			} else {
				return nil, fmt.Errorf("%w : %s", errors.ErrInvalidCronExpression, seg)
			}
		} else if anySmaller.MatchString(t) {
			if num, err := strconv.Atoi(t[1:]); err == nil {
				itrv.AddRange(math.MinInt32, int64(num))
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
				itrv.AddRange(int64(from), int64(to))
			}
		} else if stepNum.MatchString(t) {
			if stepInt, err := strconv.Atoi(t[2:]); err == nil {
				if itrv.Steps == nil {
					itrv.Steps = make(map[int64]bool)
				}
				itrv.Steps[int64(stepInt)] = true
			} else {
				return nil, fmt.Errorf("%w : %s", errors.ErrInvalidCronExpression, seg)
			}
		} else {
			return nil, fmt.Errorf("%w : %s", errors.ErrInvalidCronExpression, seg)
		}
	}
	return itrv, nil
}
