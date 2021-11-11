package probing

import (
	yaml "gopkg.in/yaml.v3"
)

const (
	version = "1.0.0"
)

type ProbePool []*Probe

type Probe struct {
	Name                 string                      `json:"name"`
	ID                   string                      `json:"id"`
	Requests             []*ProbeRequest             `json:"requests"`
	BaseURL              string                      `json:"base_url"`
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

func (this *Mailbox) Equals(that *Mailbox) bool {
	if that == nil {
		return false
	}
	if this.Name != that.Name || this.Email != that.Email {
		return false
	}
	return true
}

type ProbeRequest struct {
	Name                 string              `json:"name"`
	PathExpr             string              `json:"path_expr"`
	MethodExpr           string              `json:"method_expr"`
	HeadersExpr          map[string][]string `json:"headers_exprs"`
	BodyExpr             string              `json:"body_expr"`
	CertificateCheckExpr string              `json:"certificate_check_expr"`
	StartRequestIfExpr   string              `json:"start_request_if_expr"`
	SuccessIfExpr        string              `json:"success_if_expr"`
	FailIfExpr           string              `json:"fail_if_expr"`
}

func YAMLToProbePool(yamlBytes []byte) (probePool ProbePool, err error) {
	pool := make(ProbePool, 0)
	err = yaml.Unmarshal(yamlBytes, &pool)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func ProbePoolToYAML(pool ProbePool) (yamlBytes []byte, err error) {
	return yaml.Marshal(pool)
}
