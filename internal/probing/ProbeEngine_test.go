package probing

import (
	"context"
	"fmt"
	"github.com/newm4n/mihp/internal/probing/dummy"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestProbe_Chaining(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
	t.Log("Starting dummy server")
	srv := &dummy.DummyServer{}
	srv.Start()
	defer func() {
		srv.Stop()
		t.Log("Dummy server stopped")
	}()

	time.Sleep(3 * time.Second)
	//srv.Port

	logrus.SetLevel(logrus.TraceLevel)
	//pCtx := NewProbeContext()
	probe := &Probe{
		Name:                         "Local",
		ID:                           "1001",
		Requests:                     make([]*ProbeRequest, 0),
		Cron:                         "* * * * * * *",
		SuccessNotificationThreshold: 0,
		FailNotificationThreshold:    0,
	}
	req1 := &ProbeRequest{
		Name:   "Login",
		URL:    fmt.Sprintf("\"http://localhost:%d/login\"", srv.Port),
		Method: `"GET"`,
		Headers: map[string][]string{
			"User-Agent": {"\"mihp/1.0.0 mihp is http probe\""},
		},
		Body:             "",
		CertificateCheck: "false",
		SuccessIf:        `IsDefined("probe.Local.req.Login.resp.code") && GetInt("probe.Local.req.Login.resp.code")==200`,
		FailIf:           "",
	}
	probe.Requests = append(probe.Requests, req1)

	req2 := &ProbeRequest{
		Name:   "Dashboard",
		URL:    fmt.Sprintf("\"http://localhost:%d/dashboard\"", srv.Port),
		Method: `"GET"`,
		Headers: map[string][]string{
			"User-Agent":    {`"mihp/1.0.0 mihp is http probe"`},
			"Authorization": {`GetStringElem("probe.Local.req.Login.resp.header.Testtoken",0)`},
		},
		Body:             "",
		CertificateCheck: "false",
		StartRequestIf:   `IsDefined("probe.Local.req.Login.success") && GetBool("probe.Local.req.Login.success") == true `,
		SuccessIf:        `IsDefined("probe.Local.req.Dashboard.resp.code") && GetInt("probe.Local.req.Dashboard.resp.code")==200`,
		FailIf:           "",
	}
	probe.Requests = append(probe.Requests, req2)

	time.Sleep(1500 * time.Millisecond)

	pCtx := NewProbeContext()
	assert.NoError(t, ExecuteProbe(context.Background(), probe, pCtx))
	t.Log(pCtx.ToString(false))
}

func TestProbe_ExecuteGoogle(t *testing.T) {
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

	assert.NoError(t, ExecuteProbe(context.Background(), probe, pCtx))

	assert.True(t, pCtx["probe.Google.success"].(bool))
	t.Log(pCtx.ToString(false))
}
