package com

import (
	bytes2 "bytes"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"net"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	UDPServerStopChannel    chan bool        = make(chan bool)
	UDPServerMessageChannel chan *UDPMessage = make(chan *UDPMessage)
	UDPConn                 *net.UDPConn
	Mutex                   sync.Mutex
)

type UDPMessageHandler func(message *UDPMessage)

type UDPMessage struct {
	Conn     *net.UDPConn
	FromAddr *net.UDPAddr
	Message  string
}

func StopServer() {
	UDPConn.Close()
	UDPServerStopChannel <- true
}

func StartServer(ctx context.Context, ip net.IP, port int, handler UDPMessageHandler) error {
	ServerAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ip.String(), port))
	if err != nil {
		logrus.Errorf("Error while resolving UDP Address. got %s", err.Error())
		return err
	}

	/* Now listen at selected port */
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	if err != nil {
		logrus.Errorf("Error while start listening UDP. got %s", err.Error())
		return err
	}

	UDPConn = ServerConn

	defer func() {
		ServerConn.Close()
		logrus.Warnf("Server connection closed")
	}()

	go func() {
		for {
			if ctx.Err() != nil {
				UDPServerStopChannel <- true
				logrus.Errorf("Context canceled. got %s", ctx.Err().Error())
				break
			}
		}
	}()

	go func() {
		time.Sleep(1 * time.Second)
		buf := make([]byte, 1024)
		for {
			n, addr, err := ServerConn.ReadFromUDP(buf)
			if err != nil {
				if reflect.TypeOf(err).String() != "*net.OpError" {
					logrus.Errorf("Error: got %s of type %s", err, reflect.TypeOf(err).String())
				} else {
					logrus.Errorf("Error: got %s", err)
				}
				break
			}
			UDPServerMessageChannel <- &UDPMessage{
				Conn:     ServerConn,
				FromAddr: addr,
				Message:  string(buf[0:n]),
			}

		}
	}()

	for {
		select {
		case <-UDPServerStopChannel:
			break
		case v := <-UDPServerMessageChannel:
			handler(v)
			continue
		}
		break
	}
	logrus.Warn("Emptying channel")
	for len(UDPServerMessageChannel) > 0 {
		<-UDPServerMessageChannel
	}
	logrus.Warn("Channel emptied")
	return nil
}

// GetOutboundIP get preferred outbound ip of this machine
func GetOutboundIP() IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return ParseIP(localAddr.IP.String())
}

func SendUDPMessage(Conn *net.UDPConn, targetIP IP, targetPort int, message string) error {
	Mutex.Lock()
	defer Mutex.Unlock()
	ServerAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", targetIP.String(), targetPort))
	if err != nil {
		return err
	}
	buf := []byte(message)
	_, err = Conn.WriteToUDP(buf, ServerAddr)
	if err != nil {
		return err
	}
	return nil
}

func GetIPNetworkGroup(ip IP, mask NetMask) []IP {
	classA := bytesByMask(ip[0], mask[0])
	classB := bytesByMask(ip[1], mask[1])
	classC := bytesByMask(ip[2], mask[2])
	classD := bytesByMask(ip[3], mask[3])
	ips := make([]IP, 0)
	for _, a := range classA {
		for _, b := range classB {
			for _, c := range classC {
				for _, d := range classD {
					nip := IP{a, b, c, d}
					ips = append(ips, nip)
					if len(ips) == 256 {
						return ips
					}
				}
			}
		}
	}
	return ips
}

func bytesByMask(b, m byte) []byte {
	bytes := make([]byte, 0)
	buff := &bytes2.Buffer{}

	for i := byte(0); true; i++ {
		if b&m == i&m {
			bytes = append(bytes, i)
			buff.WriteString(fmt.Sprintf("%d ", i))
		}
		if i == 255 {
			break
		}
	}

	return bytes
}

func NetmaskForSlash(slash int) NetMask {
	if slash >= 32 {
		return NetMask{255, 255, 255, 255}
	}
	if slash <= 0 {
		return NetMask{0, 0, 0, 0}
	}
	bits := (uint32(0xFFFFFFFF) >> slash) ^ uint32(0xFFFFFFFF)
	b1 := (bits >> 24) & 0x000000FF
	b2 := (bits >> 16) & 0x000000FF
	b3 := (bits >> 8) & 0x000000FF
	b4 := bits & 0x000000FF
	return NetMask{byte(b1), byte(b2), byte(b3), byte(b4)}
}

func ParseNetMask(mask string) NetMask {
	ret := make(NetMask, 4)
	bytes := strings.Split(mask, ".")
	if len(bytes) != 4 {
		return nil
	}
	for i, s := range bytes {
		b, err := strconv.ParseUint(s, 10, 8)
		if err != nil {
			return nil
		}
		if b > 255 {
			return nil
		}
		ret[i] = byte(b)
	}
	return ret
}

func ParseIP(ip string) IP {
	nm := ParseNetMask(ip)
	if nm == nil {
		return nil
	}
	return IP{nm[0], nm[1], nm[2], nm[3]}
}

type NetMask []byte

func (nm NetMask) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", nm[0], nm[1], nm[2], nm[3])
}

type IP []byte

func (ip IP) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
}

func (ip IP) ToNetIP() net.IP {
	return net.ParseIP(ip.String())
}
