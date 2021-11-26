package com

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"
)

func HandlerTest(message *UDPMessage) {
	fmt.Println(message.Message)
}

func TestStartServer(t *testing.T) {
	t.Run("StartServerWithTimeout", func(t *testing.T) {
		timeout, _ := context.WithTimeout(context.Background(), 4*time.Second)
		ip := net.IP{0, 0, 0, 0}
		err := StartServer(timeout, ip, 54652, HandlerTest)
		if err != nil {
			t.Log(err.Error())
		}
	})
}

func TestGetOutboundIP(t *testing.T) {
	ip := GetOutboundIP()
	t.Log(ip.String())
	t.Log(ip[0])
	t.Log(ip[3])
}
