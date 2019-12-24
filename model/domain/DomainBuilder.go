package domain

import "database/sql"

func BuildDomain(rows *sql.Rows) (Domain, error) {
	var result Domain
	resultEmpty := Domain{}

	for rows.Next() {
		b := Domain{}

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

func BuildDomains(rows *sql.Rows) ([]Domain, error) {
	var results []Domain

	for rows.Next() {
		b := Domain{}

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

func BuildDetailsDomain(rows *sql.Rows) ([]DetailDomain, error) {
	var results []DetailDomain

	for rows.Next() {
		b := DetailDomain{}

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
