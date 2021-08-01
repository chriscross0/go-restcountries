package restcountries

import (
	"io/ioutil"
	"net/http"
)

// takes a url and http client (for mock testing) and returns the response text and error
func getUrlContent(url string, myClient HttpClient) (string, error) {
	req, _ := http.NewRequest("GET", url, nil)

	resp, respErr := myClient.Do(req)

	if respErr != nil {
		return "", respErr
	}

	defer resp.Body.Close()

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return "", readErr
	}

	return string(body), nil
}
