package probing

import (
	"fmt"
	"github.com/newm4n/mihp/pkg/helper"
	"strconv"
	"strings"
)

const (
	version = "1.0.0"
)

type ProbePool []*Probe

func (pool ProbePool) FromProperties(prop helper.Properties) error {
	if ver, ok := prop["version"]; !ok || ver != version {
		return fmt.Errorf("invalid version %s", prop["version"])
	}
	pool = pool[:]
	probeCount := 0
	if strPCount, ok := prop["probe.count"]; !ok {
		return fmt.Errorf("properties contains no probe")
	} else {
		i, err := strconv.Atoi(strPCount)
		if err != nil {
			return fmt.Errorf("probe.count is not an integer")
		}
		probeCount = i
	}

	for probeIdx := 0; probeIdx < probeCount; probeIdx++ {
		probe := &Probe{}

		if strProbName, ok := prop[fmt.Sprintf("probe.%d.name", probeIdx)]; ok {
			probe.Name = strProbName
		} else {
			return fmt.Errorf("probe.%d has no name", probeIdx)
		}

		if strProbID, ok := prop[fmt.Sprintf("probe.%d.id", probeIdx)]; ok {
			probe.ID = strProbID
		} else {
			return fmt.Errorf("probe.%d has no id", probeIdx)
		}

		if strProbCron, ok := prop[fmt.Sprintf("probe.%d.cron", probeIdx)]; ok {
			probe.Cron = strProbCron
		} else {
			return fmt.Errorf("probe.%d has no id", probeIdx)
		}

		// TODO continue with threshold

		pool = append(pool, probe)
	}
	return nil
}

func (pool ProbePool) ToProperties() helper.Properties {
	prop := helper.NewProperties()
	prop["version"] = version

	prop["probe.count"] = fmt.Sprintf("%d", len(pool))
	for idx, prob := range pool {
		prop[fmt.Sprintf("probe.%d.name", idx)] = prob.Name
		prop[fmt.Sprintf("probe.%d.id", idx)] = prob.ID
		prop[fmt.Sprintf("probe.%d.cron", idx)] = prob.Cron
		prop[fmt.Sprintf("probe.%d.threshold.up", idx)] = fmt.Sprintf("%d", prob.UpThreshold)
		prop[fmt.Sprintf("probe.%d.threshold.down", idx)] = fmt.Sprintf("%d", prob.UpThreshold)
		if prob.SMTPNotification != nil {
			prop[fmt.Sprintf("probe.%d.notification.smtp.host", idx)] = prob.SMTPNotification.SMTPHost
			prop[fmt.Sprintf("probe.%d.notification.smtp.port", idx)] = fmt.Sprintf("%d", prob.SMTPNotification.SMTPPort)
			prop[fmt.Sprintf("probe.%d.notification.smtp.from.email", idx)] = prob.SMTPNotification.From.Email
			prop[fmt.Sprintf("probe.%d.notification.smtp.from.name", idx)] = prob.SMTPNotification.From.Name
			if prob.SMTPNotification.To != nil {
				if len(prob.SMTPNotification.To) > 0 {
					prop[fmt.Sprintf("probe.%d.notification.smtp.to.count", idx)] = fmt.Sprintf("%d", len(prob.SMTPNotification.To))
					for eidx, mbx := range prob.SMTPNotification.To {
						prop[fmt.Sprintf("probe.%d.notification.smtp.to.%d.email", idx, eidx)] = mbx.Email
						prop[fmt.Sprintf("probe.%d.notification.smtp.to.%d.name", idx, eidx)] = mbx.Name
					}
				}
			}
			if prob.SMTPNotification.Cc != nil {
				if len(prob.SMTPNotification.Cc) > 0 {
					prop[fmt.Sprintf("probe.%d.notification.smtp.cc.count", idx)] = fmt.Sprintf("%d", len(prob.SMTPNotification.Cc))
					for eidx, mbx := range prob.SMTPNotification.Cc {
						prop[fmt.Sprintf("probe.%d.notification.smtp.cc.%d.email", idx, eidx)] = mbx.Email
						prop[fmt.Sprintf("probe.%d.notification.smtp.cc.%d.name", idx, eidx)] = mbx.Name
					}
				}
			}
			if prob.SMTPNotification.Bcc != nil {
				if len(prob.SMTPNotification.Bcc) > 0 {
					prop[fmt.Sprintf("probe.%d.notification.smtp.bcc.count", idx)] = fmt.Sprintf("%d", len(prob.SMTPNotification.Bcc))
					for eidx, mbx := range prob.SMTPNotification.Bcc {
						prop[fmt.Sprintf("probe.%d.notification.smtp.bcc.%d.email", idx, eidx)] = mbx.Email
						prop[fmt.Sprintf("probe.%d.notification.smtp.bcc.%d.name", idx, eidx)] = mbx.Name
					}
				}
			}
		}
		if prob.CallbackNotification != nil {
			prop[fmt.Sprintf("probe.%d.notification.callback.url.up", idx)] = prob.CallbackNotification.UpCall
			prop[fmt.Sprintf("probe.%d.notification.callback.url.down", idx)] = prob.CallbackNotification.DownCall
		}
		if prob.Requests != nil {
			prop[fmt.Sprintf("probe.%d.request.count", idx)] = fmt.Sprintf("%d", len(prob.Requests))
			for ridx, req := range prob.Requests {
				prop[fmt.Sprintf("probe.%d.request.%d.name", idx, ridx)] = req.Name
				prop[fmt.Sprintf("probe.%d.request.%d.url", idx, ridx)] = req.URL
				prop[fmt.Sprintf("probe.%d.request.%d.method", idx, ridx)] = req.Method
				prop[fmt.Sprintf("probe.%d.request.%d.body", idx, ridx)] = req.Body
				if len(req.Headers) > 0 {
					hNames := make([]string, 0)
					for hk, _ := range req.Headers {
						hNames = append(hNames, hk)
					}
					prop[fmt.Sprintf("probe.%d.request.%d.header", idx, ridx)] = strings.Join(hNames, ",")
					for hk, hkVals := range req.Headers {
						prop[fmt.Sprintf("probe.%d.request.%d.header.%s.count", idx, ridx, hk)] = fmt.Sprintf("%d", len(hkVals))
						for hvIdx, hvs := range hkVals {
							prop[fmt.Sprintf("probe.%d.request.%d.header.%s.%d", idx, ridx, hk, hvIdx)] = hvs
						}
					}
				}
				prop[fmt.Sprintf("probe.%d.request.%d.startRequestIf", idx, ridx)] = req.StartRequestIf
				prop[fmt.Sprintf("probe.%d.request.%d.SuccessIf", idx, ridx)] = req.SuccessIf
				prop[fmt.Sprintf("probe.%d.request.%d.FailIf", idx, ridx)] = req.FailIf
				prop[fmt.Sprintf("probe.%d.request.%d.CertificateCheck", idx, ridx)] = fmt.Sprintf("%v", req.CertificateCheck)
			}
		}
	}
	return prop
}

type Probe struct {
	Name                 string                      `json:"name"`
	ID                   string                      `json:"id"`
	Requests             []*ProbeRequest             `json:"requests"`
	Cron                 string                      `json:"cron"`
	UpThreshold          int                         `json:"up_threshold"`
	DownThreshold        int                         `json:"down_threshold"`
	SMTPNotification     *SMTPNotificationTarget     `json:"smtp_notification"`
	CallbackNotification *CallbackNotificationTarget `json:"callback_notification"`
}

type SMTPNotificationTarget struct {
	SMTPHost string
	SMTPPort int
	From     *Mailbox
	Password string
	To       []*Mailbox
	Cc       []*Mailbox
	Bcc      []*Mailbox
}

type CallbackNotificationTarget struct {
	UpCall   string
	DownCall string
}

type Mailbox struct {
	Name  string
	Email string
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
