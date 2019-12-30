package model

import (
	"../handler/command"
	modelDomain "../model/domain"
	modelSsl "../model/ssllabs"
)

func BuildServer(data modelSsl.SSL, detailsDomain []modelDomain.DetailDomain, changeServer bool, pageTitle string, pageLogo string) DataServe {

	if "IN_PROGRESS" != data.Status && "DNS" != data.Status {
		currentGrade := command.GetLowestGradeCurrent(data.Endpoints)
		var previousGrade string

		if detailsDomain == nil {
			previousGrade = currentGrade
		} else {
			previousGrade = command.GetLowestGradePrevious(detailsDomain)
		}

		dataServer := buildData(data, currentGrade, previousGrade, changeServer, pageTitle, pageLogo)

		return dataServer
	} else {
		return DataServe{}
	}
}

func buildData(data modelSsl.SSL, currentGrade string, previousGrade string, changeServer bool, pageTitle string, pageLogo string) DataServe {
	servers := make([]Serve, 0)
	dataServer := DataServe{}

	for _, dataElement := range data.Endpoints {
		serve := Serve{}
		result := command.RunWhoIs(dataElement.IpAddress)

		serve.Address = dataElement.IpAddress
		serve.SslGrade = dataElement.Grade
		serve.Country = result["Country"][0]
		serve.Owner = result["OrgName"][0]

		servers = append(servers, serve)
	}

	dataServer.Serves = servers
	dataServer.ServersChanged = changeServer
	dataServer.SslGrade = currentGrade
	dataServer.PreviousSslGrade = previousGrade
	dataServer.Logo = pageLogo
	dataServer.Title = pageTitle
	dataServer.IsDown = false

	return dataServer
}
