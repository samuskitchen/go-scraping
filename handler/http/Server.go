package http

import (
	"../../driver"
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
	"time"

	modelServe "../../model"
	modelDomain "../../model/domain"
	modelSsl "../../model/ssllabs"
	repository "../../repository"
	domain "../../repository/domain"
)

func NewServerHandler(db *driver.DB) *Domain{
	return &Domain{
		repo: domain.NewSQLDomainRepo(db.SQL),
	}
}

type Domain struct {
	repo repository.DomainRepo
}


func (rp *Domain) Create(ssl *modelSsl.SSL) {

}


func (rp *Domain) GetByAddress(w http.ResponseWriter, r *http.Request) {
	address := chi.URLParam(r, "address")
	data, err := GetDataSSl(address)
	var dataServer modelServe.DataServe
	var servers []modelServe.Serve

	/*it, err := rp.repo.GetAllDomain(r.Context())
	fmt.Println(it)*/
	payload, err := rp.repo.GetDomainByAddress(r.Context(), address)
	var changeServer bool

	if (modelDomain.Domain{}) == payload {
		dm := modelDomain.Domain{}
		dm.Address = data.Host
		dm.LastConsultation = time.Now()

		idDomain, err := rp.repo.CreateDomain(r.Context(), dm)

		if err != nil {
			//log.Fatal(err)
			respondWithError(w, http.StatusNoContent, err.Error())
		}

		SaveDetailDomain(data, idDomain, rp, w, r)
	} else {

		detailsDomain, err := rp.repo.GetDetailsByDomain(r.Context(), payload.ID,  len(data.Endpoints))

		if err != nil {
			//log.Fatal(err)
			respondWithError(w, http.StatusNoContent, err.Error())
		}

		for _, element := range detailsDomain{
			for _, dataElement := range data.Endpoints {
				if dataElement.Grade != element.Grade || dataElement.ServerName != element.ServerName || dataElement.IpAddress != element.IpAddress{
					changeServer = true
				}
			}
		}

		if changeServer {
			SaveDetailDomain(data, payload.ID, rp, w, r)
		}
	}

	//TODO Build return data
	for _, dataElement := range data.Endpoints {
		serve := modelServe.Serve{}

		serve.Address = dataElement.IpAddress
		serve.SslGrade = dataElement.Grade
		serve.Country = ""
		serve.Owner = ""

		servers = append(servers, serve)
	}

	dataServer.Serves = servers
	dataServer.ServersChanged = changeServer
	dataServer.SslGrade = "A+"
	dataServer.PreviousSslGrade = "B"
	dataServer.Logo = ""
	dataServer.Title = ""
	dataServer.IsDown = false

	if err != nil {
		respondWithError(w, http.StatusNoContent, "Address not found")
	}

	respondWithJSON(w, http.StatusOK, dataServer)
}

func SaveDetailDomain(data modelSsl.SSL, idDomain int64, rp *Domain, w http.ResponseWriter, r *http.Request) {
	for _, element := range data.Endpoints {
		dt := modelDomain.DetailDomain{}

		dt.IDDomain = idDomain
		dt.IpAddress = element.IpAddress
		dt.Grade = element.Grade
		dt.ServerName = element.ServerName
		dt.Date = time.Now()

		err := rp.repo.CreateDetailDomain(r.Context(), dt)

		if err != nil {
			//log.Fatal(err)
			respondWithError(w, http.StatusNoContent, err.Error())
		}
	}
}

func (rp *Domain) GetAllAddress(w http.ResponseWriter, r *http.Request) {
	payload, err := rp.repo.GetAllDomain(r.Context())

	if err != nil {
		respondWithError(w, http.StatusNoContent, "Address not found")
	}

	respondWithJSON(w, http.StatusOK, payload)
}

// respondWithJSON write json response format
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// respondWithError return error message
func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"message": msg})
}