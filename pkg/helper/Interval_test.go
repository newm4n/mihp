package helper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRange_IsIn2(t *testing.T) {
	assert.True(t, NewRange(45, 67).IsIn(50))
}

func TestMinint(t *testing.T) {
	assert.Equal(t, 1, minint(1, 2))
	assert.Equal(t, 2, minint(2, 2))
	assert.Equal(t, 2, minint(3, 2))
}

func TestMaxint(t *testing.T) {
	assert.Equal(t, 2, maxint(1, 2))
	assert.Equal(t, 2, maxint(2, 2))
	assert.Equal(t, 3, maxint(3, 2))
}

func TestRange_Touches(t *testing.T) {
	r1 := NewRange(4, 10)
	r2 := NewRange(11, 15)
	assert.True(t, r1.Touches(r2))
}

func TestRange_Combine(t *testing.T) {
	r1 := NewRange(4, 10)
	r2 := NewRange(11, 15)
	r3, err := r1.Combine(r2)
	assert.NoError(t, err)
	assert.Equal(t, 4, r3.From)
	assert.Equal(t, 15, r3.To)

	r1 = NewRange(4, 10)
	r2 = NewRange(6, 15)
	r3, err = r1.Combine(r2)
	assert.NoError(t, err)
	assert.Equal(t, 4, r3.From)
	assert.Equal(t, 15, r3.To)

	r1 = NewRange(4, 15)
	r2 = NewRange(6, 12)
	r3, err = r1.Combine(r2)
	assert.NoError(t, err)
	assert.Equal(t, 4, r3.From)
	assert.Equal(t, 15, r3.To)

	r1 = NewRange(4, 15)
	r2 = NewRange(17, 23)
	r3, err = r1.Combine(r2)
	assert.Error(t, err)
}

func TestInterval_Add(t *testing.T) {
	inter := &Interval{Ranges: make([]*Range, 0)}
	assert.Equal(t, 0, len(inter.Ranges))

	inter.Add(5)
	assert.Equal(t, 1, len(inter.Ranges))

	inter.Add(7)
	assert.Equal(t, 2, len(inter.Ranges))

	inter.Add(8)
	assert.Equal(t, 2, len(inter.Ranges))

	inter.Add(6)
	assert.Equal(t, 1, len(inter.Ranges))

	assert.False(t, inter.IsIn(4))
	assert.True(t, inter.IsIn(5))
	assert.True(t, inter.IsIn(6))
	assert.True(t, inter.IsIn(7))
	assert.True(t, inter.IsIn(8))
	assert.False(t, inter.IsIn(9))
}

func TestStringToInterval(t *testing.T) {
	itv, err := StringToInterval("1,3,5")
	assert.NoError(t, err)
	assert.Equal(t, 3, len(itv.Ranges))
	assert.True(t, itv.IsIn(1))
	assert.False(t, itv.IsIn(2))
	assert.True(t, itv.IsIn(3))
	assert.False(t, itv.IsIn(4))
	assert.True(t, itv.IsIn(5))

	itv, err = StringToInterval("1,2,3,4,5")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(itv.Ranges))
	assert.True(t, itv.IsIn(1))
	assert.True(t, itv.IsIn(2))
	assert.True(t, itv.IsIn(3))
	assert.True(t, itv.IsIn(4))
	assert.True(t, itv.IsIn(5))

	itv, err = StringToInterval("1-2,4,6-8")
	assert.NoError(t, err)
	assert.Equal(t, 3, len(itv.Ranges))
	assert.False(t, itv.IsIn(0))
	assert.True(t, itv.IsIn(1))
	assert.True(t, itv.IsIn(2))
	assert.False(t, itv.IsIn(3))
	assert.True(t, itv.IsIn(4))
	assert.False(t, itv.IsIn(5))
	assert.True(t, itv.IsIn(6))
	assert.True(t, itv.IsIn(7))
	assert.True(t, itv.IsIn(8))
	assert.False(t, itv.IsIn(9))

	itv, err = StringToInterval("1-2,4,6-8 ")
	assert.Error(t, err)
}

func TestInterval_IsIn_Step(t *testing.T) {
	i, err := StringToInterval("*/2")
	assert.Equal(t, 1, len(i.Steps))
	for s := range i.Steps {
		t.Logf("%d", s)
	}

	assert.NoError(t, err)
	assert.True(t, i.IsIn(10))
	assert.False(t, i.IsIn(11))
	assert.True(t, i.IsIn(12))
	assert.False(t, i.IsIn(13))
	assert.True(t, i.IsIn(14))
	assert.False(t, i.IsIn(15))
	assert.True(t, i.IsIn(16))
	assert.False(t, i.IsIn(17))
	assert.True(t, i.IsIn(18))
	assert.False(t, i.IsIn(19))
	assert.True(t, i.IsIn(20))
	assert.False(t, i.IsIn(21))
}
