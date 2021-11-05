package notification

const (
	NotifTypeEmailSMTP = "SMTP"
	NotifTypeCallBack  = "CALLBACK"
	NotifTypeTelegram  = "TELEGRAM"
	NotifTypeSlack     = "SLACK"
	NotifTypeMSTeams   = "MSTEAMS"
)

type Notification interface {
	NotifyUp(probeName, downDuration string) error
	NotifyDown(probeName, cause, upDuration string) error
}
