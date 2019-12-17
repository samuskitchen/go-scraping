package domain

import "time"

type DetailDomain struct {
	ID int64
	IDDomain int64
	IpAddress string
	ServerName string
	Grade string
	Date time.Time
}
