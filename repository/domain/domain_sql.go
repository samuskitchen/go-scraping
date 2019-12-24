package domain

import (
	"../../model/domain"
	repo "../../repository"
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"time"
)

func NewSQLDomainRepo(Conn *sql.DB) repo.DomainRepo {
	return &sqlDomainRepo{
		Conn: Conn,
	}
}

type sqlDomainRepo struct {
	Conn *sql.DB
}

func (sql *sqlDomainRepo) CreateDomain(ctx context.Context, domain domain.Domain) (int64, error) {
	var domainId int64
	query := "INSERT INTO domain(address, last_consultation) VALUES($1, $2) RETURNING id"

	stmt, err := sql.Conn.PrepareContext(ctx, query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(domain.Address, domain.LastConsultation).Scan(&domainId)
	if err != nil {
		return -1, err
	}

	return domainId, err

}

func (sql *sqlDomainRepo) CreateDetailDomain(ctx context.Context, detailDomain domain.DetailDomain) error {
	query := "INSERT INTO detail_domain(id_domain, ipaddress, servername, grade, date) VALUES($1, $2, $3, $4, $5)"

	stmt, err := sql.Conn.PrepareContext(ctx, query)

	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.ExecContext(ctx, detailDomain.IDDomain, detailDomain.IpAddress, detailDomain.ServerName, detailDomain.Grade, detailDomain.Date)
	defer stmt.Close()

	if err != nil {
		log.Fatal(err)
	}

	return err
}

func (sql *sqlDomainRepo) UpdateLastGetDomain(ctx context.Context, idDomain int64, date time.Time) error {
	query := "UPDATE domain SET last_consultation = $1 WHERE id = $2"

	stmt, err := sql.Conn.PrepareContext(ctx, query)

	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.ExecContext(ctx, date, idDomain)
	defer stmt.Close()

	if err != nil {
		log.Fatal(err)
	}

	return err
}

func (sql *sqlDomainRepo) GetAllDomain(ctx context.Context) ([]domain.Domain, error) {
	query := "SELECT dm.id, dm.address, dm.last_consultation FROM domain AS dm"

	rows, err := sql.Conn.QueryContext(ctx, query)

	if err != nil {
		log.Fatal(err)
	}

	return domain.BuildDomains(rows)
}

func (sql *sqlDomainRepo) GetDomainByAddress(ctx context.Context, address string) (domain.Domain, error) {
	query := "SELECT dm.id, dm.address, dm.last_consultation FROM domain AS dm WHERE dm.address = $1"

	rows, err := sql.Conn.QueryContext(ctx, query, address)

	if err != nil {
		log.Fatal(err)
	}

	return domain.BuildDomain(rows)
}

func (sql *sqlDomainRepo) GetDetailsByDomain(ctx context.Context, idDomain int64, countServer int) ([]domain.DetailDomain, error) {
	query := "SELECT dt.id, dt.id_domain, dt.ipaddress, dt.grade, dt.servername, dt.date FROM domain AS dm INNER JOIN detail_domain AS dt ON dt.id_domain = dm.id WHERE dm.id = $1 ORDER BY dt.id, dt.date DESC LIMIT $2"

	rows, err := sql.Conn.QueryContext(ctx, query, idDomain, countServer)

	if err != nil {
		log.Fatal(err)
	}

	return domain.BuildDetailsDomain(rows)
}
