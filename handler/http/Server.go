package http

import (
	"../../driver"
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"

	//model_domain "../../model/domain"
	model_ssl "../../model/ssllabs"
	//model_serve "../../model"
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


func (p *Domain) Create(ssl *model_ssl.SSL) {

}


func (p *Domain) GetByAddress(w http.ResponseWriter, r *http.Request) {
	address := chi.URLParam(r, "address")
	data, err := GetDataSSl(address)

	//payload, err := p.repo.GetByDomain(r.Context(), address)

	if err != nil {
		respondWithError(w, http.StatusNoContent, "Address not found")
	}

	respondWithJSON(w, http.StatusOK, data)
}

func (p *Domain) GetAllAddress(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "address"))
	payload, err := p.repo.GetByDomain(r.Context(), string(id))

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