package command

import (
	"encoding/json"
	"net/http"
	"regexp"
	"sort"
	"time"

	modelDomain "go-scraping/model/domain"
	modelSsl "go-scraping/model/ssllabs"
)

// RespondWithJSON write json response format
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// RespondWithError return error message
func RespondWithError(w http.ResponseWriter, code int, msg string) {
	RespondWithJSON(w, code, map[string]string{"message": msg})
}

// GetLowestGradeCurrent We get the lowest grade from current servers
func GetLowestGradeCurrent(data []modelSsl.Endpoint) string {
	var gradeASCII []int
	var grade string

	for _, dataElement := range data {
		if dataElement.Grade != "A+" {
			gradeASCII = append(gradeASCII, int(dataElement.Grade[0]))
		} else {
			grade = "A+"
		}
	}

	if len(gradeASCII) > 0 {
		sort.Slice(gradeASCII, func(i, j int) bool {
			return gradeASCII[i] > gradeASCII[j]
		})

		grade = string(gradeASCII[0])
	}

	return grade
}

// GetLowestGradePrevious  We get the lowest grade from previous servers
func GetLowestGradePrevious(detail []modelDomain.DetailDomain) string {
	var gradeASCII []int
	var grade string

	for _, dataElement := range detail {
		if dataElement.Grade != "A+" {
			gradeASCII = append(gradeASCII, int(dataElement.Grade[0]))
		} else {
			grade = "A+"
		}
	}

	if len(gradeASCII) > 0 {
		sort.Slice(gradeASCII, func(i, j int) bool {
			return gradeASCII[i] > gradeASCII[j]
		})

		grade = string(gradeASCII[0])
	}

	return grade
}

// ValidateChangeServer Validate if there is a change in the main data of the servers
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

// DiffHours Total hours difference
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

// ValidateURL Validate that the URL is cleaned
func ValidateURL(address string) string {
	space := regexp.MustCompile(`^(http:\/\/www\.|https:\/\/www\.|http:\/\/|https:\/\/)`)
	address = space.ReplaceAllString(address, "")
	return address
}
