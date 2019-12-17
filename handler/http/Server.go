package http

import (
	"../../driver"
	"encoding/json"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"time"

	model_ssl "../../model/ssllabs"
	//model_serve "../../model"
	model_domain "../../model/domain"
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


func (rp *Domain) Create(ssl *model_ssl.SSL) {

}


func (rp *Domain) GetByAddress(w http.ResponseWriter, r *http.Request) {
	address := chi.URLParam(r, "address")
	data, err := GetDataSSl(address)

	payload, err := rp.repo.GetDomainByAddress(address)

	if (model_domain.Domain{}) == payload {
		dm := model_domain.Domain{}
		dm.Address = data.Host

		idDomain, err := rp.repo.CreateDomain(dm)

		if err != nil {
			log.Fatal(err)
			respondWithError(w, http.StatusNoContent, err.Error())
		}

		for _, element := range data.Endpoints {
			dt := model_domain.DetailDomain{}

			dt.IDDomain = idDomain
			dt.IpAddress = element.IpAddress
			dt.Grade = element.Grade
			dt.ServerName = element.ServerName
			dt.Date = time.Now()

			err = rp.repo.CreateDetailDomain(dt)

			if err != nil{
				log.Fatal(err)
				respondWithError(w, http.StatusNoContent, err.Error())
			}
		}
	} else {

	}

	if err != nil {
		respondWithError(w, http.StatusNoContent, "Address not found")
	}

	respondWithJSON(w, http.StatusOK, data)
}

func (rp *Domain) GetAllAddress(w http.ResponseWriter, r *http.Request) {
	payload, err := rp.repo.GetAllDomain()

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