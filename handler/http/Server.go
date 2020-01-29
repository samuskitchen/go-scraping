package http

import (
	"net/http"
	"time"

	"../../driver"
	"github.com/go-chi/chi"

	"../../handler/command"
	modelServe "../../model"
	modelDomain "../../model/domain"
	modelSsl "../../model/ssllabs"
	"../../repository"
	"../../repository/domain"
)

// NewServerHandler ...
func NewServerHandler(db *driver.DB) *Domain {
	return &Domain{
		repo: domain.NewSQLDomainRepo(db.SQL),
	}
}

// Domain ...
type Domain struct {
	repo repository.DomainRepo
}

// GetByAddress We get all the information by address
func (rp *Domain) GetByAddress(w http.ResponseWriter, r *http.Request) {
	address := chi.URLParam(r, "address")
	address = command.ValidateURL(address)

	data, err := GetDataSSl(address)
	if err != nil {
		command.RespondWithError(w, http.StatusNoContent, err.Error())
		return
	}

	if "IN_PROGRESS" == data.Status || "DNS" == data.Status || "" == data.Status{
		command.RespondWithJSON(w, http.StatusPartialContent, "Try later the server data is not yet available, Thank you!")
		return
	}

	pageTitle, pageLogo, err := GetTitleAndLogo(address)
	if err != nil {
		command.RespondWithError(w, http.StatusPartialContent, "Address not found")
		return
	}

	loc, _ := time.LoadLocation("America/Bogota")
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
			return
		}

		saveDetailDomain(data, idDomain, rp, w, r)
	} else {

		detailsDomain, err := rp.repo.GetDetailsByDomain(r.Context(), payload.ID, len(data.Endpoints))
		if err != nil {
			command.RespondWithError(w, http.StatusNoContent, err.Error())
			return
		}

		changeServer = command.ValidateChangeServer(loc, payload, data, detailsDomain, changeServer)
		if changeServer {
			err = rp.repo.UpdateLastGetDomain(r.Context(), payload.ID, time.Now())
			saveDetailDomain(data, payload.ID, rp, w, r)
		}
	}

	if err != nil {
		command.RespondWithError(w, http.StatusNoContent, "Address not found")
		return
	}

	dataServer := modelServe.BuildServer(data, detailsDomain, changeServer, pageTitle, pageLogo)
	command.RespondWithJSON(w, http.StatusOK, dataServer)
}

// GetAllAddress We get all the addresses we have consulted
func (rp *Domain) GetAllAddress(w http.ResponseWriter, r *http.Request) {
	payload, err := rp.repo.GetAllDomain(r.Context())

	if err != nil {
		command.RespondWithError(w, http.StatusNoContent, "Address not found")
		return
	}

	payloadItems := modelServe.BuilderAddress(payload)
	command.RespondWithJSON(w, http.StatusOK, payloadItems)
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
