package internal

import (
	"fmt"
	yaml "gopkg.in/yaml.v3"
	"net"
	"regexp"
	"strings"
)

const (
	version = "1.0.0"
)

type MIHPConfig struct {
	Version   string         `yaml:"version"`
	ProbePool ProbePool      `yaml:"probe_pool"`
	Central   *CentralConfig `yaml:"central"`
	Minion    *MinionConfig  `yaml:"minion"`
}

func YAMLToMIHPConfig(yamlBytes []byte) (probePool *MIHPConfig, err error) {
	config := &MIHPConfig{}
	err = yaml.Unmarshal(yamlBytes, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func MIHPConfigToYAML(config *MIHPConfig) (yamlBytes []byte, err error) {
	return yaml.Marshal(config)
}

type CentralConfig struct {
	ListenHost              string `yaml:"listen_host"`
	ListenPort              int    `yaml:"listen_port"`
	ReadHeaderTimeoutSecond int    `yaml:"read_header_timeout_second"`
	ReadTimeoutSecond       int    `yaml:"read_timeout"`
	WriteTimeoutSecond      int    `yaml:"write_timeout"`
	IdleTimeoutSecond       int    `yaml:"idle_timeout"`

	AdminUser     string `yaml:"admin_user"`
	AdminPassword string `yaml:"admin_password"`

	JWTIssuer              string `yaml:"jwt_issuer"`
	JWTKey                 string `yaml:"jwt_key"`
	JWTAccessKeyAgeMinute  int    `yaml:"jwt_access_key_age_minute"`
	JWTRefreshKeyAgeMinute int    `yaml:"jwt_refresh_key_age_minute"`

	MySQLConfig      *DBConfig `yaml:"my_sql_config"`
	PostgreSQLConfig *DBConfig `yaml:"postgre_sql_config"`
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

type MinionConfig struct {
	CentralBaseURL string
	ReportCron     string
	Name           string
	CountryISO     string
	Datacenter     string
	MinionUID      string
}

type ProbePool []*Probe

type Probe struct {
	Name                 string                      `json:"name" yaml:"name"`
	ID                   string                      `json:"id" yaml:"id"`
	Requests             []*ProbeRequest             `json:"requests" yaml:"requests"`
	BaseURL              string                      `json:"base_url" yaml:"base_url""`
	Cron                 string                      `json:"cron" yaml:"cron"`
	UpThreshold          int                         `json:"up_threshold" yaml:"up_threshold"`
	DownThreshold        int                         `json:"down_threshold" yaml:"down_threshold"`
	SMTPNotification     *SMTPNotificationTarget     `json:"smtp_notification" yaml:"SMTP_notification"`
	CallbackNotification *CallbackNotificationTarget `json:"callback_notification" yaml:"callback_notification"`
}

type SMTPNotificationTarget struct {
	SMTPHost string     `yaml:"smtp_host"`
	SMTPPort int        `yaml:"smtp_port"`
	From     *Mailbox   `yaml:"from"`
	Password string     `yaml:"password"`
	To       []*Mailbox `yaml:"to"`
	Cc       []*Mailbox `yaml:"cc"`
	Bcc      []*Mailbox `yaml:"bcc"`
}

type CallbackNotificationTarget struct {
	UpCall   string `yaml:"up_call"`
	DownCall string `yaml:"down_call"`
}

type Mailbox struct {
	Name  string `yaml:"name"`
	Email string `yaml:"email"`
}

func NewMailbox(e string) (*Mailbox, error) {
	var display, email string
	if strings.Contains(e, "<") && strings.Contains(e, ">") {
		if strings.Index(e, "<") > strings.Index(e, ">") {
			return nil, fmt.Errorf("%s is not a valid email address", e)
		}
		display = e[:strings.Index(e, "<")]
		email = e[strings.Index(e, "<")+1 : strings.Index(e, ">")]
	} else if !strings.Contains(e, "<") && !strings.Contains(e, ">") {
		email = strings.TrimSpace(e)
	} else {
		return nil, fmt.Errorf("%s is not a valid email address", e)
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return nil, fmt.Errorf("%s is not a valid email address", e)
	}
	if parts[0] == "" {
		return nil, fmt.Errorf("%s is not a valid email address", e)
	}
	if b, _ := regexp.MatchString(`\s`, parts[0]); b {
		return nil, fmt.Errorf("%s is not a valid email address", e)
	}
	mx, err := net.LookupMX(parts[1])
	if err != nil || len(mx) == 0 {
		return nil, fmt.Errorf("%s is not a valid email address", e)
	}
	return &Mailbox{
		Name:  display,
		Email: email,
	}, nil
}

func (this *Mailbox) Equals(that *Mailbox) bool {
	if that == nil {
		return false
	}
	if this.Name != that.Name || this.Email != that.Email {
		return false
	}
	return true
}

func (this *Mailbox) String() string {
	if this.Name == "" {
		return this.Email
	}
	return fmt.Sprintf("%s <%s>", this.Name, this.Email)
}

type ProbeRequest struct {
	Name                 string              `json:"name" yaml:"name"`
	PathExpr             string              `json:"path_expr" yaml:"path_expr"`
	MethodExpr           string              `json:"method_expr" yaml:"method_expr"`
	HeadersExpr          map[string][]string `json:"headers_exprs" yaml:"headers_expr"`
	BodyExpr             string              `json:"body_expr" yaml:"body_expr"`
	CertificateCheckExpr string              `json:"certificate_check_expr" yaml:"certificate_check_expr"`
	StartRequestIfExpr   string              `json:"start_request_if_expr" yaml:"start_request_if_expr"`
	SuccessIfExpr        string              `json:"success_if_expr" yaml:"success_if_expr"`
	FailIfExpr           string              `json:"fail_if_expr" yaml:"fail_if_expr"`
}

func YAMLToProbePool(yamlBytes []byte) (probePool ProbePool, err error) {
	pool := make(ProbePool, 0)
	err = yaml.Unmarshal(yamlBytes, &pool)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func ProbePoolToYAML(pool ProbePool) (yamlBytes []byte, err error) {
	return yaml.Marshal(pool)
}
