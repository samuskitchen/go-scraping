package http

import (
	model "../../model/ssllabs"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func GetDataSSl(address string) (model.SSL, error) {
	response, err := http.Get("https://api.ssllabs.com/api/v3/analyze?host=​" + address)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseSSL model.SSL
	json.Unmarshal(responseData, &responseSSL)

	return responseSSL, err
}