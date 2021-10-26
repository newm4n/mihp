package probing

import (
	"bytes"
	"context"
	"fmt"
	"github.com/newm4n/mihp/pkg/errors"
	"github.com/newm4n/mihp/pkg/helper"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"
)

const (
	NotifTypeEmailSMTP = "SMTP"
	NotifTypeCallBack  = "CALLBACK"
	NotifTypeTelegram  = "TELEGRAM"
	NotifTypeSlack     = "SLACK"
	NotifTypeMSTeams   = "MSTEAMS"
)

type Probe struct {
	Name     string          `json:"name"`
	ID       string          `json:"id"`
	Requests []*ProbeRequest `json:"requests"`

	Cron string `json:"cron"`

	SuccessNotificationThreshold int `json:"success_notification_threshold"`
	FailNotificationThreshold    int `json:"fail_notification_threshold"`

	NotificationType string `json:"notification_type"`
}

func (p *Probe) CanStart(ctx context.Context) bool {
	if ctx.Err() != nil {
		logrus.Warnf("Will never ever start since context has error. got %s", ctx.Err())
		return false
	}
	cs, err := helper.NewCronStruct(p.Cron)
	if err != nil {
		logrus.Errorf("Will never ever start since cron is wrong. got %s", err.Error())
		return false
	}
	return cs.IsIn(time.Now())
}

func (p *Probe) Execute(ctx context.Context, pctx ProbeContext) error {
	if ctx.Err() != nil {
		return fmt.Errorf("%w : context probably timed-out", ctx.Err())
	}
	if p.CanStart(ctx) {
		pctx[fmt.Sprintf("probe.%s.id", p.Name)] = p.ID
		startTime := time.Now()
		pctx[fmt.Sprintf("probe.%s.starttime", p.Name)] = startTime

		defer func() {
			pctx[fmt.Sprintf("probe.%s.endtime", p.Name)] = time.Now()
			pctx[fmt.Sprintf("probe.%s.duration", p.Name)] = int(time.Now().Sub(startTime) / time.Millisecond)
		}()

		reqNames := make([]string, 0)
		for _, reqs := range p.Requests {
			reqNames = append(reqNames, reqs.Name)
		}

		pctx[fmt.Sprintf("probe.%s.req", p.Name)] = strings.Join(reqNames, ",")

		for seq, reqs := range p.Requests {
			err := reqs.Execute(ctx, p, seq, pctx)
			if err != nil {
				pctx[fmt.Sprintf("probe.%s.fail", p.Name)] = true
				pctx[fmt.Sprintf("probe.%s.success", p.Name)] = false
				return err
			}
		}
		pctx[fmt.Sprintf("probe.%s.fail", p.Name)] = false
		pctx[fmt.Sprintf("probe.%s.success", p.Name)] = true
	} else {
		logrus.Tracef("probe.%s can't start", p.Name)
	}
	return nil
}

type ProbeRequest struct {
	Name             string              `json:"name"`
	URL              string              `json:"url"`
	Method           string              `json:"method"`
	Headers          map[string][]string `json:"headers"`
	Body             string              `json:"body"`
	CertificateCheck string              `json:"certificate_check"`
	StartRequestIf   string              `json:"start_request_if"`
	SuccessIf        string              `json:"success_if"`
	FailIf           string              `json:"fail_if"`
}

func (pr *ProbeRequest) Execute(ctx context.Context, probe *Probe, sequence int, pctx ProbeContext) error {
	if ctx.Err() != nil {
		pctx[fmt.Sprintf("probe.%s.req.%s.error", probe.Name, pr.Name)] = ctx.Err()
		return fmt.Errorf("%w : context probably timed-out", ctx.Err())
	}
	pctx[fmt.Sprintf("probe.%s.req.%s.sequence", probe.Name, pr.Name)] = sequence
	if len(pr.StartRequestIf) > 0 {
		out, err := GoCelEvaluate(ctx, pr.StartRequestIf, pctx, reflect.Bool)
		if err != nil {
			pctx[fmt.Sprintf("probe.%s.req.%s.canstart", probe.Name, pr.Name)] = false
			pctx[fmt.Sprintf("probe.%s.req.%s.error", probe.Name, pr.Name)] = err
			return err
		}
		if !out.(bool) {
			pctx[fmt.Sprintf("probe.%s.req.%s.canstart", probe.Name, pr.Name)] = false
			return errors.ErrContextValueIsNotBool
		}
	}
	pctx[fmt.Sprintf("probe.%s.req.%s.canstart", probe.Name, pr.Name)] = true

	client := NewHttpClient(10, 10, false)
	var request *http.Request
	var err error

	urlItv, err := GoCelEvaluate(ctx, pr.URL, pctx, reflect.String)
	if err != nil {
		pctx[fmt.Sprintf("probe.%s.req.%s.error", probe.Name, pr.Name)] = err
		return err
	}
	URL := urlItv.(string)
	pctx[fmt.Sprintf("probe.%s.req.%s.url", probe.Name, pr.Name)] = URL

	methodItv, err := GoCelEvaluate(ctx, pr.Method, pctx, reflect.String)
	if err != nil {
		pctx[fmt.Sprintf("probe.%s.req.%s.error", probe.Name, pr.Name)] = err
		return err
	}
	METHOD := methodItv.(string)
	pctx[fmt.Sprintf("probe.%s.req.%s.method", probe.Name, pr.Name)] = METHOD

	if pr.Body != "" {
		bodyItv, err := GoCelEvaluate(ctx, pr.Body, pctx, reflect.String)
		if err != nil {
			pctx[fmt.Sprintf("probe.%s.req.%s.error", probe.Name, pr.Name)] = err
			return err
		}
		pctx[fmt.Sprintf("probe.%s.req.%s.body", probe.Name, pr.Name)] = bodyItv.(string)
		req, err := http.NewRequest(METHOD, URL, bytes.NewBuffer([]byte(bodyItv.(string))))
		if err != nil {
			pctx[fmt.Sprintf("probe.%s.req.%s.error", probe.Name, pr.Name)] = err
			return err
		}
		request = req
	} else {
		req, err := http.NewRequest(METHOD, URL, nil)
		if err != nil {
			pctx[fmt.Sprintf("probe.%s.req.%s.error", probe.Name, pr.Name)] = err
			return err
		}
		request = req
	}

	if pr.Headers != nil && len(pr.Headers) > 0 {
		headerKeys := make([]string, 0)
		for hKey, _ := range pr.Headers {
			headerKeys = append(headerKeys, hKey)
		}
		pctx[fmt.Sprintf("probe.%s.req.%s.header", probe.Name, pr.Name)] = headerKeys
		for hKey, hVals := range pr.Headers {
			headerValArr := make([]string, len(hVals))
			for idx, expr := range hVals {
				iv, err := GoCelEvaluate(ctx, expr, pctx, reflect.String)
				if err != nil {
					pctx[fmt.Sprintf("probe.%s.req.%s.error", probe.Name, pr.Name)] = err
					return err
				}
				headerValArr[idx] = iv.(string)
			}
			for _, hV := range headerValArr {
				request.Header.Add(hKey, hV)
			}
			pctx[fmt.Sprintf("probe.%s.req.%s.header.%s", probe.Name, pr.Name, hKey)] = headerValArr
		}
	}

	reqStartTime := time.Now()
	pctx[fmt.Sprintf("probe.%s.req.%s.starttime", probe.Name, pr.Name)] = reqStartTime

	response, err := client.Do(request)

	pctx[fmt.Sprintf("probe.%s.req.%s.duration", probe.Name, pr.Name)] = int(time.Now().Sub(reqStartTime) / time.Millisecond)

	if err != nil {
		pctx[fmt.Sprintf("probe.%s.req.%s.error", probe.Name, pr.Name)] = err
		return err
	}

	pctx[fmt.Sprintf("probe.%s.req.%s.resp.code", probe.Name, pr.Name)] = response.StatusCode

	if response.Header != nil && len(response.Header) > 0 {
		headerKeys := make([]string, 0)
		for hKey, _ := range response.Header {
			headerKeys = append(headerKeys, hKey)
		}
		pctx[fmt.Sprintf("probe.%s.req.%s.resp.header", probe.Name, pr.Name)] = headerKeys
		for hKey, hVals := range response.Header {
			pctx[fmt.Sprintf("probe.%s.req.%s.resp.header.%s", probe.Name, pr.Name, hKey)] = hVals
		}
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		pctx[fmt.Sprintf("probe.%s.req.%s.resp.body", probe.Name, pr.Name)] = ""
		pctx[fmt.Sprintf("probe.%s.req.%s.resp.body.size", probe.Name, pr.Name)] = 0
	} else {
		pctx[fmt.Sprintf("probe.%s.req.%s.resp.body.size", probe.Name, pr.Name)] = len(bodyBytes)
		if hVals := response.Header.Values("Content-Type"); hVals != nil {
			if strings.Contains(hVals[0], "text/") {
				pctx[fmt.Sprintf("probe.%s.req.%s.resp.body", probe.Name, pr.Name)] = string(bodyBytes)
			} else {
				pctx[fmt.Sprintf("probe.%s.req.%s.resp.body", probe.Name, pr.Name)] = "<binary>"
			}
		}
	}

	// Check success if
	if len(pr.SuccessIf) > 0 {
		out, err := GoCelEvaluate(ctx, pr.SuccessIf, pctx, reflect.Bool)
		if err != nil {
			pctx[fmt.Sprintf("probe.%s.req.%s.success", probe.Name, pr.Name)] = false
			pctx[fmt.Sprintf("probe.%s.req.%s.fail", probe.Name, pr.Name)] = true
			pctx[fmt.Sprintf("probe.%s.req.%s.error", probe.Name, pr.Name)] = err
			fmt.Errorf("error when evaluating SuccessIf. got %s", err.Error())
			return err
		}
		if !out.(bool) {
			pctx[fmt.Sprintf("probe.%s.req.%s.success", probe.Name, pr.Name)] = false
			pctx[fmt.Sprintf("probe.%s.req.%s.fail", probe.Name, pr.Name)] = true
			return nil
		}
		pctx[fmt.Sprintf("probe.%s.req.%s.success", probe.Name, pr.Name)] = true
		pctx[fmt.Sprintf("probe.%s.req.%s.fail", probe.Name, pr.Name)] = false
	} else if len(pr.FailIf) > 0 { // Check fail if
		out, err := GoCelEvaluate(ctx, pr.FailIf, pctx, reflect.Bool)
		if err != nil {
			pctx[fmt.Sprintf("probe.%s.req.%s.fail", probe.Name, pr.Name)] = true
			pctx[fmt.Sprintf("probe.%s.req.%s.success", probe.Name, pr.Name)] = false
			pctx[fmt.Sprintf("probe.%s.req.%s.error", probe.Name, pr.Name)] = err
			fmt.Errorf("error when evaluating FailIf. got %s", err.Error())
			return err
		}
		if !out.(bool) {
			pctx[fmt.Sprintf("probe.%s.req.%s.fail", probe.Name, pr.Name)] = true
			pctx[fmt.Sprintf("probe.%s.req.%s.success", probe.Name, pr.Name)] = false
			return nil
		}
		pctx[fmt.Sprintf("probe.%s.req.%s.fail", probe.Name, pr.Name)] = false
		pctx[fmt.Sprintf("probe.%s.req.%s.success", probe.Name, pr.Name)] = true
	}
	pctx[fmt.Sprintf("probe.%s.req.%s.fail", probe.Name, pr.Name)] = false
	pctx[fmt.Sprintf("probe.%s.req.%s.success", probe.Name, pr.Name)] = true

	return nil
}
