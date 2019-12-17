package ssllabs

type Endpoint struct {
	IpAddress         string `json:"ipAddress"`
	ServerName        string `json:"serverName"`
	StatusMessage     string `json:"statusMessage"`
	Grade             string `json:"grade"`
	GradeTrustIgnored string `json:"gradeTrustIgnored"`
	HasWarnings       bool   `json:"hasWarnings"`
	IsExceptional     bool   `json:"isExceptional"`
	Progress          int64  `json:"progress"`
	Duration          int64  `json:"duration"`
	Delegation        int64  `json:"delegation"`
}
