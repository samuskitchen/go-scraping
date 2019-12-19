package http

import (
	"../../driver"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
	"regexp"
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

	loc, _ := time.LoadLocation("America/Bogota")
	var dataServer modelServe.DataServe
	var servers []modelServe.Serve
	var changeServer bool

	address = ValidateURL(address)
	payload, err := rp.repo.GetDomainByAddress(r.Context(), address)

	if (modelDomain.Domain{}) == payload {
		dm := modelDomain.Domain{}
		dm.Address = address
		dm.LastConsultation = time.Now().In(loc)

		idDomain, err := rp.repo.CreateDomain(r.Context(), dm)

		if err != nil {
			//log.Fatal(err)
			respondWithError(w, http.StatusNoContent, err.Error())
		}

		SaveDetailDomain(data, idDomain, rp, w, r)
	} else {

		//detailsDomain, err := rp.repo.GetDetailsByDomain(r.Context(), payload.ID, len(data.Endpoints))

		if err != nil {
			//log.Fatal(err)
			respondWithError(w, http.StatusNoContent, err.Error())
		}

		loc, _ := time.LoadLocation("America/Bogota")

		now := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), 0, loc)
		//now := time.Now().Round(.000000).In(loc)
		fmt.Println("\nToday : ", loc, " Time : ", now)

		pastDate := payload.LastConsultation
		fmt.Println("Time : ", pastDate) //
		fmt.Printf("###############################################################\n")
		diff := now.Sub(pastDate)

		hrs := int(diff.Hours()/60.0)
		fmt.Printf("Diffrence in Hours : %d Hours\n", hrs)


		/*hours := time.Now().In(loc).Sub(payload.LastConsultation).Hours()
		fmt.Println("Hours", hours)

		if hours >= 1 {
			if len(data.Endpoints) == len(detailsDomain){
				for i := 0; i < len(data.Endpoints); i++ {
					if data.Endpoints[i].Grade != detailsDomain[i].Grade ||
						data.Endpoints[i].ServerName != detailsDomain[i].ServerName ||
						data.Endpoints[i].IpAddress != detailsDomain[i].IpAddress{
						changeServer = true
					}
				}
			}else if len(data.Endpoints) > len(detailsDomain){
				changeServer = true
			}
		}



		if changeServer {
			SaveDetailDomain(data, payload.ID, rp, w, r)
		}*/
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

func ValidateURL(address string) string {
	space := regexp.MustCompile(`^(http:\/\/www\.|https:\/\/www\.|http:\/\/|https:\/\/)`)
	address = space.ReplaceAllString(address, "$1\u200B")
	//address = strings.TrimSpace(address)
	//address = strings.TrimLeft(address, "")
	//address = strings.TrimRight(address, "")
	return address
}

func SaveDetailDomain(data modelSsl.SSL, idDomain int64, rp *Domain, w http.ResponseWriter, r *http.Request) {
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