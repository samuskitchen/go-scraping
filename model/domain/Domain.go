package domain

import "time"

type Domain struct {
	ID               int64
	Address          string
	LastConsultation time.Time
}
