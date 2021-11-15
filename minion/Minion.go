package minion

import (
	"github.com/newm4n/mihp/internal"
	"github.com/newm4n/mihp/internal/probing"
)

var (
	Config               *internal.MIHPConfig
	EmailNotifChannel    = make(map[string]*probing.ProbeEventProcessor)
	CallbackNotifChannel = make(map[string]*probing.ProbeEventProcessor)
	LogNotifChannel      = make(map[string]*probing.ProbeEventProcessor)
)

func Initialize(MIHPConfig *internal.MIHPConfig) {
	Config = MIHPConfig
	for _, probe := range Config.ProbePool {
		AcceptProbe(probe)
	}
}

func AcceptProbe(probe *internal.Probe) {
	if probe.SMTPNotification != nil {
		EmailNotifChannel[probe.ID] = probing.NewProbeEventProcessor(probing.LogTrigger)
	}
	// todo finish this MINION
}
