package minion

import (
	"bytes"
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
	MinionUDPPort = 62891
)

var (
	Config               *internal.MIHPConfig
	EmailNotifChannel    = make(map[string]*probing.ProbeEventProcessor)
	CallbackNotifChannel = make(map[string]*probing.ProbeEventProcessor)
	LogNotifChannel      = make(map[string]*probing.ProbeEventProcessor)
	Rank                 uint64
	VoteCount            uint64
	MyIP                 com.IP
	MyNetmask            com.NetMask
	LeaderIP             com.IP
	LeaderRank           uint64
	VoteTickDuration     = 10 * time.Second // should be arround every 5 minutes ?
	PingTickDuration     = 30 * time.Second
	MinionGroupList      = make(map[string]*PingPong)
	CanPing              = true
)

func init() {
	Rank = rand.Uint64()
	LeaderRank = Rank
	MyIP = com.GetOutboundIP()
	LeaderIP = com.IP{MyIP[0], MyIP[1], MyIP[2], MyIP[3]}
}

func Initialize(MIHPConfig *internal.MIHPConfig) {
	Config = MIHPConfig
	if Config.ProbePool != nil {
		for _, probe := range Config.ProbePool {
			AcceptProbe(probe)
		}
	}

	MyIP = com.ParseIP(Config.Minion.MinionIP)
	MyNetmask = com.ParseNetMask(Config.Minion.MinionNetwork)

	if MyIP == nil {
		MyIP = com.GetOutboundIP()
		LeaderIP = com.IP{MyIP[0], MyIP[1], MyIP[2], MyIP[3]}
		fmt.Printf("Bind IP missing from config, Minion will bind to ip %s", MyIP.String())
	}

	if MyNetmask == nil {
		MyNetmask = com.NetMask{255, 255, 255, 0}
		fmt.Printf("Net Mask missing from config, Minion will us ip mask %s", MyNetmask.String())
	}

	fmt.Printf("BIND IP     : %s\n", MyIP.String())
	fmt.Printf("NET MASK    : %s\n", MyNetmask.String())
	fmt.Printf("RANK        : %d\n", Rank)
}

func AcceptProbe(probe *internal.Probe) {
	if probe.SMTPNotification != nil {
		EmailNotifChannel[probe.ID] = probing.NewProbeEventProcessor(probing.LogTrigger)
	}
	// todo finish this MINION
}

func MinionDaemonHandler(message *com.UDPMessage) {
	go func() {
		fromIP := com.ParseIP(message.FromAddr.IP.String())
		if strings.HasPrefix(message.Message, "VREQ ") {
			logrus.Infof("Receive vote request from %s", fromIP)
			logrus.Infof("Sending vote response to %s", fromIP)
			err := com.SendUDPMessage(message.Conn, fromIP, MinionUDPPort, fmt.Sprintf("VRES %d", Rank))
			if err != nil {
				logrus.Errorf("error while sending vote response to %s. got %s", fromIP.String(), err.Error())
			}
			vCountStr := strings.TrimSpace(message.Message[5:])
			vCount, err := strconv.ParseUint(vCountStr, 10, 64)
			if vCount == 0 {
				SendVoteRequest(message.Conn)
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
			if pp, ok := MinionGroupList[fromIP.String()]; !ok {
				MinionGroupList[fromIP.String()] = &PingPong{
					Ping:         time.Now(),
					Pong:         time.Now(),
					PongReceived: true,
				}
			} else {
				pp.Pong = time.Now()
				pp.PongReceived = true
			}
		} else if message.Message == "PING" {
			err := com.SendUDPMessage(message.Conn, fromIP, MinionUDPPort, "PONG")
			if err != nil {
				logrus.Errorf("error while sending vote response to %s. got %s", fromIP.String(), err.Error())
			}
		} else if message.Message == "PONG" {
			if pp, ok := MinionGroupList[fromIP.String()]; ok {
				pp.Pong = time.Now()
				pp.PongReceived = true
			}
		}
	}()
}

func SendPingRequests(conn *net.UDPConn) {
	if CanPing {
		for k, pp := range MinionGroupList {
			if MyIP.String() == k {
				continue
			}
			if pp.PongReceived == false {
				if time.Since(pp.Ping) > 10*time.Second {
					logrus.Warnf("Node %s not respoinding to ping for %s. It probably dead and removed on the next vote.", k, time.Since(pp.Ping).String())
				}
			}
			target := com.ParseIP(k)
			err := com.SendUDPMessage(conn, target, MinionUDPPort, fmt.Sprintf("PING", VoteCount))
			pp.PongReceived = false
			pp.Ping = time.Now()
			if err != nil {
				logrus.Errorf("error while sending PING request to %s. got %s", k, err.Error())
			}
		}
	}
}

func SendVoteRequest(conn *net.UDPConn) {
	logrus.Info("Sending vote requests ... ")

	LeaderIP = MyIP
	LeaderRank = Rank

	CanPing = false
	for _, ip := range com.GetIPNetworkGroup(MyIP, MyNetmask) {
		if bytes.Equal(MyIP, ip) {
			continue
		}

		err := com.SendUDPMessage(conn, ip, MinionUDPPort, fmt.Sprintf("VREQ %d", VoteCount))
		if err != nil {
			logrus.Errorf("error while sending vote request to %s. got %s", ip.String(), err.Error())
		}
		VoteCount++
	}
	go func() {
		timer := time.NewTimer(5 * time.Second)
		<-timer.C
		if LeaderRank == Rank {
			logrus.Info("Current leader is my self")
		} else {
			logrus.Infof("Current leader is %s with rank %d", LeaderIP, LeaderRank)
		}
		CanPing = true
		defer timer.Stop()
	}()
}

func Start(ctx context.Context, config *internal.MIHPConfig) {
	Initialize(config)

	go func() {
		err := com.StartServer(ctx, MyIP.ToNetIP(), MinionUDPPort, MinionDaemonHandler)
		if err != nil {
			fmt.Sprintf(err.Error())
			os.Exit(1)
		}
	}()

	voteTicker := time.NewTicker(VoteTickDuration)
	stopVoteTicker := make(chan bool)
	go func() {
		for {
			select {
			case <-stopVoteTicker:
				return
			case <-voteTicker.C:
				SendVoteRequest(com.UDPConn)
			}
		}
	}()

	pingTicker := time.NewTicker(PingTickDuration)
	stopPingTicker := make(chan bool)
	go func() {
		for {
			select {
			case <-stopPingTicker:
				return
			case <-pingTicker.C:
				SendPingRequests(com.UDPConn)
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
		stopVoteTicker <- true

		pingTicker.Stop()
		stopPingTicker <- true
	}()

	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	logrus.Info("shutting down minion........ bye")
}

type PingPong struct {
	Ping         time.Time
	Pong         time.Time
	PongReceived bool
}

func (pp *PingPong) Duration() time.Duration {
	if pp.PongReceived {
		return pp.Pong.Sub(pp.Ping)
	}
	return 0
}
