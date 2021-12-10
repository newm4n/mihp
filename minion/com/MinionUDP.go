package com

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"net"
	"reflect"
)

var (
	UDPServerStopChannel    chan bool        = make(chan bool)
	UDPServerMessageChannel chan *UDPMessage = make(chan *UDPMessage)
	UDPConn                 *net.UDPConn
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
	logrus.Infof("UDP Server listening at %s:%d", ip.String(), port)
	ServerAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ip.String(), port))
	if err != nil {
		return err
	}

	/* Now listen at selected port */
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	if err != nil {
		return err
	}

	UDPConn = ServerConn

	defer ServerConn.Close()

	go func() {
		for {
			if ctx.Err() != nil {
				UDPServerStopChannel <- true
				break
			}
		}
	}()

	go func() {
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
		}
		break
	}
	fmt.Println("Emptying channel")
	for len(UDPServerMessageChannel) > 0 {
		<-UDPServerMessageChannel
	}
	fmt.Println("Channel emptied")
	return nil
}

// GetOutboundIP get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func SendUDPMessage(localIP net.IP, localPort int, targetIP net.IP, targetPort int, message string) error {
	ServerAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", targetIP.String(), targetPort))
	if err != nil {
		return err
	}
	LocalAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", localIP.String(), localPort))
	if err != nil {
		return err
	}
	Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
	if err != nil {
		return err
	}
	defer Conn.Close()
	buf := []byte(message)
	_, err = Conn.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func GetIPNetworkGroup(ip net.IP, mask net.IPMask) []net.IP {
	classA := bytesByMask(ip[0], mask[0])
	classB := bytesByMask(ip[1], mask[1])
	classC := bytesByMask(ip[2], mask[2])
	classD := bytesByMask(ip[3], mask[3])
	ips := make([]net.IP, 0)
	for _, a := range classA {
		for _, b := range classB {
			for _, c := range classC {
				for _, d := range classD {
					ips = append(ips, net.IP{a, b, c, d})
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
	for i := byte(0); true; i++ {
		if b&m == i&m {
			bytes = append(bytes, i)
		}
		if i == 255 {
			break
		}
	}
	return bytes
}

func NetmaskForSlash(slash int) net.IPMask {
	if slash >= 32 {
		return net.IPMask{255, 255, 255, 255}
	}
	if slash <= 0 {
		return net.IPMask{0, 0, 0, 0}
	}
	bits := (uint32(0xFFFFFFFF) >> slash) ^ uint32(0xFFFFFFFF)
	b1 := (bits >> 24) & 0x000000FF
	b2 := (bits >> 16) & 0x000000FF
	b3 := (bits >> 8) & 0x000000FF
	b4 := bits & 0x000000FF
	return net.IPMask{byte(b1), byte(b2), byte(b3), byte(b4)}
}
