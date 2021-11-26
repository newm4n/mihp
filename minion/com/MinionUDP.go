package com

import (
	"context"
	"fmt"
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
					fmt.Println("Error: ", err, reflect.TypeOf(err).String())
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
