package probing

import (
	"fmt"
	"github.com/newm4n/mihp/pkg/errors"
	"time"
)

type Probe struct {
	Name             string
	ID               string
	Requests         []*ProbeRequest
	SecondsInterval  int
	Since            time.Time
	Until            time.Time
	SuccessThreshold int
	FailThreshold    int
}

func (p *Probe) Execute(pctx ProbeContext) error {
	if p.Since.Before(time.Now()) && p.Until.After(time.Now()) {

		pctx[fmt.Sprintf("%s.id", p.Name)] = p.ID
		startTime := time.Now()
		pctx[fmt.Sprintf("%s.starttime", p.Name)] = startTime

		defer func() {
			pctx[fmt.Sprintf("%s.endtime", p.Name)] = time.Now()
			pctx[fmt.Sprintf("%s.duration", p.Name)] = time.Now().Sub(startTime) / time.Millisecond
		}()

	}
	return nil
}

type ProbeRequest struct {
	Name             string
	URL              string
	Method           string
	CertificateCheck string
	StartRequestIf   string
	SuccessIf        string
	FailIf           string
}

func (pr *ProbeRequest) Execute(probe *Probe, sequence int, pctx ProbeContext) error {
	pctx[fmt.Sprintf("%s.%s.req.%s.sequence", probe.Name, probe.ID, pr.Name)] = sequence
	if len(pr.StartRequestIf) > 0 {
		out, err := GoCelEvaluate(pr.StartRequestIf, pctx)
		if err != nil {
			return err
		}
		if canStart, ok := out.(bool); !ok {
			pctx[fmt.Sprintf("%s.req.%s.error", probe.Name, pr.Name)] = errors.ErrContextValueIsNotBool
			return errors.ErrContextValueIsNotBool
		} else if !canStart {
			pctx[fmt.Sprintf("%s.req.%s.canstart", probe.Name, pr.Name)] = false
			return errors.ErrContextValueIsNotBool
		} else {
			pctx[fmt.Sprintf("%s.req.%s.canstart", probe.Name, pr.Name)] = true
		}
	}

	// start probing.
	return nil

}
