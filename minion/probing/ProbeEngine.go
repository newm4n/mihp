package probing

import (
	"bytes"
	"context"
	"fmt"
	"github.com/newm4n/mihp/internal"
	"github.com/newm4n/mihp/pkg/errors"
	"github.com/newm4n/mihp/pkg/helper/cron"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"
)

const (
	TimeoutsSecond = 10
)

var (
	engineLog = logrus.WithField("module", "ProbeEngine")
)

func ProbeCanStartBySchedule(ctx context.Context, probe *internal.Probe) bool {
	if ctx.Err() != nil {
		engineLog.Warnf("Will never ever start since context has error. got %s", ctx.Err())
		return false
	}
	cs, err := cron.NewSchedule(probe.Cron)
	if err != nil {
		engineLog.Errorf("Will never ever start since cron is wrong. got %s", err.Error())
		return false
	}
	return cs.IsIn(time.Now())
}

func ExecuteProbe(ctx context.Context, probe *internal.Probe, pctx internal.ProbeContext, timeoutSecond int, ignoreTLS, ignoreSchedule bool) error {
	probeLog := engineLog.WithField("probe", probe.Name)
	if ctx.Err() != nil {
		probeLog.Errorf("context error. got %s", ctx.Err())
		return fmt.Errorf("%w : context probably timed-out. got %s", errors.ErrContextError, ctx.Err())
	}
	pctx["probe"] = probe.Name
	if ProbeCanStartBySchedule(ctx, probe) || ignoreSchedule {
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
			err := ExecuteProbeRequest(ctx, probe, reqs, seq, timeoutSecond, ignoreTLS, pctx)
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

func ExecuteProbeRequest(ctx context.Context, probe *internal.Probe, probeRequest *internal.ProbeRequest,
	sequence int, timeoutSecond int, ignoreTLS bool, pctx internal.ProbeContext) error {

	requestLog := engineLog.WithField("probe", probe.Name).WithField("request", probeRequest.Name)

	if ctx.Err() != nil {
		pctx[fmt.Sprintf("probe.%s.req.%s.error", probe.Name, probeRequest.Name)] = ctx.Err()
		requestLog.Errorf("context error. got %s", ctx.Err())
		return fmt.Errorf("%w : context probably timed-out. got %s", errors.ErrContextError, ctx.Err())
	}
	pctx[fmt.Sprintf("probe.%s.req.%s.sequence", probe.Name, probeRequest.Name)] = sequence
	if len(probeRequest.StartRequestIfExpr) > 0 {
		out, err := GoCelEvaluate(ctx, probeRequest.StartRequestIfExpr, pctx, reflect.Bool)
		if err != nil {
			requestLog.Errorf("error during evaluating StartRequestIfExpr. got %s", err.Error())
			pctx[fmt.Sprintf("probe.%s.req.%s.canstart", probe.Name, probeRequest.Name)] = false
			pctx[fmt.Sprintf("probe.%s.req.%s.error", probe.Name, probeRequest.Name)] = err
			return fmt.Errorf("%w : probe %s request %s parsing StartRequestIfExpr parsing error [%s]", err, probe.Name, probeRequest.Name, probeRequest.StartRequestIfExpr)
		}
		if !out.(bool) {
			requestLog.Tracef("evaluation of StartRequestIfExpr [%s] says that probe should not start", probeRequest.StartRequestIfExpr)
			pctx[fmt.Sprintf("probe.%s.req.%s.canstart", probe.Name, probeRequest.Name)] = false
			return fmt.Errorf("%w : probe %s request %s can not start", errors.ErrStartRequestIfIsFalse, probe.Name, probeRequest.Name)
		}
	}
	requestLog.Tracef("StartRequestIfExpr is OK")
	pctx[fmt.Sprintf("probe.%s.req.%s.canstart", probe.Name, probeRequest.Name)] = true

	requestLog.Tracef("Retriefing HTTP client with %d second timeout", TimeoutsSecond)
	client := NewHttpClient(timeoutSecond, timeoutSecond, ignoreTLS)
	var request *http.Request
	var err error

	requestLog.Tracef("Evaluating PathExpr [%s]", probeRequest.PathExpr)
	urlItv, err := GoCelEvaluate(ctx, probeRequest.PathExpr, pctx, reflect.String)
	if err != nil {
		requestLog.Errorf("Error evaluating PathExpr [%s] got %s", probeRequest.PathExpr, err.Error())
		pctx[fmt.Sprintf("probe.%s.req.%s.error", probe.Name, probeRequest.Name)] = err
		return fmt.Errorf("%w : error while parsing URLExpr", err)
	}
	URL := fmt.Sprintf("%s%s", probe.BaseURL, urlItv.(string))

	requestLog.Tracef("URLExpr [%s] evaluated as [%s]", probeRequest.PathExpr, URL)
	pctx[fmt.Sprintf("probe.%s.req.%s.url", probe.Name, probeRequest.Name)] = URL

	requestLog.Tracef("Evaluating MethodExpr [%s]", probeRequest.MethodExpr)
	methodItv, err := GoCelEvaluate(ctx, probeRequest.MethodExpr, pctx, reflect.String)
	if err != nil {
		requestLog.Errorf("Error evaluating MethodExpr [%s] got %s", probeRequest.MethodExpr, err.Error())
		pctx[fmt.Sprintf("probe.%s.req.%s.error", probe.Name, probeRequest.Name)] = err
		return fmt.Errorf("%w : error while parsing METHOD", err)
	}
	METHOD := methodItv.(string)
	requestLog.Tracef("MethodExpr [%s] evaluated as [%s]", probeRequest.MethodExpr, METHOD)
	pctx[fmt.Sprintf("probe.%s.req.%s.method", probe.Name, probeRequest.Name)] = METHOD

	if probeRequest.BodyExpr != "" {
		requestLog.Tracef("Evaluating BodyExpr [%s]", probeRequest.BodyExpr)
		bodyItv, err := GoCelEvaluate(ctx, probeRequest.BodyExpr, pctx, reflect.String)
		if err != nil {
			requestLog.Errorf("Error evaluating BodyExpr [%s] got %s", probeRequest.BodyExpr, err.Error())
			pctx[fmt.Sprintf("probe.%s.req.%s.error", probe.Name, probeRequest.Name)] = err
			return fmt.Errorf("%w : error while parsing Request BodyExpr", err)
		}
		pctx[fmt.Sprintf("probe.%s.req.%s.body", probe.Name, probeRequest.Name)] = bodyItv.(string)

		requestLog.Tracef("BodyExpr [%s] evaluated as [%s]", probeRequest.BodyExpr, bodyItv.(string))

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

	if probeRequest.HeadersExpr != nil && len(probeRequest.HeadersExpr) > 0 {
		requestLog.Tracef("Parsing %d request headers", len(probeRequest.HeadersExpr))
		headerKeys := make([]string, 0)
		for hKey, _ := range probeRequest.HeadersExpr {
			headerKeys = append(headerKeys, hKey)
		}
		pctx[fmt.Sprintf("probe.%s.req.%s.header", probe.Name, probeRequest.Name)] = headerKeys
		for hKey, hVals := range probeRequest.HeadersExpr {
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
	if len(probeRequest.SuccessIfExpr) > 0 {
		requestLog.Tracef("Evaluating SuccessIfExpr [%s]", probeRequest.SuccessIfExpr)
		out, err := GoCelEvaluate(ctx, probeRequest.SuccessIfExpr, pctx, reflect.Bool)
		if err != nil {
			pctx[fmt.Sprintf("probe.%s.req.%s.success", probe.Name, probeRequest.Name)] = false
			pctx[fmt.Sprintf("probe.%s.req.%s.fail", probe.Name, probeRequest.Name)] = true
			pctx[fmt.Sprintf("probe.%s.req.%s.error", probe.Name, probeRequest.Name)] = err
			requestLog.Errorf("error when evaluating SuccessIfExpr [%s]. got %s", probeRequest.SuccessIfExpr, err.Error())
			return fmt.Errorf("%w : probe %s request %s parsing SuccessIfExpr parsing error [%s]", err, probe.Name, probeRequest.Name, probeRequest.SuccessIfExpr)
		}
		if !out.(bool) {
			requestLog.Errorf("evaluation of SuccessIfExpr [%s] yields a %v", probeRequest.SuccessIfExpr, out.(bool))
			pctx[fmt.Sprintf("probe.%s.req.%s.success", probe.Name, probeRequest.Name)] = false
			pctx[fmt.Sprintf("probe.%s.req.%s.fail", probe.Name, probeRequest.Name)] = true
			return fmt.Errorf("%w : probe %s request %s SuccessIfExpr criteria returns false", errors.ErrSuccessIfIsFalse, probe.Name, probeRequest.Name)
		}
		pctx[fmt.Sprintf("probe.%s.req.%s.success", probe.Name, probeRequest.Name)] = true
		pctx[fmt.Sprintf("probe.%s.req.%s.fail", probe.Name, probeRequest.Name)] = false
	} else if len(probeRequest.FailIfExpr) > 0 { // Check fail if
		requestLog.Tracef("Evaluating FailIfExpr [%s]", probeRequest.FailIfExpr)
		out, err := GoCelEvaluate(ctx, probeRequest.FailIfExpr, pctx, reflect.Bool)
		if err != nil {
			pctx[fmt.Sprintf("probe.%s.req.%s.fail", probe.Name, probeRequest.Name)] = true
			pctx[fmt.Sprintf("probe.%s.req.%s.success", probe.Name, probeRequest.Name)] = false
			pctx[fmt.Sprintf("probe.%s.req.%s.error", probe.Name, probeRequest.Name)] = err
			requestLog.Errorf("error when evaluating FailIfExpr. got %s", err.Error())
			return fmt.Errorf("%w : probe %s request %s parsing FailIfExpr parsing error [%s]", err, probe.Name, probeRequest.Name, probeRequest.FailIfExpr)
		}
		if out.(bool) {
			requestLog.Errorf("evaluation of FailIfExpr [%s] yields a %v", probeRequest.FailIfExpr, out.(bool))
			pctx[fmt.Sprintf("probe.%s.req.%s.fail", probe.Name, probeRequest.Name)] = true
			pctx[fmt.Sprintf("probe.%s.req.%s.success", probe.Name, probeRequest.Name)] = false
			return fmt.Errorf("%w : probe %s request %s FailIfExpr criteria returns true", errors.ErrFailIfIsTrue, probe.Name, probeRequest.Name)
		}
		pctx[fmt.Sprintf("probe.%s.req.%s.fail", probe.Name, probeRequest.Name)] = false
		pctx[fmt.Sprintf("probe.%s.req.%s.success", probe.Name, probeRequest.Name)] = true
	}
	pctx[fmt.Sprintf("probe.%s.req.%s.fail", probe.Name, probeRequest.Name)] = false
	pctx[fmt.Sprintf("probe.%s.req.%s.success", probe.Name, probeRequest.Name)] = true

	return nil
}
