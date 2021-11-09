package probing

import (
	"fmt"
	"github.com/newm4n/mihp/internal"
	"testing"
	"time"
)

func DummyContext(name, id string, probeTime time.Time, span time.Duration, success []bool) []internal.ProbeContext {
	contexts := make([]internal.ProbeContext, len(success))
	for i := 0; i < len(success); i++ {
		dur := span * time.Duration(i)
		t := probeTime.Add(dur)
		contexts[i] = internal.NewProbeContext()
		contexts[i]["probe"] = name
		contexts[i][fmt.Sprintf("probe.%s.id", name)] = id
		contexts[i][fmt.Sprintf("probe.%s.fail", name)] = !success[i]
		contexts[i][fmt.Sprintf("probe.%s.success", name)] = success[i]
		contexts[i][fmt.Sprintf("probe.%s.starttime", name)] = t
		contexts[i][fmt.Sprintf("probe.%s.endtime", name)] = t.Add(2 * time.Second)
		contexts[i][fmt.Sprintf("probe.%s.duration", name)] = 2 * time.Second
	}
	return contexts
}

func TestNewProbeEventProcessor(t *testing.T) {
	eventProc := NewProbeEventProcessor(nil)
	dummy := DummyContext("dummy", "123456789", time.Now(), 1*time.Second,
		[]bool{true, true, true, true, false, false, false, false, true, true, true,
			true, false, false, true, false, false, false, true, false, true, false, true, true, true, true,
			false, false, false, false, false, false, false, true, false, true, false, true, true, true})
	for idx, ctx := range dummy {
		t.Logf("----- %d  ---- %v", idx, ctx["probe.dummy.success"])
		eventProc.AcceptProbeContext(ctx)
		time.Sleep(1 * time.Second)
	}
}
