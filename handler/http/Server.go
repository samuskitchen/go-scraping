package http

import (
	"../../driver"
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
	"regexp"
	"sort"
	"time"

	modelServe "../../model"
	modelDomain "../../model/domain"
	modelSsl "../../model/ssllabs"
	repository "../../repository"
	domain "../../repository/domain"
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
	data, err := GetDataSSl(address)

	loc, _ := time.LoadLocation("America/Bogota")
	var dataServer modelServe.DataServe
	var servers []modelServe.Serve
	var detailsDomain []modelDomain.DetailDomain
	var changeServer bool

	address = validateURL(address)
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

		saveDetailDomain(data, idDomain, rp, w, r)
	} else {

		detailsDomain, err := rp.repo.GetDetailsByDomain(r.Context(), payload.ID, len(data.Endpoints))

		if err != nil {
			//log.Fatal(err)
			respondWithError(w, http.StatusNoContent, err.Error())
		}

		changeServer = validateChangeServer(loc, payload, data, detailsDomain, changeServer)

		if changeServer {
			saveDetailDomain(data, payload.ID, rp, w, r)
		}
	}

	currentGrade := getLowestGradeCurrent(data.Endpoints)
	var previousGrade string

	if detailsDomain == nil {
		previousGrade = currentGrade
	} else {
		previousGrade = getLowestGradePrevious(detailsDomain)
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
	dataServer.SslGrade = currentGrade
	dataServer.PreviousSslGrade = previousGrade
	dataServer.Logo = ""
	dataServer.Title = ""
	dataServer.IsDown = false

	if err != nil {
		respondWithError(w, http.StatusNoContent, "Address not found")
	}

	respondWithJSON(w, http.StatusOK, dataServer)
}

func (rp *Domain) GetAllAddress(w http.ResponseWriter, r *http.Request) {
	payload, err := rp.repo.GetAllDomain(r.Context())

	if err != nil {
		respondWithError(w, http.StatusNoContent, "Address not found")
	}

	respondWithJSON(w, http.StatusOK, payload)
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
			//log.Fatal(err)
			respondWithError(w, http.StatusNoContent, err.Error())
		}
	}
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

// We get the lowest grade from current servers
func getLowestGradeCurrent(data []modelSsl.Endpoint) string {
	var gradeAscii []int
	var grade string

	for _, dataElement := range data {
		if dataElement.Grade != "A+" {
			gradeAscii = append(gradeAscii, int(dataElement.Grade[0]))
		} else {
			grade = "A+"
		}
	}

	sort.Slice(gradeAscii, func(i, j int) bool {
		return gradeAscii[i] > gradeAscii[j]
	})

	grade = string(gradeAscii[0])
	return grade
}

// We get the lowest grade from previous servers
func getLowestGradePrevious(detail []modelDomain.DetailDomain) string {
	var gradeAscii []int
	var grade string

	for _, dataElement := range detail {
		if dataElement.Grade != "A+" {
			gradeAscii = append(gradeAscii, int(dataElement.Grade[0]))
		} else {
			grade = "A+"
		}
	}

	sort.Slice(gradeAscii, func(i, j int) bool {
		return gradeAscii[i] > gradeAscii[j]
	})

	grade = string(gradeAscii[0])
	return grade
}

// Validate if there is a change in the main data of the servers
func validateChangeServer(loc *time.Location, payload modelDomain.Domain, data modelSsl.SSL, detailsDomain []modelDomain.DetailDomain, changeServer bool) bool {
	hours := diffHours(loc, payload)
	if hours >= 1 {
		if len(data.Endpoints) == len(detailsDomain) {
			for i := 0; i < len(data.Endpoints); i++ {
				if data.Endpoints[i].Grade != detailsDomain[i].Grade ||
					data.Endpoints[i].ServerName != detailsDomain[i].ServerName ||
					data.Endpoints[i].IpAddress != detailsDomain[i].IpAddress {
					changeServer = true
				}
			}
		} else if len(data.Endpoints) > len(detailsDomain) {
			changeServer = true
		}
	}
	return changeServer
}

// Total hours difference
func diffHours(loc *time.Location, payload modelDomain.Domain) float64 {
	t1 := time.Date(time.Now().Year(),
		time.Now().Month(),
		time.Now().Day(),
		time.Now().Hour(),
		0, 0, 0, loc)

	t2 := time.Date(payload.LastConsultation.Year(),
		payload.LastConsultation.Month(),
		payload.LastConsultation.Day(),
		payload.LastConsultation.Hour(),
		0, 0, 0, loc)

	return t1.Sub(t2).Hours()
}

// Validate that the URL is cleaned
func validateURL(address string) string {
	space := regexp.MustCompile(`^(http:\/\/www\.|https:\/\/www\.|http:\/\/|https:\/\/)`)
	address = space.ReplaceAllString(address, "")
	return address
}