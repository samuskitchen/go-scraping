package repository

import (
	"../model"
	"../model/domain"
	"context"
)

type DomainRepo interface {
	GetByDomain(ctx context.Context, address string)(*model.DataServe, error)
	CreateDomain(domain *domain.Domain)
	CreateDetailDomain(detailDomain *domain.DetailDomain)
}