package notification

import (
	"bytes"
	"embed"
	"github.com/sirupsen/logrus"
	"html/template"
	"log"
	"strings"
)

//go:embed static
var staticFolder embed.FS

const (
	mime = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
)

var (
	upMailBodyTmpl   *template.Template
	downMailBodyTmpl *template.Template
	upSubjectTmpl    *template.Template
	downSubjectTmpl  *template.Template
)

func init() {
	tmpl, err := template.ParseFS(staticFolder, "static/smtp_up_body.html")
	if err != nil {
		log.Fatal(err)
	}
	upMailBodyTmpl = tmpl
	tmpl, err = template.ParseFS(staticFolder, "static/smtp_down_body.html")
	if err != nil {
		log.Fatal(err)
	}
	downMailBodyTmpl = tmpl
	tmpl, err = template.ParseFS(staticFolder, "static/smtp_up_subject.txt")
	if err != nil {
		log.Fatal(err)
	}
	upSubjectTmpl = tmpl
	tmpl, err = template.ParseFS(staticFolder, "static/smtp_down_subject.txt")
	if err != nil {
		log.Fatal(err)
	}
	downSubjectTmpl = tmpl
}

type EmailNotification struct {
	FromField string `json:"from_field"`

	ToList  []*Recipient `json:"to_list"`
	CcList  []*Recipient `json:"cc_list"`
	BccList []*Recipient `json:"bcc_list"`
}

func (notif *EmailNotification) Headers(subject string) string {
	var bodyBuffer bytes.Buffer
	if notif.ToList != nil && len(notif.ToList) > 0 {
		bodyBuffer.WriteString("To: ")
		count := 0
		for _, r := range notif.ToList {
			if count > 0 {
				bodyBuffer.WriteString(",")
			}
			bodyBuffer.WriteString(r.String())
			count++
		}
		bodyBuffer.WriteString("\r\n")
	}
	if notif.CcList != nil && len(notif.CcList) > 0 {
		bodyBuffer.WriteString("Cc: ")
		count := 0
		for _, r := range notif.CcList {
			if count > 0 {
				bodyBuffer.WriteString(",")
			}
			bodyBuffer.WriteString(r.String())
			count++
		}
		bodyBuffer.WriteString("\r\n")
	}
	bodyBuffer.WriteString("Subject: ")
	bodyBuffer.WriteString(subject)
	bodyBuffer.WriteString("\r\n")
	return bodyBuffer.String()
}

func (notif *EmailNotification) Receivers() []string {
	recList := make(map[string]bool)
	if notif.ToList != nil && len(notif.ToList) > 0 {
		for _, r := range notif.ToList {
			recList[r.Email] = true
		}
	}
	if notif.CcList != nil && len(notif.CcList) > 0 {
		for _, r := range notif.CcList {
			recList[r.Email] = true
		}
	}
	if notif.BccList != nil && len(notif.BccList) > 0 {
		for _, r := range notif.BccList {
			recList[r.Email] = true
		}
	}
	receivers := make([]string, 0)
	for eml, _ := range recList {
		receivers = append(receivers, eml)
	}
	return receivers
}

func (notif *EmailNotification) SendNotification(subject, body string) error {
	sendmailLog := logrus.WithField("mailer", "sendmail").WithField("from", notif.FromField).WithField("mode", "DUMMY")

	sendmailLog.Infof("Using PlainAuth u=%s p=******", notif.FromField)

	var bodyBuffer bytes.Buffer

	bodyBuffer.WriteString(notif.Headers(subject))
	bodyBuffer.WriteString("\r\n")
	bodyBuffer.WriteString(mime)
	bodyBuffer.WriteString("\r\n")
	bodyBuffer.WriteString(body)

	receivers := notif.Receivers()
	sendingLog := sendmailLog.WithField("to", strings.Join(receivers, ","))
	sendingLog.Infof("sending email Body ... \n%s", bodyBuffer.String())
	sendingLog.Info("send DUMMY mail success")
	return nil
}
