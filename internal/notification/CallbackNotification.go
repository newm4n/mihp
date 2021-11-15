package notification

import (
	"github.com/newm4n/mihp/internal/probing"
	"net/http"
	"time"
)

type CallbackNotification struct {
	UpURL     string
	DownURL   string
	EventType EventType
}

func (notif *CallbackNotification) Notify() error {
	client := probing.NewHttpClient(10, 10, true)
	if notif.EventType == EventUp {
		req, _ := http.NewRequest("GET", notif.UpURL, nil)
		_, err := client.Do(req)
		return err
	} else {
		req, _ := http.NewRequest("GET", notif.DownURL, nil)
		_, err := client.Do(req)
		return err
	}
}

func (notif *CallbackNotification) CallbackNotificationTrigger(probeName, probeId string, down bool, firstUp, lastUp, firstDown, lastDown time.Time) {

}
