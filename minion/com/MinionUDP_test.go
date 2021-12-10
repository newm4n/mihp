package com

import (
	"bytes"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
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

func TestIPEqual(t *testing.T) {
	assert.Equal(t, net.IP{1, 2, 3, 4}, net.IP{1, 2, 3, 4})
	assert.True(t, bytes.Equal(net.IP{1, 2, 3, 4}, net.IP{1, 2, 3, 4}))
}

func TestGetIPNetworkGroup(t *testing.T) {
	ips := GetIPNetworkGroup(net.IP{1, 2, 3, 4}, net.IPMask{255, 255, 255, 255})
	assert.Equal(t, 1, len(ips))
	assert.Equal(t, byte(1), ips[0][0])
	assert.Equal(t, byte(2), ips[0][1])
	assert.Equal(t, byte(3), ips[0][2])
	assert.Equal(t, byte(4), ips[0][3])

	ips = GetIPNetworkGroup(net.IP{1, 2, 3, 4}, net.IPMask{255, 255, 255, 0})
	assert.Equal(t, 256, len(ips))

	for i, ip := range ips {
		assert.Equal(t, byte(1), ip[0])
		assert.Equal(t, byte(2), ip[1])
		assert.Equal(t, byte(3), ip[2])
		assert.Equal(t, byte(i), ip[3])
	}

	ips = GetIPNetworkGroup(net.IP{1, 2, 3, 4}, net.IPMask{255, 255, 0, 0})
	assert.Equal(t, 256, len(ips))
}

func TestGetOutboundIP(t *testing.T) {
	ip := GetOutboundIP()
	t.Log(ip.String())
	t.Log(ip[0])
	t.Log(ip[3])
}

func TestNetmaskForSlash(t *testing.T) {
	t.Run("Slash24", func(t *testing.T) {
		mask := NetmaskForSlash(24)
		assert.Equal(t, byte(255), mask[0])
		assert.Equal(t, byte(255), mask[1])
		assert.Equal(t, byte(255), mask[2])
		assert.Equal(t, byte(0), mask[3])
	})
	t.Run("Slash16", func(t *testing.T) {
		mask := NetmaskForSlash(16)
		assert.Equal(t, byte(255), mask[0])
		assert.Equal(t, byte(255), mask[1])
		assert.Equal(t, byte(0), mask[2])
		assert.Equal(t, byte(0), mask[3])
	})
	t.Run("Slash8", func(t *testing.T) {
		mask := NetmaskForSlash(8)
		assert.Equal(t, byte(255), mask[0])
		assert.Equal(t, byte(0), mask[1])
		assert.Equal(t, byte(0), mask[2])
		assert.Equal(t, byte(0), mask[3])
	})
}
