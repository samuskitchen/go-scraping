package repository

import (
	"../model/domain"
)

type DomainRepo interface {
	CreateDomain(domain domain.Domain) (int64, error)
	CreateDetailDomain(detailDomain domain.DetailDomain) error
	GetAllDomain()([]domain.Domain, error)
	GetDomainByAddress(address string) (domain.Domain, error)
	GetDetailsByDomain(idDomain int64, countServer int64) ([]domain.DetailDomain, error)
}
