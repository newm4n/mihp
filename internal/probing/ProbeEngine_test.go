package probing

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestProbe_Execute(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
	pCtx := NewProbeContext()
	probe := &Probe{
		Name:                         "Google",
		ID:                           "123",
		Requests:                     make([]*ProbeRequest, 0),
		Cron:                         "* * * * * * *",
		SuccessNotificationThreshold: 0,
		FailNotificationThreshold:    0,
	}
	req1 := &ProbeRequest{
		Name:   "GoogleHome",
		URL:    `"https://google.com"`,
		Method: `"GET"`,
		Headers: map[string][]string{
			"User-Agent": {"\"mihp/1.0.0 mihp is http probe\""},
		},
		Body:             "",
		CertificateCheck: "false",
		StartRequestIf:   "",
		SuccessIf:        "",
		FailIf:           "",
	}
	probe.Requests = append(probe.Requests, req1)

	time.Sleep(1500 * time.Millisecond)

	assert.NoError(t, probe.Execute(context.Background(), pCtx))
	t.Log(pCtx.ToString(false))
}
