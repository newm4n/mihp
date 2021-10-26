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

const (
	TimeoutsSecond = 10
)

var (
	engineLog = logrus.WithField("module", "ProbeEngine")
)

func ProbeCanStartBySchedule(ctx context.Context, probe *Probe) bool {
	if ctx.Err() != nil {
		engineLog.Warnf("Will never ever start since context has error. got %s", ctx.Err())
		return false
	}
	cs, err := helper.NewCronStruct(probe.Cron)
	if err != nil {
		engineLog.Errorf("Will never ever start since cron is wrong. got %s", err.Error())
		return false
	}
	return cs.IsIn(time.Now())
}

func ExecuteProbe(ctx context.Context, probe *Probe, pctx ProbeContext) error {
	probeLog := engineLog.WithField("probe", probe.Name)
	if ctx.Err() != nil {
		probeLog.Errorf("context error. got %s", ctx.Err())
		return fmt.Errorf("%w : context probably timed-out. got %s", errors.ErrContextError, ctx.Err())
	}
	if ProbeCanStartBySchedule(ctx, probe) {
		pctx[fmt.Sprintf("probe.%s.id", probe.Name)] = probe.ID
		startTime := time.Now()
		pctx[fmt.Sprintf("probe.%s.starttime", probe.Name)] = startTime

		defer func() {
			pctx[fmt.Sprintf("probe.%s.endtime", probe.Name)] = time.Now()
			pctx[fmt.Sprintf("probe.%s.duration", probe.Name)] = time.Now().Sub(startTime)
		}()

		reqNames := make([]string, 0)
		for _, reqs := range probe.Requests {
			reqNames = append(reqNames, reqs.Name)
		}

		pctx[fmt.Sprintf("probe.%s.req", probe.Name)] = strings.Join(reqNames, ",")

		for seq, reqs := range probe.Requests {
			err := ExecuteProbeRequest(ctx, probe, reqs, seq, pctx)
			if err != nil {
				pctx[fmt.Sprintf("probe.%s.fail", probe.Name)] = true
				pctx[fmt.Sprintf("probe.%s.success", probe.Name)] = false
				probeLog.Errorf("error when execute probe request %s. got %s", reqs.Name, err.Error())
				return err
			}
		}
		pctx[fmt.Sprintf("probe.%s.fail", probe.Name)] = false
		pctx[fmt.Sprintf("probe.%s.success", probe.Name)] = true
	} else {
		probeLog.Tracef("probe.%s can't start", probe.Name)
	}
	return nil
}

func ExecuteProbeRequest(ctx context.Context, probe *Probe, probeRequest *ProbeRequest, sequence int, pctx ProbeContext) error {

	requestLog := engineLog.WithField("probe", probe.Name).WithField("request", probeRequest.Name)

	if ctx.Err() != nil {
		pctx[fmt.Sprintf("probe.%s.req.%s.error", probe.Name, probeRequest.Name)] = ctx.Err()
		requestLog.Errorf("context error. got %s", ctx.Err())
		return fmt.Errorf("%w : context probably timed-out. got %s", errors.ErrContextError, ctx.Err())
	}
	pctx[fmt.Sprintf("probe.%s.req.%s.sequence", probe.Name, probeRequest.Name)] = sequence
	if len(probeRequest.StartRequestIf) > 0 {
		out, err := GoCelEvaluate(ctx, probeRequest.StartRequestIf, pctx, reflect.Bool)
		if err != nil {
			requestLog.Errorf("error during evaluating StartRequestIf. got %s", err.Error())
			pctx[fmt.Sprintf("probe.%s.req.%s.canstart", probe.Name, probeRequest.Name)] = false
			pctx[fmt.Sprintf("probe.%s.req.%s.error", probe.Name, probeRequest.Name)] = err
			return fmt.Errorf("%w : probe %s request %s parsing StartRequestIf parsing error [%s]", err, probe.Name, probeRequest.Name, probeRequest.StartRequestIf)
		}
		if !out.(bool) {
			requestLog.Tracef("evaluation of StartRequestIf [%s] says that probe should not start", probeRequest.StartRequestIf)
			pctx[fmt.Sprintf("probe.%s.req.%s.canstart", probe.Name, probeRequest.Name)] = false
			return fmt.Errorf("%w : probe %s request %s can not start", errors.ErrStartRequestIfIsFalse, probe.Name, probeRequest.Name)
		}
	}
	requestLog.Tracef("StartRequestIf is OK")
	pctx[fmt.Sprintf("probe.%s.req.%s.canstart", probe.Name, probeRequest.Name)] = true

	requestLog.Tracef("Retriefing HTTP client with %d second timeout", TimeoutsSecond)
	client := NewHttpClient(TimeoutsSecond, TimeoutsSecond, false)
	var request *http.Request
	var err error

	requestLog.Tracef("Evaluating URL [%s]", probeRequest.URL)
	urlItv, err := GoCelEvaluate(ctx, probeRequest.URL, pctx, reflect.String)
	if err != nil {
		requestLog.Errorf("Error evaluating URL [%s] got %s", probeRequest.URL, err.Error())
		pctx[fmt.Sprintf("probe.%s.req.%s.error", probe.Name, probeRequest.Name)] = err
		return fmt.Errorf("%w : error while parsing URL", err)
	}
	URL := urlItv.(string)
	requestLog.Tracef("URL [%s] evaluated as [%s]", probeRequest.URL, URL)
	pctx[fmt.Sprintf("probe.%s.req.%s.url", probe.Name, probeRequest.Name)] = URL

	requestLog.Tracef("Evaluating Method [%s]", probeRequest.Method)
	methodItv, err := GoCelEvaluate(ctx, probeRequest.Method, pctx, reflect.String)
	if err != nil {
		requestLog.Errorf("Error evaluating Method [%s] got %s", probeRequest.Method, err.Error())
		pctx[fmt.Sprintf("probe.%s.req.%s.error", probe.Name, probeRequest.Name)] = err
		return fmt.Errorf("%w : error while parsing METHOD", err)
	}
	METHOD := methodItv.(string)
	requestLog.Tracef("Method [%s] evaluated as [%s]", probeRequest.Method, METHOD)
	pctx[fmt.Sprintf("probe.%s.req.%s.method", probe.Name, probeRequest.Name)] = METHOD

	if probeRequest.Body != "" {
		requestLog.Tracef("Evaluating Body [%s]", probeRequest.Body)
		bodyItv, err := GoCelEvaluate(ctx, probeRequest.Body, pctx, reflect.String)
		if err != nil {
			requestLog.Errorf("Error evaluating Body [%s] got %s", probeRequest.Body, err.Error())
			pctx[fmt.Sprintf("probe.%s.req.%s.error", probe.Name, probeRequest.Name)] = err
			return fmt.Errorf("%w : error while parsing Request Body", err)
		}
		pctx[fmt.Sprintf("probe.%s.req.%s.body", probe.Name, probeRequest.Name)] = bodyItv.(string)

		requestLog.Tracef("Body [%s] evaluated as [%s]", probeRequest.Body, bodyItv.(string))

		req, err := http.NewRequest(METHOD, URL, bytes.NewBuffer([]byte(bodyItv.(string))))
		if err != nil {
			requestLog.Errorf("Error while creating new http Request. got %s", err.Error())
			pctx[fmt.Sprintf("probe.%s.req.%s.error", probe.Name, probeRequest.Name)] = err
			return fmt.Errorf("%w : got %s", errors.ErrCreateHttpClient, err.Error())
		}
		request = req
	} else {
		req, err := http.NewRequest(METHOD, URL, nil)
		if err != nil {
			requestLog.Errorf("Error while creating new http Request. got %s", err.Error())
			pctx[fmt.Sprintf("probe.%s.req.%s.error", probe.Name, probeRequest.Name)] = err
			return fmt.Errorf("%w : got %s", errors.ErrCreateHttpRequest, err.Error())
		}
		request = req
	}

	if probeRequest.Headers != nil && len(probeRequest.Headers) > 0 {
		requestLog.Tracef("Parsing %d request headers", len(probeRequest.Headers))
		headerKeys := make([]string, 0)
		for hKey, _ := range probeRequest.Headers {
			headerKeys = append(headerKeys, hKey)
		}
		pctx[fmt.Sprintf("probe.%s.req.%s.header", probe.Name, probeRequest.Name)] = headerKeys
		for hKey, hVals := range probeRequest.Headers {
			headerValArr := make([]string, len(hVals))
			for idx, expr := range hVals {
				requestLog.Tracef("Parsing request headers [%s] = [%s]", hKey, expr)
				iv, err := GoCelEvaluate(ctx, expr, pctx, reflect.String)
				if err != nil {
					requestLog.Errorf("Error Parsing request headers [%s] = [%s]. got %s", hKey, expr, err.Error())
					pctx[fmt.Sprintf("probe.%s.req.%s.error", probe.Name, probeRequest.Name)] = err
					return fmt.Errorf("%w : error while evaluating probe %s request %s header %s expression [%s]", err, probe.Name, probeRequest.Name, hKey, expr)
				}
				requestLog.Tracef("Parsing request headers [%s] = [%s] --> [%s]", hKey, expr, iv.(string))
				headerValArr[idx] = iv.(string)
			}
			for _, hV := range headerValArr {
				request.Header.Add(hKey, hV)
			}
			pctx[fmt.Sprintf("probe.%s.req.%s.header.%s", probe.Name, probeRequest.Name, hKey)] = headerValArr
		}
	}

	reqStartTime := time.Now()
	pctx[fmt.Sprintf("probe.%s.req.%s.starttime", probe.Name, probeRequest.Name)] = reqStartTime
	requestLog.Tracef("Start calling http request.")

	response, err := client.Do(request)

	pctx[fmt.Sprintf("probe.%s.req.%s.duration", probe.Name, probeRequest.Name)] = time.Now().Sub(reqStartTime)
	requestLog.Tracef("Calling http request. Takes %s", time.Now().Sub(reqStartTime))

	if err != nil {
		requestLog.Errorf("Calling http request. Got %s", err.Error())
		pctx[fmt.Sprintf("probe.%s.req.%s.error", probe.Name, probeRequest.Name)] = err
		return fmt.Errorf("%w : http error for probe %s request %s got %s", errors.ErrHttpCallError, probe.Name, probeRequest.Name, err.Error())
	}

	pctx[fmt.Sprintf("probe.%s.req.%s.resp.code", probe.Name, probeRequest.Name)] = response.StatusCode
	requestLog.Tracef("Http response code is %d", response.StatusCode)

	requestLog.Tracef("Http response has %d headers", len(response.Header))
	if response.Header != nil && len(response.Header) > 0 {
		headerKeys := make([]string, 0)
		for hKey, _ := range response.Header {
			headerKeys = append(headerKeys, hKey)
		}
		pctx[fmt.Sprintf("probe.%s.req.%s.resp.header", probe.Name, probeRequest.Name)] = headerKeys
		for hKey, hVals := range response.Header {
			requestLog.Tracef("Http response header [%s] = [%s]", hKey, strings.Join(hVals, ","))
			pctx[fmt.Sprintf("probe.%s.req.%s.resp.header.%s", probe.Name, probeRequest.Name, hKey)] = hVals
		}
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		requestLog.Errorf("Error http response body reading. got %s", err.Error())
		pctx[fmt.Sprintf("probe.%s.req.%s.resp.body", probe.Name, probeRequest.Name)] = ""
		pctx[fmt.Sprintf("probe.%s.req.%s.resp.body.size", probe.Name, probeRequest.Name)] = 0
	} else {
		requestLog.Tracef("Http response body size is %d bytes", len(bodyBytes))
		pctx[fmt.Sprintf("probe.%s.req.%s.resp.body.size", probe.Name, probeRequest.Name)] = len(bodyBytes)
		if hVals := response.Header.Values("Content-Type"); hVals != nil {
			if strings.Contains(hVals[0], "text/") {
				requestLog.Tracef("Http response body is [%s]", string(bodyBytes))
				pctx[fmt.Sprintf("probe.%s.req.%s.resp.body", probe.Name, probeRequest.Name)] = string(bodyBytes)
			} else {
				requestLog.Tracef("Http response body is binary")
				pctx[fmt.Sprintf("probe.%s.req.%s.resp.body", probe.Name, probeRequest.Name)] = "<binary>"
			}
		}
	}

	// Check success if
	if len(probeRequest.SuccessIf) > 0 {
		requestLog.Tracef("Evaluating SuccessIf [%s]", probeRequest.SuccessIf)
		out, err := GoCelEvaluate(ctx, probeRequest.SuccessIf, pctx, reflect.Bool)
		if err != nil {
			pctx[fmt.Sprintf("probe.%s.req.%s.success", probe.Name, probeRequest.Name)] = false
			pctx[fmt.Sprintf("probe.%s.req.%s.fail", probe.Name, probeRequest.Name)] = true
			pctx[fmt.Sprintf("probe.%s.req.%s.error", probe.Name, probeRequest.Name)] = err
			requestLog.Errorf("error when evaluating SuccessIf [%s]. got %s", probeRequest.SuccessIf, err.Error())
			return fmt.Errorf("%w : probe %s request %s parsing SuccessIf parsing error [%s]", err, probe.Name, probeRequest.Name, probeRequest.SuccessIf)
		}
		if !out.(bool) {
			requestLog.Errorf("evaluation of SuccessIf [%s] yields a %v", probeRequest.SuccessIf, out.(bool))
			pctx[fmt.Sprintf("probe.%s.req.%s.success", probe.Name, probeRequest.Name)] = false
			pctx[fmt.Sprintf("probe.%s.req.%s.fail", probe.Name, probeRequest.Name)] = true
			return fmt.Errorf("%w : probe %s request %s SuccessIf criteria returns false", errors.ErrSuccessIfIsFalse, probe.Name, probeRequest.Name)
		}
		pctx[fmt.Sprintf("probe.%s.req.%s.success", probe.Name, probeRequest.Name)] = true
		pctx[fmt.Sprintf("probe.%s.req.%s.fail", probe.Name, probeRequest.Name)] = false
	} else if len(probeRequest.FailIf) > 0 { // Check fail if
		requestLog.Tracef("Evaluating FailIf [%s]", probeRequest.FailIf)
		out, err := GoCelEvaluate(ctx, probeRequest.FailIf, pctx, reflect.Bool)
		if err != nil {
			pctx[fmt.Sprintf("probe.%s.req.%s.fail", probe.Name, probeRequest.Name)] = true
			pctx[fmt.Sprintf("probe.%s.req.%s.success", probe.Name, probeRequest.Name)] = false
			pctx[fmt.Sprintf("probe.%s.req.%s.error", probe.Name, probeRequest.Name)] = err
			requestLog.Errorf("error when evaluating FailIf. got %s", err.Error())
			return fmt.Errorf("%w : probe %s request %s parsing FailIf parsing error [%s]", err, probe.Name, probeRequest.Name, probeRequest.FailIf)
		}
		if out.(bool) {
			requestLog.Errorf("evaluation of FailIf [%s] yields a %v", probeRequest.FailIf, out.(bool))
			pctx[fmt.Sprintf("probe.%s.req.%s.fail", probe.Name, probeRequest.Name)] = true
			pctx[fmt.Sprintf("probe.%s.req.%s.success", probe.Name, probeRequest.Name)] = false
			return fmt.Errorf("%w : probe %s request %s FailIf criteria returns true", errors.ErrFailIfIsTrue, probe.Name, probeRequest.Name)
		}
		pctx[fmt.Sprintf("probe.%s.req.%s.fail", probe.Name, probeRequest.Name)] = false
		pctx[fmt.Sprintf("probe.%s.req.%s.success", probe.Name, probeRequest.Name)] = true
	}
	pctx[fmt.Sprintf("probe.%s.req.%s.fail", probe.Name, probeRequest.Name)] = false
	pctx[fmt.Sprintf("probe.%s.req.%s.success", probe.Name, probeRequest.Name)] = true

	return nil
}
