package domain

import (
	domain "../../model/domain"
	repo "../../repository"
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

func (s *sqlDomainRepo) CreateDomain(domain domain.Domain) (int64, error){
	var domainId int64
	query := "INSERT INTO domain(address) VALUES($1) RETURNING id"

	err := s.Conn.QueryRow(query, domain.Address).Scan(&domainId)

	if err != nil {
		log.Fatal(err)
	}

	return domainId, err

}

func (s *sqlDomainRepo) CreateDetailDomain(detailDomain domain.DetailDomain) error {
	query := "INSERT INTO detail_domain(id_domain, ipaddress, servername, grade, date) VALUES($1, $2, $3, $4, $5)"

	stmt, err := s.Conn.Prepare(query)

	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(detailDomain.IDDomain, detailDomain.IpAddress, detailDomain.ServerName, detailDomain.Grade, detailDomain.Date)

	if err != nil {
		log.Fatal(err)
	}

	return err
}

func (s *sqlDomainRepo) GetAllDomain() ([]domain.Domain, error) {
	query := "SELECT dm.id, dm.address FROM domain AS dm"

	rows, err := s.Conn.Query(query)

	if err != nil {
		log.Fatal(err)
	}

	return buildDomains(rows)
}

func (s *sqlDomainRepo) GetDomainByAddress(address string) (domain.Domain, error) {
	query := "SELECT dm.id, dm.address FROM domain AS dm WHERE dm.address = $1"

	rows, err := s.Conn.Query(query, address)

	if err != nil {
		log.Fatal(err)
	}

	return buildDomain(rows)
}

func (s *sqlDomainRepo) GetDetailsByDomain(idDomain int64, countServer int64) ([]domain.DetailDomain, error) {
	query := "SELECT dt.id, dt.id_domain, dt.ipaddress, dt.grade, dt.servername, dt.date FROM domain AS dm INNER JOIN detail_domain AS dt ON dt.id_domain = dm.id WHERE dm.address = $1 ORDER BY dt.date DESC LIMIT $2"

	rows, err := s.Conn.Query(query, idDomain, countServer)

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

		if err := rows.Scan(&b.ID, &b.Address); err != nil {
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

		if err := rows.Scan(&b.ID, &b.Address); err != nil {
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

		if err := rows.Scan(&b.ID, &b.IDDomain, &b.IpAddress, &b.ServerName, &b.Grade, &b.Date); err != nil {
			return nil, err
		}
		results = append(results, b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
