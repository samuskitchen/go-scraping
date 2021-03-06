package repository

import (
	"context"
	"go-scraping/model/domain"
	"time"
)

type DomainRepo interface {
	CreateDomain(ctx context.Context, domain domain.Domain) (int64, error)
	CreateDetailDomain(ctx context.Context, detailDomain domain.DetailDomain) error
	UpdateLastGetDomain(ctx context.Context, idDomain int64, date time.Time) error
	GetAllDomain(ctx context.Context)([]domain.Domain, error)
	GetDomainByAddress(ctx context.Context, address string) (domain.Domain, error)
	GetDetailsByDomain(ctx context.Context, idDomain int64, countServer int) ([]domain.DetailDomain, error)
}
