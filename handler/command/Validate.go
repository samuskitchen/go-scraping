package command

import (
	"encoding/json"
	"net/http"
	"regexp"
	"sort"
	"time"

	modelDomain "../../model/domain"
	modelSsl "../../model/ssllabs"
)

// respondWithJSON write json response format
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// respondWithError return error message
func RespondWithError(w http.ResponseWriter, code int, msg string) {
	RespondWithJSON(w, code, map[string]string{"message": msg})
}

// We get the lowest grade from current servers
func GetLowestGradeCurrent(data []modelSsl.Endpoint) string {
	var gradeAscii []int
	var grade string

	for _, dataElement := range data {
		if dataElement.Grade != "A+" {
			gradeAscii = append(gradeAscii, int(dataElement.Grade[0]))
		} else {
			grade = "A+"
		}
	}

	if len(gradeAscii) > 0 {
		sort.Slice(gradeAscii, func(i, j int) bool {
			return gradeAscii[i] > gradeAscii[j]
		})

		grade = string(gradeAscii[0])
	}

	return grade
}

// We get the lowest grade from previous servers
func GetLowestGradePrevious(detail []modelDomain.DetailDomain) string {
	var gradeAscii []int
	var grade string

	for _, dataElement := range detail {
		if dataElement.Grade != "A+" {
			gradeAscii = append(gradeAscii, int(dataElement.Grade[0]))
		} else {
			grade = "A+"
		}
	}

	if len(gradeAscii) > 0 {
		sort.Slice(gradeAscii, func(i, j int) bool {
			return gradeAscii[i] > gradeAscii[j]
		})

		grade = string(gradeAscii[0])
	}

	return grade
}

// Validate if there is a change in the main data of the servers
func ValidateChangeServer(loc *time.Location, payload modelDomain.Domain, data modelSsl.SSL, detailsDomain []modelDomain.DetailDomain, changeServer bool) bool {
	hours := DiffHours(loc, payload)
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
func DiffHours(loc *time.Location, payload modelDomain.Domain) float64 {
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
func ValidateURL(address string) string {
	space := regexp.MustCompile(`^(http:\/\/www\.|https:\/\/www\.|http:\/\/|https:\/\/)`)
	address = space.ReplaceAllString(address, "")
	return address
}
