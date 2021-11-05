package internal

import (
	"encoding/binary"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestProbeContext_String(t *testing.T) {
	pb := NewProbeContext()
	pb["row1Int"] = 123
	pb["row2String"] = "One Two Three"
	pb["row3Float"] = 123.456
	pb["row4Bool"] = true
	pb["row5Time"] = time.Now()

	t.Log(pb.String())
}

func TestSerialize(t *testing.T) {
	i := int32(-543905)
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(i))

	ui := binary.BigEndian.Uint32(b)
	assert.Equal(t, i, int32(ui))
}
