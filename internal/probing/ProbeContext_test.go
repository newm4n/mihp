package probing

import (
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
