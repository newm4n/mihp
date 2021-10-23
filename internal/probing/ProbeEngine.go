package probing

import (
	"bytes"
	"context"
	"fmt"
	"github.com/newm4n/mihp/pkg/errors"
	"io/ioutil"
	"net/http"
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

func (p *Probe) Execute(ctx context.Context, pctx ProbeContext) error {
	if p.Since.Before(time.Now()) && p.Until.After(time.Now()) {

		pctx[fmt.Sprintf("%s.id", p.Name)] = p.ID
		startTime := time.Now()
		pctx[fmt.Sprintf("%s.starttime", p.Name)] = startTime

		defer func() {
			pctx[fmt.Sprintf("%s.endtime", p.Name)] = time.Now()
			pctx[fmt.Sprintf("%s.duration", p.Name)] = time.Now().Sub(startTime) / time.Millisecond
		}()

		for seq, reqs := range p.Requests {
			reqs.Execute(ctx, p, seq, pctx)
		}
	}
	return nil
}

type ProbeRequest struct {
	Name             string
	URL              string
	Method           string
	Headers          map[string][]string
	Body             string
	CertificateCheck string
	StartRequestIf   string
	SuccessIf        string
	FailIf           string
}

func (pr *ProbeRequest) Execute(ctx context.Context, probe *Probe, sequence int, pctx ProbeContext) error {
	pctx[fmt.Sprintf("%s.%s.req.%s.sequence", probe.Name, probe.ID, pr.Name)] = sequence
	if len(pr.StartRequestIf) > 0 {
		out, err := GoCelEvaluate(pr.StartRequestIf, pctx)
		if err != nil {
			pctx[fmt.Sprintf("%s.req.%s.error", probe.Name, pr.Name)] = err
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

	client := NewHttpClient(10, 10, false)
	var request *http.Request
	var err error

	var URL string
	urlItv, err := GoCelEvaluate(pr.URL, pctx)
	if err != nil {
		pctx[fmt.Sprintf("%s.req.%s.error", probe.Name, pr.Name)] = err
		return err
	}
	if theUrl, ok := urlItv.(string); !ok {
		pctx[fmt.Sprintf("%s.req.%s.error", probe.Name, pr.Name)] = err
		return err
	} else {
		URL = theUrl
	}

	var METHOD string
	methodItv, err := GoCelEvaluate(pr.Method, pctx)
	if err != nil {
		pctx[fmt.Sprintf("%s.req.%s.error", probe.Name, pr.Name)] = err
		return err
	}
	if theMethod, ok := methodItv.(string); !ok {
		pctx[fmt.Sprintf("%s.req.%s.error", probe.Name, pr.Name)] = err
		return err
	} else {
		METHOD = theMethod
	}

	if pr.Body != "" {
		bodyItv, err := GoCelEvaluate(pr.Body, pctx)
		if theBody, ok := bodyItv.(string); !ok {
			pctx[fmt.Sprintf("%s.req.%s.error", probe.Name, pr.Name)] = err
			return err
		} else {
			req, err := http.NewRequest(METHOD, URL, bytes.NewBuffer([]byte(theBody)))
			if err != nil {
				pctx[fmt.Sprintf("%s.req.%s.error", probe.Name, pr.Name)] = err
				return err
			}
			request = req
		}
	} else {
		req, err := http.NewRequest(METHOD, URL, nil)
		if err != nil {
			pctx[fmt.Sprintf("%s.req.%s.error", probe.Name, pr.Name)] = err
			return err
		}
		request = req
	}

	for hKey, hVals := range pr.Headers {
		for _, hVal := range hVals {
			hValItv, err := GoCelEvaluate(hVal, pctx)
			if err != nil {
				pctx[fmt.Sprintf("%s.req.%s.error", probe.Name, pr.Name)] = err
				return err
			}
			if hStrVal, ok := hValItv.(string); !ok {
				pctx[fmt.Sprintf("%s.req.%s.error", probe.Name, pr.Name)] = err
				return err
			} else {
				request.Header.Add(hKey, hStrVal)
			}
		}
	}

	reqStartTime := time.Now()
	pctx[fmt.Sprintf("%s.req.%s.starttime", probe.Name, pr.Name)] = reqStartTime

	response, err := client.Do(request)

	pctx[fmt.Sprintf("%s.req.%s.duration", probe.Name, pr.Name)] = time.Now().Sub(reqStartTime) / time.Millisecond

	if err != nil {
		pctx[fmt.Sprintf("%s.req.%s.error", probe.Name, pr.Name)] = err
		return err
	}

	pctx[fmt.Sprintf("%s.req.%s.resp.code", probe.Name, pr.Name)] = response.StatusCode

	for rhKey, rhVals := range response.Header {
		pctx[fmt.Sprintf("%s.req.%s.resp.header.%s.count", probe.Name, pr.Name, rhKey)] = len(rhVals)
		for hIdx, rhVal := range rhVals {
			pctx[fmt.Sprintf("%s.req.%s.resp.header.%s.%d", probe.Name, pr.Name, rhKey, hIdx+1)] = rhVal
		}
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		pctx[fmt.Sprintf("%s.req.%s.resp.body", probe.Name, pr.Name)] = ""
		pctx[fmt.Sprintf("%s.req.%s.resp.body.size", probe.Name, pr.Name)] = 0
	} else {
		pctx[fmt.Sprintf("%s.req.%s.resp.body", probe.Name, pr.Name)] = string(bodyBytes)
		pctx[fmt.Sprintf("%s.req.%s.resp.body.size", probe.Name, pr.Name)] = len(bodyBytes)
	}

	// start probing.
	return nil

}
