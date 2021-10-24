package notification

type Notification interface {
	NotifyUp(probeName, downDuration string) error
	NotifyDown(probeName, cause, upDuration string) error
}
