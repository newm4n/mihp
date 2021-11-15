package notification

const (
	NotifTypeEmailSMTP = "SMTP"
	NotifTypeCallBack  = "CALLBACK"
	NotifTypeTelegram  = "TELEGRAM"
	NotifTypeSlack     = "SLACK"
	NotifTypeMSTeams   = "MSTEAMS"

	EventUp EventType = iota
	EventDown
)

type EventType int

type Notification interface {
	Notify() error
}

var (
	notificationChannel chan Notification
)

func init() {
	notificationChannel = make(chan Notification)
}
