package domain

import (
	domain "../../model/domain"
	repo "../../repository"
	"context"
	"database/sql"
	_ "github.com/lib/pq"
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

func (sql *sqlDomainRepo) CreateDomain(ctx context.Context, domain domain.Domain) (int64, error){
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

func (sql *sqlDomainRepo) GetAllDomain(ctx context.Context) ([]domain.Domain, error) {
	query := "SELECT dm.id, dm.address, dm.last_consultation FROM domain AS dm"

	rows, err := sql.Conn.QueryContext(ctx, query)

	if err != nil {
		log.Fatal(err)
	}

	return buildDomains(rows)
}

func (sql *sqlDomainRepo) GetDomainByAddress(ctx context.Context, address string) (domain.Domain, error) {
	query := "SELECT dm.id, dm.address, dm.last_consultation FROM domain AS dm WHERE dm.address = $1"

	rows, err := sql.Conn.QueryContext(ctx, query, address)

	if err != nil {
		log.Fatal(err)
	}

	return buildDomain(rows)
}

func (sql *sqlDomainRepo) GetDetailsByDomain(ctx context.Context, idDomain int64, countServer int) ([]domain.DetailDomain, error) {
	query := "SELECT dt.id, dt.id_domain, dt.ipaddress, dt.grade, dt.servername, dt.date FROM domain AS dm INNER JOIN detail_domain AS dt ON dt.id_domain = dm.id WHERE dm.id = $1 ORDER BY dt.id, dt.date DESC LIMIT $2"

	rows, err := sql.Conn.QueryContext(ctx, query, idDomain, countServer)

	if err != nil {
		log.Fatal(err)
	}

	return buildDetailsDomain(rows)
}

func buildDomain(rows *sql.Rows) (domain.Domain, error) {
	var result domain.Domain
	resultEmpty := domain.Domain{}

	for rows.Next() {
		b := domain.Domain{}

		if err := rows.Scan(&b.ID, &b.Address, &b.LastConsultation); err != nil {
			return resultEmpty, err
		}

		result = b
	}

	if err := rows.Err(); err != nil {
		return resultEmpty, err
	}

	return result, nil
}

func buildDomains(rows *sql.Rows) ([]domain.Domain, error) {
	var results []domain.Domain

	for rows.Next() {
		b := domain.Domain{}

		if err := rows.Scan(&b.ID, &b.Address, &b.LastConsultation); err != nil {
			return nil, err
		}
		results = append(results, b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}


func buildDetailsDomain(rows *sql.Rows) ([]domain.DetailDomain, error) {
	var results []domain.DetailDomain

	for rows.Next() {
		b := domain.DetailDomain{}

		if err := rows.Scan(&b.ID, &b.IDDomain, &b.IpAddress, &b.Grade, &b.ServerName, &b.Date); err != nil {
			return nil, err
		}
		results = append(results, b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
