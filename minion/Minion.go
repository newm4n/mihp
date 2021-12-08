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
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	MinionUDPServerPort = 62891
	MinionUDPClientPort = 62892
)

var (
	Config               *internal.MIHPConfig
	EmailNotifChannel    = make(map[string]*probing.ProbeEventProcessor)
	CallbackNotifChannel = make(map[string]*probing.ProbeEventProcessor)
	LogNotifChannel      = make(map[string]*probing.ProbeEventProcessor)
	Rank                 uint64
	MyIP                 net.IP
	LeaderIP             net.IP
	LeaderRank           uint64
	VoteTickDuration     = 10 * time.Second // should be arround every 5 minutes ?
)

func init() {
	Rank = rand.Uint64()
	LeaderRank = Rank
	MyIP = com.GetOutboundIP()
	LeaderIP = net.IP{MyIP[0], MyIP[1], MyIP[2], MyIP[3]}
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
	go func() {
		fromIP := message.FromAddr.IP
		if message.Message == "VREQ" {
			err := com.SendUDPMessage(MyIP, MinionUDPClientPort, fromIP, MinionUDPServerPort, fmt.Sprintf("VRES %d", Rank))
			if err != nil {
				logrus.Errorf("error while sending vote response to %s. got %s", fromIP.String(), err.Error())
			}
		} else if strings.HasPrefix(message.Message, "VRES ") {
			theirRankStr := strings.TrimSpace(message.Message[5:])
			theirRank, err := strconv.ParseUint(theirRankStr, 10, 64)
			if err != nil {
				logrus.Errorf("error while receiving vote response from %s. got invalid rank number format %s", fromIP.String(), theirRankStr)
			} else {
				if theirRank > LeaderRank {
					LeaderIP = fromIP
					LeaderRank = theirRank
					logrus.Infof("Choosen new leader %s of rank %d", LeaderIP, LeaderRank)
				}
			}
		}
	}()
}

func SendVoteRequest() {
	logrus.Info("Sending vote requests ... ")
	for lastByte := byte(0); true; lastByte++ {

		if lastByte != MyIP[3] {
			target := net.IP{MyIP[0], MyIP[1], MyIP[2], lastByte}
			err := com.SendUDPMessage(MyIP, MinionUDPClientPort, target, MinionUDPServerPort, "VREQ")
			if err != nil {
				logrus.Errorf("error while sending vote request to %s. got %s", target.String(), err.Error())
			}
		}
		if lastByte == 255 {
			break
		}
	}
	go func() {
		timer := time.NewTimer(5 * time.Second)
		<-timer.C
		if LeaderRank == Rank {
			logrus.Info("Current leader is my self")
		} else {
			logrus.Infof("Current leader is %s with rank %d", LeaderIP, LeaderRank)
		}
		defer timer.Stop()
	}()
}

func Start(ctx context.Context) {
	go func() {
		err := com.StartServer(ctx, MyIP, MinionUDPServerPort, MinionDaemonHandler)
		if err != nil {
			fmt.Sprintf(err.Error())
		}
	}()

	voteTicker := time.NewTicker(VoteTickDuration)
	stopTicker := make(chan bool)
	go func() {
		for {
			select {
			case <-stopTicker:
				return
			case <-voteTicker.C:
				SendVoteRequest()
			}
		}
	}()

	gracefulStop := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(gracefulStop, os.Interrupt)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	// Block until we receive our signal.
	logrus.Warn("Warming UP")
	<-gracefulStop

	defer func() {
		com.StopServer()
		voteTicker.Stop()
		stopTicker <- true
	}()

	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	logrus.Info("shutting down minion........ bye")

}
