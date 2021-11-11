package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProbePoolSerialization(t *testing.T) {
	var yamlByte []byte
	var pool ProbePool
	pool = make(ProbePool, 0)
	probe := &Probe{
		Name: "ProbeMaster",
		ID:   "1234-56789",
		Requests: []*ProbeRequest{
			{
				Name:       "request1",
				PathExpr:   "/req1",
				MethodExpr: "\"GET\"",
				HeadersExpr: map[string][]string{
					"Content-Type": []string{"\"application/yaml\""},
					"User-Agent":   []string{"\"someuser agent/123\""},
				},
				BodyExpr:             "",
				CertificateCheckExpr: "true",
				StartRequestIfExpr:   "IsDefined(\"some.property.in.context.start\")",
				SuccessIfExpr:        "IsDefined(\"some.property.in.context\") && GetInt(\"some.property.code\") == 200",
				FailIfExpr:           "IsDefined(\"some.property.in.context\") && GetInt(\"some.property.code\") != 200",
			},
			{
				Name:       "request2",
				PathExpr:   "/req2",
				MethodExpr: "\"POST\"",
				HeadersExpr: map[string][]string{
					"Content-Type": []string{"\"application/yaml\""},
					"User-Agent":   []string{"\"someuser agent/123\""},
				},
				BodyExpr:             "\"this is the body of something\"",
				CertificateCheckExpr: "true",
				StartRequestIfExpr:   "IsDefined(\"some.property.in.context.start\")",
				SuccessIfExpr:        "IsDefined(\"some.property.in.context\") && GetInt(\"some.property.code\") == 200",
				FailIfExpr:           "IsDefined(\"some.property.in.context\") && GetInt(\"some.property.code\") != 200",
			},
		},
		BaseURL:       "https://baseURL:1234",
		Cron:          "0 */3 * * * * *",
		UpThreshold:   4,
		DownThreshold: 4,
		SMTPNotification: &SMTPNotificationTarget{
			SMTPHost: "smpt.google.com",
			SMTPPort: 32,
			From: &Mailbox{
				Name:  "Domain Owner",
				Email: "owner@domain.com",
			},
			Password: "smtppassword",
			To: []*Mailbox{
				{
					Name:  "Domain Owner 1",
					Email: "owner1@domain.com",
				},
				{
					Name:  "Domain Owner 2",
					Email: "owner2@domain.com",
				},
			},
			Cc: []*Mailbox{
				{
					Name:  "Domain Owner 3",
					Email: "owner3@domain.com",
				},
				{
					Name:  "Domain Owner 4",
					Email: "owner4@domain.com",
				},
			},
			Bcc: []*Mailbox{
				{
					Name:  "Domain Owner 5",
					Email: "owner5@domain.com",
				},
				{
					Name:  "Domain Owner 6",
					Email: "owner6@domain.com",
				},
			},
		},
		CallbackNotification: &CallbackNotificationTarget{
			UpCall:   "http://targetmonitor.com/up",
			DownCall: "http://targetmonitor.com/down",
		},
	}
	pool = append(pool, probe)

	yamlBytes, err := ProbePoolToYAML(pool)
	assert.NoError(t, err)
	yamlByte = yamlBytes
	t.Log("\n" + string(yamlByte))

	pool2, err := YAMLToProbePool(yamlByte)
	assert.NoError(t, err)
	assert.True(t, ProbePoolEquals(pool, pool2, t))
}

func ProbePoolEquals(this, that ProbePool, t *testing.T) bool {
	assert.Equal(t, len(this), len(that))
	for idx, thisPrb := range this {
		thatPrb := that[idx]
		assert.Equal(t, thisPrb.Name, thatPrb.Name)
		assert.Equal(t, thisPrb.ID, thatPrb.ID)
		assert.Equal(t, thisPrb.BaseURL, thatPrb.BaseURL)
		assert.Equal(t, thisPrb.Cron, thatPrb.Cron)
		assert.Equal(t, thisPrb.UpThreshold, thatPrb.UpThreshold)
		assert.Equal(t, thisPrb.DownThreshold, thatPrb.DownThreshold)

		if thisPrb.Requests != nil && thatPrb.Requests != nil {
			assert.Equal(t, len(thisPrb.Requests), len(thatPrb.Requests))
			for ridx, thisreq := range thisPrb.Requests {
				thatreq := thatPrb.Requests[ridx]
				assert.Equal(t, thisreq.Name, thatreq.Name)
				assert.Equal(t, thisreq.PathExpr, thatreq.PathExpr)
				assert.Equal(t, thisreq.MethodExpr, thatreq.MethodExpr)
				if thisreq.HeadersExpr != nil && thatreq.HeadersExpr != nil {
					if len(thisreq.HeadersExpr) != len(thatreq.HeadersExpr) {
						return false
					}
					// todo compare the map[string][]string this and that
				} else if thisreq.HeadersExpr == nil && thatreq.HeadersExpr == nil {
					// ok
				} else {
					return false
				}
				assert.Equal(t, thisreq.BodyExpr, thatreq.BodyExpr)
				assert.Equal(t, thisreq.CertificateCheckExpr, thatreq.CertificateCheckExpr)
				assert.Equal(t, thisreq.StartRequestIfExpr, thatreq.StartRequestIfExpr)
				assert.Equal(t, thisreq.SuccessIfExpr, thatreq.SuccessIfExpr)
				assert.Equal(t, thisreq.FailIfExpr, thatreq.FailIfExpr)
			}
		} else if thisPrb.Requests == nil && thatPrb.Requests == nil {
			// ok
		} else {
			return false
		}
		if thisPrb.SMTPNotification != nil && thatPrb.SMTPNotification != nil {
			assert.Equal(t, thisPrb.SMTPNotification.SMTPHost, thatPrb.SMTPNotification.SMTPHost)
			assert.Equal(t, thisPrb.SMTPNotification.SMTPPort, thatPrb.SMTPNotification.SMTPPort)
			assert.Equal(t, thisPrb.SMTPNotification.Password, thatPrb.SMTPNotification.Password)
			if thisPrb.SMTPNotification.From != nil {
				if !thisPrb.SMTPNotification.From.Equals(thatPrb.SMTPNotification.From) {
					return false
				}
			}

			if thisPrb.SMTPNotification.To != nil && thatPrb.SMTPNotification.To != nil {
				if len(thisPrb.SMTPNotification.To) != len(thatPrb.SMTPNotification.To) {
					return false
				}
				// todo compare element of mailboxes
			} else if thisPrb.SMTPNotification.To == nil && thatPrb.SMTPNotification.To == nil {
				//OK
			} else {
				return false
			}
			if thisPrb.SMTPNotification.Cc != nil && thatPrb.SMTPNotification.Cc != nil {
				if len(thisPrb.SMTPNotification.Cc) != len(thatPrb.SMTPNotification.Cc) {
					return false
				}
				// todo compare element of mailboxes
			} else if thisPrb.SMTPNotification.Cc == nil && thatPrb.SMTPNotification.Cc == nil {
				//OK
			} else {
				return false
			}
			if thisPrb.SMTPNotification.Bcc != nil && thatPrb.SMTPNotification.Bcc != nil {
				if len(thisPrb.SMTPNotification.Bcc) != len(thatPrb.SMTPNotification.Bcc) {
					return false
				}
				// todo compare element of mailboxes
			} else if thisPrb.SMTPNotification.Bcc == nil && thatPrb.SMTPNotification.Bcc == nil {
				//OK
			} else {
				return false
			}
		} else if thisPrb.SMTPNotification == nil && thatPrb.SMTPNotification == nil {
			// ok
		} else {
			return false
		}
		if thisPrb.CallbackNotification != nil && thatPrb.CallbackNotification != nil {
			assert.Equal(t, thisPrb.CallbackNotification.UpCall, thatPrb.CallbackNotification.UpCall)
			assert.Equal(t, thisPrb.CallbackNotification.DownCall, thatPrb.CallbackNotification.DownCall)
		} else if thisPrb.CallbackNotification == nil && thatPrb.CallbackNotification == nil {
			// ok
		} else {
			return false
		}
	}
	return true
}
