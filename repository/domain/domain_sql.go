package domain

import (
	serve "../../model"
	domain "../../model/domain"
	repo "../../repository"
	"context"
	"database/sql"
	"log"
)

func NewSQLDomainRepo(Conn *sql.DB) repo.DomainRepo {
	return &sqlDomainRepo{
		Conn: Conn,
	}
}

type sqlDomainRepo struct {
	Conn *sql.DB
}

func (s *sqlDomainRepo) GetByDomain(ctx context.Context, address string) (*serve.DataServe, error) {
	log.Println("implement me")

	return nil, nil
}

func (s *sqlDomainRepo) CreateDomain(domain *domain.Domain) {
	panic("implement me")
}

func (s *sqlDomainRepo) CreateDetailDomain(detailDomain *domain.DetailDomain) {
	panic("implement me")
}
