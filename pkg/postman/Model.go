package postman

type Collection struct {
	Info     *Information `json:"info"`
	Item     []*Items     `json:"item"`
	Event    []*Event     `json:"event"`
	Variable []*Variable  `json:"variable"`
	Auth     *Auth        `json:"auth,omitempty"`
}

type Information struct {
	Name        string   `json:"name"`
	PostmanID   string   `json:"_postman_id"`
	Description string   `json:"description"`
	Version     *Version `json:"version"`
	Schema      string   `json:"schema"`
}

type Version struct { // todo need manual parsing
	Major       int    `json:"major"`
	Minor       int    `json:"minor"`
	Patch       int    `json:"patch"`
	Identifier  string `json:"identifier"`
	FullVersion string `json:"full_version"`
}

type Items struct { // todo need manual parsing
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Variable    []*Variable `json:"variable"`
	Event       []*Event    `json:"event"`
}

type Item struct {
	Items
	ID       string      `json:"id"`
	Request  *Request    `json:"request"`
	Response []*Response `json:"response"`
}

type Folder struct {
	Items
	Item []*Items `json:"item"`
	Auth *Auth    `json:"auth,omitempty"`
}

type Event struct {
	ID       string `json:"id"`
	Listen   string `json:"listen"`
	Script   Script `json:"script"`
	Disabled bool   `json:"disabled"`
}

type Script struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Exec string `json:"exec"`
	Src  URL    `json:"src"`
	Name string `json:"name"`
}

type URL struct {
	Raw      string        `json:"raw"`
	Protocol string        `json:"protocol"`
	Host     string        `json:"host"`
	Path     string        `json:"path"`
	Port     string        `json:"port"`
	Query    *[]QueryParam `json:"query"`
	Hash     string        `json:"hash"`
	Variable *[]Variable   `json:"variable"`
}

type QueryParam struct {
	Key         string       `json:"key"`
	Value       string       `json:"value"`
	Disabled    bool         `json:"disabled"`
	Description *Description `json:"description"`
}

type Description struct { // todo need manual parsing
	Text    string `json:"text"`
	Content string `json:"content"`
	Type    string `json:"type"`
	Version string `json:"version"`
}

type Variable struct {
	ID          string       `json:"id"`
	Key         string       `json:"key"`
	Value       string       `json:"value"`
	Type        string       `json:"type"`
	Name        string       `json:"name"`
	Description *Description `json:"description"`
	System      string       `json:"system"`
	Disabled    bool         `json:"disabled"`
}

type Auth struct {
	Type     string
	Apikey   *AuthInfo `json:"apikey,omitempty"`
	Awsv4    *AuthInfo `json:"awsv4,omitempty"`
	Basic    *AuthInfo `json:"basic,omitempty"`
	Bearer   *AuthInfo `json:"bearer,omitempty"`
	Digest   *AuthInfo `json:"digest,omitempty"`
	Edgegrid *AuthInfo `json:"edgegrid,omitempty"`
	Hawk     *AuthInfo `json:"hawk,omitempty"`
	Ntlm     *AuthInfo `json:"ntlm,omitempty"`
	Oauth1   *AuthInfo `json:"oauth1,omitempty"`
	Oauth2   *AuthInfo `json:"oauth2,omitempty"`
}

type AuthInfo struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type Request struct { // todo requires manual parsing
	RawRequest  string
	URL         *URL         `json:"url"`
	Auth        *Auth        `json:"auth"`
	Proxy       *Proxy       `json:"proxy"`
	Certificate *Certificate `json:"certificate"`
	Method      string       `json:"method"`
	Description *Description `json:"description"`
	Header      []*Header    `json:"header"`
	Body        *Body        `json:"body"`
}

type Proxy struct {
	Match    string `json:"match"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Tunnel   bool   `json:"tunnel"`
	Disabled bool   `json:"disabled"`
}

type Certificate struct {
	Name       string   `json:"name"`
	Matches    []string `json:"matches"`
	Key        *KeySrc  `json:"key"`
	Cert       *KeySrc  `json:"cert"`
	Passphrase string   `json:"passphrase"`
}

type KeySrc struct {
	Src string `json:"src"`
}

type Header struct {
	Key         string       `json:"key"`
	Value       string       `json:"value"`
	Disabled    bool         `json:"disabled"`
	Description *Description `json:"description"`
}

type Body struct {
	Mode       string                 `json:"mode"`
	Raw        string                 `json:"raw"`
	GraphQL    map[string]interface{} `json:"graphql"`
	URLEncoded []*URLEncoded          `json:"url_encoded"`
	Formdata   []*FormData            `json:"formdata"`
	File       *File                  `json:"file"`
	Options    map[string]interface{} `json:"options"`
	Disabled   bool                   `json:"disabled"`
}

type URLEncoded struct {
	Key         string       `json:"key"`
	Value       string       `json:"value"`
	Disabled    bool         `json:"disabled"`
	Description *Description `json:"description"`
}

type FormData struct {
	Key         string       `json:"key"`
	Value       string       `json:"value"`
	Src         string       `json:"src"`
	Disabled    bool         `json:"disabled"`
	Type        string       `json:"type"`
	ContentType string       `json:"contentType"`
	Description *Description `json:"description"`
}

type File struct {
	Src     string `json:"src"`
	Content string `json:"content"`
}

type Response struct {
	ID              string
	OriginalRequest *Request               `json:"originalRequest"`
	ResponseTime    int                    `json:"responseTime"` // todo need manual parsing
	Timings         map[string]interface{} `json:"timings"`
	Header          []*Header              `json:"header"`
	Cookie          []*Cookie              `json:"cookie"`
	Body            string                 `json:"body"`
	Status          string                 `json:"status"`
	Code            int                    `json:"code"`
}

type Cookie struct {
	Domain     string   `json:"domain"`
	Expires    string   `json:"expires"`
	MaxAge     string   `json:"maxAge"`
	HostOnly   bool     `json:"hostOnly"`
	HttpOnly   bool     `json:"httpOnly"`
	Name       string   `json:"name"`
	Path       string   `json:"path"`
	Secure     bool     `json:"secure"`
	Session    bool     `json:"session"`
	Value      string   `json:"value"`
	Extensions []string `json:"extensions"`
}
