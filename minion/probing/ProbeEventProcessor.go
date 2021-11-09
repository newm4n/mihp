package probing

import (
	"container/list"
	"fmt"
	"github.com/newm4n/mihp/internal"
	"log"
	"time"
)

type Trigger func(probeName, probeId string, down bool, firstUp, lastUp, firstDown, lastDown time.Time)

func LogTrigger(probeName, probeId string, down bool, firstUp, lastUp, firstDown, lastDown time.Time) {
	if !down {
		if lastUp.Sub(firstDown) > 24*365*time.Hour {
			log.Printf("Probe %s [%s] is detected online for the first time this year at %s", probeName, probeId, firstUp.Format(time.RFC3339))
		} else {
			log.Printf("Probe %s [%s] is back up at %s after downed for %s", probeName, probeId, firstUp.Format(time.RFC3339), firstUp.Sub(firstDown).String())
		}
	} else {
		if lastDown.Sub(firstUp) > 24*365*time.Hour {
			log.Printf("Probe %s [%s] is downed for the first time this year at %s", probeName, probeId, firstDown.Format(time.RFC3339))
		} else {
			log.Printf("Probe %s [%s] is downed at %s after on-line for %s", probeName, probeId, firstDown.Format(time.RFC3339), firstDown.Sub(firstUp).String())
		}
	}
}

func NewProbeEventProcessor(trigger Trigger) *ProbeEventProcessor {
	if trigger != nil {
		return &ProbeEventProcessor{Trigger: trigger}
	}
	return &ProbeEventProcessor{Trigger: LogTrigger}
}

type ProbeEventProcessor struct {
	Trackers []*ProbeEventTracker
	Trigger  Trigger
}

func (proc *ProbeEventProcessor) AcceptProbeContext(pbctx internal.ProbeContext) *ProbeEventTracker {
	if proc.Trackers == nil {
		proc.Trackers = make([]*ProbeEventTracker, 0)
	}
	for _, t := range proc.Trackers {
		if t.ProbeName == pbctx["probe"].(string) {
			t.AcceptProbeContext(pbctx, proc.Trigger)
			return t
		}
	}
	name := pbctx["probe"].(string)
	id := pbctx[fmt.Sprintf("probe.%s.id", name)].(string)
	t := &ProbeEventTracker{
		ProbeID:          id,
		ProbeName:        name,
		FailThreshold:    2,
		SuccessThreshold: 2,
		FailCount:        0,
		SuccessCount:     0,
		LastDown:         time.UnixMilli(0),
		LastUp:           time.UnixMilli(0),
		LastStatusDown:   true,
	}
	t.AcceptProbeContext(pbctx, proc.Trigger)
	proc.Trackers = append(proc.Trackers, t)
	return t
}

type ProbeEventTracker struct {
	ProbeID          string
	ProbeName        string
	FailThreshold    int
	SuccessThreshold int

	FailCount    int
	SuccessCount int

	FirstUp   time.Time
	LastUp    time.Time
	FirstDown time.Time
	LastDown  time.Time

	LastStatusDown bool
	UpDownHistory  *list.List
}

type History struct {
	Down bool
	Time time.Time
}

func (t *ProbeEventTracker) AcceptProbeContext(pbctx internal.ProbeContext, trigger Trigger) {
	if t.UpDownHistory == nil {
		t.UpDownHistory = list.New().Init()
	}

	name := pbctx["probe"].(string)
	id := pbctx[fmt.Sprintf("probe.%s.id", name)].(string)

	if name != t.ProbeName {
		return
	}
	if pbctx[fmt.Sprintf("probe.%s.success", name)].(bool) {
		t.FailCount = 0
		t.SuccessCount++
		t.UpDownHistory.PushBack(&History{
			Down: false,
			Time: pbctx[fmt.Sprintf("probe.%s.starttime", t.ProbeName)].(time.Time),
		})
	} else {
		t.FailCount++
		t.SuccessCount = 0
		t.UpDownHistory.PushBack(&History{
			Down: true,
			Time: pbctx[fmt.Sprintf("probe.%s.starttime", t.ProbeName)].(time.Time),
		})
	}

	if t.UpDownHistory.Len() > 30 {
		t.UpDownHistory.Remove(t.UpDownHistory.Front())
	}

	if t.FailCount > t.FailThreshold && !t.LastStatusDown {
		t.LastStatusDown = true

		if t.UpDownHistory.Len() < t.FailThreshold {
			t.FirstDown = t.UpDownHistory.Front().Value.(*History).Time
		} else {
			ele := t.UpDownHistory.Back()
			for i := 0; i < t.FailThreshold-1; i++ {
				ele = ele.Prev()
			}
			t.FirstDown = ele.Value.(*History).Time
			t.LastUp = ele.Prev().Value.(*History).Time
		}

		trigger(name, id, true, t.FirstUp, t.LastUp, t.FirstDown, t.LastDown)
	} else if t.SuccessCount > t.SuccessThreshold && t.LastStatusDown {
		t.LastStatusDown = false

		if t.UpDownHistory.Len() < t.SuccessThreshold {
			t.FirstUp = t.UpDownHistory.Front().Value.(*History).Time
		} else {
			ele := t.UpDownHistory.Back()
			for i := 0; i < t.SuccessThreshold-1; i++ {
				ele = ele.Prev()
			}
			t.FirstUp = ele.Value.(*History).Time
			t.LastDown = ele.Prev().Value.(*History).Time
		}

		trigger(name, id, false, t.FirstUp, t.LastUp, t.FirstDown, t.LastDown)
	}
}
