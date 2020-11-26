package http

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	model "go-scraping/model/ssllabs"
)

// GetDataSSl obtains the SSL information from the indicated address
func GetDataSSl(address string) (model.SSL, error) {
	response, err := http.Get("https://api.ssllabs.com/api/v3/analyze?host=" + address)

	if err != nil {
		log.Println(err)
	}

	defer response.Body.Close()
	if response.StatusCode != 200 {
		log.Printf("status code error: %d %s", response.StatusCode, response.Status)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}

	var responseSSL model.SSL
	json.Unmarshal(responseData, &responseSSL)

	return responseSSL, err
}
