package notification

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSMTPNotification_SendNotification(t *testing.T) {
	notif := &SMTPNotification{
		EmailNotification: EmailNotification{
			FromField: "dummysender@mihp.com",
			ToList: []*Recipient{
				{
					Name:  "Dummy Target Notif",
					Email: "dummytargetnotif@mihp.com",
				},
			},
			CcList:  nil,
			BccList: nil,
		},
		PasswordField: "",
		SMTPHost:      "mail.smtpbucket.com",
		SMTPPort:      8025,
	}

	err := notif.NotifyUp("dummyprobe", "100 minutes")
	assert.NoError(t, err)
}
