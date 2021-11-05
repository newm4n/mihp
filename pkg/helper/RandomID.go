package helper

import (
	"bytes"
	"math/rand"
	"time"
)

const (
	RandomIDLength = 20
)

var (
	CharSet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyz"
)

func init() {
	rand.Seed(time.Now().UnixMicro())
}

func RandomID() string {
	buff := &bytes.Buffer{}
	for buff.Len() < RandomIDLength {
		offset := rand.Intn(len(CharSet))
		buff.WriteString(CharSet[offset : rand.Intn(len(CharSet))+1])
	}
	return buff.String()
}
