package minion

import (
	"context"
	"fmt"
	"github.com/newm4n/mihp/internal"
	"github.com/newm4n/mihp/internal/probing"
	"github.com/newm4n/mihp/minion/com"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const (
	MinionServerPort = 62891
)

var (
	Config               *internal.MIHPConfig
	EmailNotifChannel    = make(map[string]*probing.ProbeEventProcessor)
	CallbackNotifChannel = make(map[string]*probing.ProbeEventProcessor)
	LogNotifChannel      = make(map[string]*probing.ProbeEventProcessor)
	Rank                 uint64
)

func init() {
	Rank = rand.Uint64()
}

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

func MinionDaemonHandler(message *com.UDPMessage) {

}

func Start(ctx context.Context) {
	go func() {
		err := com.StartServer(ctx, net.IP{0, 0, 0, 0}, MinionServerPort, MinionDaemonHandler)
		if err != nil {
			fmt.Sprintf(err.Error())
		}
	}()

	gracefulStop := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(gracefulStop, os.Interrupt)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	// Block until we receive our signal.
	<-gracefulStop

	com.StopServer()

	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	logrus.Info("shutting down minion........ bye")

}
