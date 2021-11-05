package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGraph_GetGraph(t *testing.T) {
	g := make(Graph)
	for i := int64(100); i >= 0; i-- {
		g[i] = &GraphUnit{
			Idx: i,
			Max: i,
			Min: i,
			Avg: i,
			Tot: i,
			Cnt: 1,
		}
	}
	gs := g.GetGraph(30, 40)
	assert.Equal(t, 11, len(gs))
	assert.Equal(t, int64(30), gs[0].Idx)
	assert.Equal(t, int64(31), gs[1].Idx)
	assert.Equal(t, int64(30), gs[0].Avg)
	assert.Equal(t, int64(31), gs[1].Avg)

	gs = g.GetGraph(130, 140)
	assert.Equal(t, 11, len(gs))
	assert.Equal(t, int64(130), gs[0].Idx)
	assert.Equal(t, int64(131), gs[1].Idx)
	assert.Equal(t, int64(0), gs[0].Avg)
	assert.Equal(t, int64(0), gs[1].Avg)

}
