package probing

import (
	"context"
	"fmt"
	"github.com/newm4n/mihp/internal/probing/dummy"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func TestGoCelEvaluateExistBool(t *testing.T) {
	pc := NewProbeContext()
	pc["ref.existingBool"] = true
	expr := `ref.existingBool`
	out, err := GoCelEvaluate(context.Background(), expr, pc, reflect.Bool)
	assert.NoError(t, err)
	assert.True(t, out.(bool))
}

func TestGoCelEvaluateNonExistBool(t *testing.T) {
	pc := NewProbeContext()
	pc["ref.existingBool"] = true

	out, err := GoCelEvaluate(context.Background(), `IsDefined("ref.existingBool")`, pc, reflect.Bool)
	assert.NoError(t, err)
	assert.True(t, out.(bool))

	out, err = GoCelEvaluate(context.Background(), `IsDefined("ref.nonExistingBool") `, pc, reflect.Bool)
	assert.NoError(t, err)
	assert.False(t, out.(bool))
}

func TestProbe_Chaining(t *testing.T) {
	t.Log("Starting dummy server")
	srv := &dummy.DummyServer{}
	srv.Start()

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
		StartRequestIf:   "",
		SuccessIf:        "",
		FailIf:           "",
	}
	probe.Requests = append(probe.Requests, req1)

	defer func() {
		srv.Stop()
		t.Log("Dummy server stopped")
	}()
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

	assert.NoError(t, probe.Execute(context.Background(), pCtx))
	t.Log(pCtx.ToString(false))
}
