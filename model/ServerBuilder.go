package model

import (
	"go-scraping/handler/command"
	modelDomain "go-scraping/model/domain"
	modelSsl "go-scraping/model/ssllabs"
)

func BuilderAddress(payload []modelDomain.Domain) Items {
	payloadItems := Items{}
	items := make([]ItemServe, 0)

	for _, element := range payload {
		items = append(items, ItemServe{Address: element.Address})
	}
	payloadItems.Items = items
	return payloadItems
}

func BuildServer(data modelSsl.SSL, detailsDomain []modelDomain.DetailDomain, changeServer bool, pageTitle string, pageLogo string) DataServe {

	currentGrade := command.GetLowestGradeCurrent(data.Endpoints)
	var previousGrade string

	if detailsDomain == nil {
		previousGrade = currentGrade
	} else {
		previousGrade = command.GetLowestGradePrevious(detailsDomain)
	}

	dataServer := buildData(data, currentGrade, previousGrade, changeServer, pageTitle, pageLogo)

	return dataServer
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
