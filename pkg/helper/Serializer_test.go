package helper

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestPut(t *testing.T) {
	a := byte(53)
	b := int(436232)
	c := int8(62)
	d := int16(634)
	e := int32(23654547)
	f := int64(330459609856)
	g := uint(62247547)
	h := uint8(62)
	i := uint16(23455)
	j := uint32(454564345)
	k := uint64(983475948725928)
	l := "jshefjk sehjf haskefj asef"
	m := time.Date(1977, time.August, 2, 12, 15, 19, 0, time.Local)
	n := strings.Split("quick,brown,fox,jumps", ",")

	buff := &bytes.Buffer{}
	assert.NoError(t, Put(buff, a))
	assert.NoError(t, Put(buff, b))
	assert.NoError(t, Put(buff, c))
	assert.NoError(t, Put(buff, d))
	assert.NoError(t, Put(buff, e))
	assert.NoError(t, Put(buff, f))
	assert.NoError(t, Put(buff, g))
	assert.NoError(t, Put(buff, h))
	assert.NoError(t, Put(buff, i))
	assert.NoError(t, Put(buff, j))
	assert.NoError(t, Put(buff, k))
	assert.NoError(t, Put(buff, l))
	assert.NoError(t, Put(buff, m))
	assert.NoError(t, Put(buff, n))

	buff2 := bytes.NewBuffer(buff.Bytes())
	if aa, err := Read(buff2); err != nil {
		t.Error(err.Error())
	} else {
		assert.Equal(t, a, aa.(byte))
	}

	if bb, err := Read(buff2); err != nil {
		t.Error(err.Error())
	} else {
		assert.Equal(t, b, bb.(int))
	}
	if cc, err := Read(buff2); err != nil {
		t.Error(err.Error())
	} else {
		assert.Equal(t, c, cc.(int8))
	}
	if dd, err := Read(buff2); err != nil {
		t.Error(err.Error())
	} else {
		assert.Equal(t, d, dd.(int16))
	}
	if ee, err := Read(buff2); err != nil {
		t.Error(err.Error())
	} else {
		assert.Equal(t, e, ee.(int32))
	}
	if ff, err := Read(buff2); err != nil {
		t.Error(err.Error())
	} else {
		assert.Equal(t, f, ff.(int64))
	}
	if gg, err := Read(buff2); err != nil {
		t.Error(err.Error())
	} else {
		assert.Equal(t, g, gg.(uint))
	}
	if hh, err := Read(buff2); err != nil {
		t.Error(err.Error())
	} else {
		assert.Equal(t, h, hh.(uint8))
	}
	if ii, err := Read(buff2); err != nil {
		t.Error(err.Error())
	} else {
		assert.Equal(t, i, ii.(uint16))
	}
	if jj, err := Read(buff2); err != nil {
		t.Error(err.Error())
	} else {
		assert.Equal(t, j, jj.(uint32))
	}
	if kk, err := Read(buff2); err != nil {
		t.Error(err.Error())
	} else {
		assert.Equal(t, k, kk.(uint64))
	}
	if ll, err := Read(buff2); err != nil {
		t.Error(err.Error())
	} else {
		assert.Equal(t, l, ll.(string))
	}
	if mm, err := Read(buff2); err != nil {
		t.Error(err.Error())
	} else {
		assert.Equal(t, m, mm.(time.Time))
	}
	if nn, err := Read(buff2); err != nil {
		t.Error(err.Error())
	} else {
		sArr := nn.([]string)
		assert.Equal(t, len(n), len(sArr))
		for i := 0; i < len(sArr); i++ {
			assert.Equal(t, n[i], sArr[i])
		}
	}
}

func TestPutTime(t *testing.T) {
	buff := &bytes.Buffer{}
	time := time.Now()
	assert.NoError(t, PutTime(buff, time))

	nBuff := bytes.NewBuffer(buff.Bytes())
	e, err := ReadTime(nBuff)
	assert.NoError(t, err)
	t.Log(e)
}

func TestPutStringArray(t *testing.T) {
	buff := &bytes.Buffer{}
	sarr := []string{"one", "two", "three"}
	assert.NoError(t, PutStringArray(buff, sarr))

	nBuff := bytes.NewBuffer(buff.Bytes())
	e, err := ReadStringArray(nBuff)
	assert.NoError(t, err)
	assert.Equal(t, len(sarr), len(e))
}
