package notification

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/smtp"
	"strings"
)

type SMTPNotification struct {
	EmailNotification
	PasswordField string `json:"password_field"`

	SMTPHost string `json:"smtp_host"`
	SMTPPort int    `json:"smtp_port"`
}

func (notif *SMTPNotification) SendNotification(subject, body string) error {
	sendmailLog := logrus.WithField("mailer", "sendmail").WithField("from", notif.FromField)

	auth := smtp.PlainAuth("", notif.FromField, notif.PasswordField, notif.SMTPHost)

	var bodyBuffer bytes.Buffer

	bodyBuffer.WriteString(notif.Headers(subject))
	bodyBuffer.WriteString("\r\n")
	bodyBuffer.WriteString(mime)
	bodyBuffer.WriteString("\r\n")
	bodyBuffer.WriteString(body)

	receivers := notif.Receivers()
	sendingLog := sendmailLog.WithField("to", strings.Join(receivers, ","))

	sendingLog.Debugf("sending using server %s:%d > BodyExpr ... \n%s", notif.SMTPHost, notif.SMTPPort, bodyBuffer.String())
	err := smtp.SendMail(fmt.Sprintf("%s:%d", notif.SMTPHost, notif.SMTPPort), auth, notif.FromField, receivers, bodyBuffer.Bytes())
	if err != nil {
		sendingLog.Error(err)
		return err
	}
	sendingLog.Debug("send SMTP mail success")
	return nil
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
