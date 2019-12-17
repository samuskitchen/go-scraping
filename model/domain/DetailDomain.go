package domain

import "time"

type DetailDomain struct {
	ID int
	IDDomain int
	IpAddress string
	ServerName string
	Grade string
	Date time.Time
}
