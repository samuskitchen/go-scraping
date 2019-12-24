package http

import (
	"../../driver"

	"github.com/go-chi/chi"
	"net/http"
	"time"

	"../../handler/command"
	modelServe "../../model"
	modelDomain "../../model/domain"
	modelSsl "../../model/ssllabs"
	"../../repository"
	"../../repository/domain"
)

func NewServerHandler(db *driver.DB) *Domain {
	return &Domain{
		repo: domain.NewSQLDomainRepo(db.SQL),
	}
}

type Domain struct {
	repo repository.DomainRepo
}

func (rp *Domain) GetByAddress(w http.ResponseWriter, r *http.Request) {
	address := chi.URLParam(r, "address")
	address = command.ValidateURL(address)

	data, err := GetDataSSl(address)

	loc, _ := time.LoadLocation("America/Bogota")
	var dataServer modelServe.DataServe
	var servers []modelServe.Serve
	var detailsDomain []modelDomain.DetailDomain
	var changeServer bool

	payload, err := rp.repo.GetDomainByAddress(r.Context(), address)

	if (modelDomain.Domain{}) == payload {
		dm := modelDomain.Domain{}
		dm.Address = address
		dm.LastConsultation = time.Now().In(loc)

		idDomain, err := rp.repo.CreateDomain(r.Context(), dm)

		if err != nil {
			command.RespondWithError(w, http.StatusNoContent, err.Error())
		}

		saveDetailDomain(data, idDomain, rp, w, r)
	} else {

		detailsDomain, err := rp.repo.GetDetailsByDomain(r.Context(), payload.ID, len(data.Endpoints))

		if err != nil {
			command.RespondWithError(w, http.StatusNoContent, err.Error())
		}

		changeServer = command.ValidateChangeServer(loc, payload, data, detailsDomain, changeServer)

		if changeServer {
			err = rp.repo.UpdateLastGetDomain(r.Context(), payload.ID, time.Now())
			saveDetailDomain(data, payload.ID, rp, w, r)
		}
	}

	if err != nil {
		command.RespondWithError(w, http.StatusNoContent, "Address not found")
	}

	buildServer(data, detailsDomain, servers, dataServer, changeServer, w)

}

func buildServer(data modelSsl.SSL, detailsDomain []modelDomain.DetailDomain, servers []modelServe.Serve,
	dataServer modelServe.DataServe, changeServer bool, w http.ResponseWriter) {

	if "IN_PROGRESS" != data.Status && "DNS" != data.Status {
		currentGrade := command.GetLowestGradeCurrent(data.Endpoints)
		var previousGrade string

		if detailsDomain == nil {
			previousGrade = currentGrade
		} else {
			previousGrade = command.GetLowestGradePrevious(detailsDomain)
		}

		//TODO Build return data
		dataServer = modelServe.BuildServer(data, servers, changeServer, currentGrade, previousGrade)

		command.RespondWithJSON(w, http.StatusOK, dataServer)
	} else {
		command.RespondWithJSON(w, http.StatusOK, "Try later the server data is not yet available, Thank you!")
	}
}

func (rp *Domain) GetAllAddress(w http.ResponseWriter, r *http.Request) {
	payload, err := rp.repo.GetAllDomain(r.Context())

	if err != nil {
		command.RespondWithError(w, http.StatusNoContent, "Address not found")
	}

	command.RespondWithJSON(w, http.StatusOK, payload)
}

// Method that saves the main details of the domain or server
func saveDetailDomain(data modelSsl.SSL, idDomain int64, rp *Domain, w http.ResponseWriter, r *http.Request) {
	loc, _ := time.LoadLocation("America/Bogota")

	for _, element := range data.Endpoints {
		dt := modelDomain.DetailDomain{}

		dt.IDDomain = idDomain
		dt.IpAddress = element.IpAddress
		dt.Grade = element.Grade
		dt.ServerName = element.ServerName
		dt.Date = time.Now().In(loc)

		err := rp.repo.CreateDetailDomain(r.Context(), dt)

		if err != nil {
			command.RespondWithError(w, http.StatusNoContent, err.Error())
		}
	}
}
