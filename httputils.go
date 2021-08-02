package restcountries

import (
	"io/ioutil"
	"net/http"
)

// getUrlContent takes a url and http client (for mock testing) and makes a GET request, returning the response text and error
func getUrlContent(url string, myClient httpClient) (string, error) {
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
