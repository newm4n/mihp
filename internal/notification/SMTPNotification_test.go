package notification

import (
	"github.com/newm4n/mihp/internal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSMTPNotification_SendNotification(t *testing.T) {
	from := &internal.Mailbox{
		Name:  "Dummy Sender",
		Email: "dummysender@mihp.com",
	}
	notif := &SMTPNotification{
		EmailNotification: EmailNotification{
			FromField: from,
			ToList: []*internal.Mailbox{
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
