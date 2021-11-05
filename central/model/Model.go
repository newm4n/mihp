package model

import (
	"time"
)

type ProbeNodeRepository interface {
	RegisterNewNode(node *ProbeNode) (bool, error)
	UpdateNode(node *ProbeNode) (bool, error)
	UnRegNode(node *ProbeNode) (bool, error)
	GetNodeByID(probeID string) (*ProbeNode, error)
	GetNodeByCountry(country string) ([]*ProbeNode, error)
	GetNodeByNetworkProvider(networkProvider string) ([]*ProbeNode, error)
	GetNodeByCity(city string) ([]*ProbeNode, error)
}

type ProbeNode struct {
	ProbeID         string    `json:"probe_id"`
	Country         string    `json:"country"`
	NetworkProvider string    `json:"network_provider"`
	City            string    `json:"city"`
	ProbeCount      int       `json:"probe_count"`
	StatusUp        bool      `json:"status_up"`
	LastContact     time.Time `json:"last_contact"`
	ProbeToken      string    `json:"probe_token"`
}

type UserRepository interface {
	GetUserRole(email, password string) (string, error)
	SaveNewUser(name, email, password, role string) (bool, error)
	UpdateUser(user *User) (bool, error)
	DeleteUser(user *User) (bool, error)
	GetUserByEmail(email *User) (*User, error)
	GetUsersByNamePrefix(prefix string) ([]*User, error)
	GetUsersByRole(role string) ([]*User, error)
	GetUserOrganizations(user *User) ([]*Organization, error)
}

type User struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	Role         string `json:"role"`
}

type OrganizationRepository interface {
	SaveNewOrganization(organization *Organization) (bool, error)
	UpdateOrganization(organization *Organization) (bool, error)
	DeleteOrganization(organization *Organization) (bool, error)
	GetOrganizationByName(name string) (*Organization, error)
	GetOrganizationUsers(organization *Organization) ([]*User, error)
	GetOrganisationProbeDatas(organization *Organization) ([]*ProbeData, error)
}

type Organization struct {
	Name               string                 `json:"name"`
	ProbesID           []string               `json:"probes_id"`
	NotificationMethod string                 `json:"notification_method"`
	SMTPConfig         *SMTPConfiguration     `json:"smtp_config"`
	CallbackConfig     *CallbackConfiguration `json:"callback_config"`
}

type CallbackConfiguration struct {
	UpNotificationURL   string
	DownNotificationURL string
}

type SMTPConfiguration struct {
	Host  string   `json:"host"`
	Port  string   `json:"port"`
	From  string   `json:"from"`
	Email string   `json:"email"`
	CC    []string `json:"cc"`
	BCC   []string `json:"bcc"`
}

type ProbeDataRepository interface {
	SaveNewProbeData(probe *ProbeData) (bool, error)
	UpdateProbeData(probe *ProbeData) (bool, error)
	DeleteProbeData(probe *ProbeData) (bool, error)
	GetProbeDataByID(id string) (*ProbeData, error)
	GetProbeRequestDatas(probe *ProbeData) ([]*ProbeRequestData, error)
}

type ProbeData struct {
	ID                   string `json:"id"`
	Name                 string `json:"name"`
	Scheme               string `json:"scheme"`
	UserHttp             string `json:"user"`
	PasswordHttp         string `json:"http_password"`
	HostHttp             string `json:"host_http"`
	PortHttp             int    `json:"port_http"`
	CertificateCheckExpr string `json:"certificate_check_expr"`

	Cron          string `json:"cron"`
	UpThreshold   int    `json:"up_threshold"`
	DownThreshold int    `json:"down_threshold"`

	ProbeCountry         string `json:"probe_country"`
	ProbeNetworkProvider string `json:"probe_network_provider"`
	ProbeCity            string `json:"probe_city"`

	AssignedToProbeID string `json:"assigned_to_probe_id"`
}

type ProbeRequestDataRepository interface {
	SaveNewProbeRequestData(probe *ProbeRequestData) (bool, error)
	UpdateProbeRequestData(probe *ProbeRequestData) (bool, error)
	DeleteProbeRequestData(probe *ProbeRequestData) (bool, error)
	GetProbeRequestRequestDataByID(id string) (*ProbeRequestData, error)
}

type ProbeRequestData struct {
	Name               string              `json:"name"`
	ProbeDataID        string              `json:"probe_data_id"`
	Sequence           int                 `json:"sequence"`
	Description        string              `json:"description"`
	PathExpr           string              `json:"path_expr"`
	MethodExpr         string              `json:"method_expr"`
	HeadersExpr        map[string][]string `json:"headers_expr"`
	BodyExpr           string              `json:"body_expr"`
	StartRequestIfExpr string              `json:"start_request_if_expr"`
	SuccessIfExpr      string              `json:"success_if_expr"`
	FailIfExpr         string              `json:"fail_if_expr"`
}
