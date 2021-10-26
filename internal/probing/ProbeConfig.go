package probing

type ProbePool []*Probe

type Probe struct {
	Name                         string          `json:"name"`
	ID                           string          `json:"id"`
	Requests                     []*ProbeRequest `json:"requests"`
	Cron                         string          `json:"cron"`
	SuccessNotificationThreshold int             `json:"success_notification_threshold"`
	FailNotificationThreshold    int             `json:"fail_notification_threshold"`
	NotificationType             string          `json:"notification_type"`
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
