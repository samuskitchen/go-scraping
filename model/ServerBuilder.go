package model

import (
	"../handler/command"
	modelSsl "../model/ssllabs"
)

func BuildServer(data modelSsl.SSL, servers []Serve,
	changeServer bool, currentGrade string, previousGrade string) DataServe {

	var dataServer DataServe

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
	dataServer.Logo = ""
	dataServer.Title = ""
	dataServer.IsDown = false

	return dataServer
}
