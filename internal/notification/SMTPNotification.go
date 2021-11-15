package notification

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/smtp"
	"strings"
)

func NewSMTPDownNotification() *SMTPNotification {
	return &SMTPNotification{
		EmailNotification: EmailNotification{
			FromField: nil,
			ToList:    nil,
			CcList:    nil,
			BccList:   nil,
		},
		EventType:     EventDown,
		PasswordField: "",
		ProbeName:     "",
		Cause:         "",
		UpDuration:    "",
		DownDuration:  "",
		SMTPHost:      "",
		SMTPPort:      0,
	}
}

func NewSMTPUpNotification() *SMTPNotification {
	return &SMTPNotification{
		EmailNotification: EmailNotification{
			FromField: nil,
			ToList:    nil,
			CcList:    nil,
			BccList:   nil,
		},
		EventType:     EventUp,
		PasswordField: "",
		ProbeName:     "",
		Cause:         "",
		UpDuration:    "",
		DownDuration:  "",
		SMTPHost:      "",
		SMTPPort:      0,
	}
}

type SMTPNotification struct {
	EmailNotification

	EventType     EventType `json:"event_type"`
	PasswordField string    `json:"password_field"`

	ProbeName    string
	Cause        string
	UpDuration   string
	DownDuration string

	SMTPHost string `json:"smtp_host"`
	SMTPPort int    `json:"smtp_port"`
}

func (notif *SMTPNotification) SendNotification(subject, body string) error {
	sendmailLog := logrus.WithField("mailer", "sendmail").WithField("from", notif.FromField)

	auth := smtp.PlainAuth("", notif.FromField.Email, notif.PasswordField, notif.SMTPHost)

	var bodyBuffer bytes.Buffer

	bodyBuffer.WriteString(notif.Headers(subject))
	bodyBuffer.WriteString("\r\n")
	bodyBuffer.WriteString(mime)
	bodyBuffer.WriteString("\r\n")
	bodyBuffer.WriteString(body)

	receivers := notif.Receivers()
	sendingLog := sendmailLog.WithField("to", strings.Join(receivers, ","))

	sendingLog.Debugf("sending using server %s:%d > BodyExpr ... \n%s", notif.SMTPHost, notif.SMTPPort, bodyBuffer.String())
	err := smtp.SendMail(fmt.Sprintf("%s:%d", notif.SMTPHost, notif.SMTPPort), auth, notif.FromField.Email, receivers, bodyBuffer.Bytes())
	if err != nil {
		sendingLog.Error(err)
		return err
	}
	sendingLog.Debug("send SMTP mail success")
	return nil
}

func (notif *SMTPNotification) Notify() error {
	if notif.EventType == EventUp {
		return notif.NotifyUp(notif.ProbeName, notif.DownDuration)
	} else {
		return notif.NotifyDown(notif.ProbeName, notif.Cause, notif.UpDuration)
	}
}

func (notif *SMTPNotification) NotifyDown(probeName, cause, upDuration string) error {
	data := map[string]string{
		"probe":      probeName,
		"cause":      cause,
		"upDuration": upDuration,
	}
	subjectbuff := &bytes.Buffer{}
	err := downSubjectTmpl.Execute(subjectbuff, data)
	if err != nil {
		return err
	}

	bodybuff := &bytes.Buffer{}
	err = downMailBodyTmpl.Execute(bodybuff, data)
	if err != nil {
		return err
	}

	return notif.SendNotification(subjectbuff.String(), bodybuff.String())
}

func (notif *SMTPNotification) NotifyUp(probeName, downDuration string) error {
	data := map[string]string{
		"probe":        probeName,
		"downDuration": downDuration,
	}
	subjectbuff := &bytes.Buffer{}
	err := upSubjectTmpl.Execute(subjectbuff, data)
	if err != nil {
		return err
	}

	bodybuff := &bytes.Buffer{}
	err = upMailBodyTmpl.Execute(bodybuff, data)
	if err != nil {
		return err
	}

	return notif.SendNotification(subjectbuff.String(), bodybuff.String())
}
