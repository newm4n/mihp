package probing

import (
	"context"
	"fmt"
	"github.com/newm4n/mihp/internal"
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
	probe := &internal.Probe{
		Name:          "Local",
		ID:            "1001",
		Requests:      make([]*internal.ProbeRequest, 0),
		BaseURL:       fmt.Sprintf("http://localhost:%d", srv.Port),
		Cron:          "* * * * * * *",
		UpThreshold:   0,
		DownThreshold: 0,
	}
	req1 := &internal.ProbeRequest{
		Name:       "Login",
		PathExpr:   "\"/login\"",
		MethodExpr: `"GET"`,
		HeadersExpr: map[string][]string{
			"User-Agent": {"\"mihp/1.0.0 mihp is http probe\""},
		},
		BodyExpr:             "",
		CertificateCheckExpr: "false",
		SuccessIfExpr:        `IsDefined("probe.Local.req.Login.resp.code") && GetInt("probe.Local.req.Login.resp.code")==200`,
		FailIfExpr:           "",
	}
	probe.Requests = append(probe.Requests, req1)

	req2 := &internal.ProbeRequest{
		Name:       "Dashboard",
		PathExpr:   "\"/dashboard\"",
		MethodExpr: `"GET"`,
		HeadersExpr: map[string][]string{
			"User-Agent":    {`"mihp/1.0.0 mihp is http probe"`},
			"Authorization": {`GetStringElem("probe.Local.req.Login.resp.header.Testtoken",0)`},
		},
		BodyExpr:             "",
		CertificateCheckExpr: "false",
		StartRequestIfExpr:   `IsDefined("probe.Local.req.Login.success") && GetBool("probe.Local.req.Login.success") == true `,
		SuccessIfExpr:        `IsDefined("probe.Local.req.Dashboard.resp.code") && GetInt("probe.Local.req.Dashboard.resp.code")==200`,
		FailIfExpr:           "",
	}
	probe.Requests = append(probe.Requests, req2)

	time.Sleep(1500 * time.Millisecond)

	pCtx := internal.NewProbeContext()
	assert.NoError(t, ExecuteProbe(context.Background(), probe, pCtx, 10, true, true))
	t.Log(pCtx.ToString(false))
}

func TestProbe_ExecuteGoogle(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
	pCtx := internal.NewProbeContext()
	probe := &internal.Probe{
		Name:          "Google",
		ID:            "123",
		Requests:      make([]*internal.ProbeRequest, 0),
		BaseURL:       `https://google.com`,
		Cron:          "* * * * * * *",
		UpThreshold:   0,
		DownThreshold: 0,
	}
	req1 := &internal.ProbeRequest{
		Name:       "GoogleHome",
		PathExpr:   `"/"`,
		MethodExpr: `"GET"`,
		HeadersExpr: map[string][]string{
			"User-Agent": {"\"mihp/1.0.0 mihp is http probe\""},
		},
		BodyExpr:             "",
		CertificateCheckExpr: "false",
		StartRequestIfExpr:   "",
		SuccessIfExpr:        "",
		FailIfExpr:           "",
	}
	probe.Requests = append(probe.Requests, req1)

	time.Sleep(1500 * time.Millisecond)

	assert.NoError(t, ExecuteProbe(context.Background(), probe, pCtx, 10, true, true))

	assert.True(t, pCtx["probe.Google.success"].(bool))
	t.Log(pCtx.ToString(false))
}
