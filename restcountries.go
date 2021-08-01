package restcountries

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"
)

type Country struct {
	Name           string    `json:"name"`
	TopLevelDomain []string  `json:"topLevelDomain"`
	Alpha2Code     string    `json:"alpha2Code"`
	Alpha3Code     string    `json:"alpha3Code"`
	CallingCodes   []string  `json:"callingCodes"`
	Capital        string    `json:"capital"`
	AltSpellings   []string  `json:"altSpellings"`
	Region         string    `json:"region"`
	Subregion      string    `json:"subregion"`
	Population     int       `json:"population"`
	Latlng         []float64 `json:"latlng"`
	Demonym        string    `json:"demonym"`
	Area           float64   `json:"area"`
	Gini           float64   `json:"gini"`
	Timezones      []string  `json:"timezones"`
	Borders        []string  `json:"borders"`
	NativeName     string    `json:"nativeName"`
	NumericCode    string    `json:"numericCode"`
	Currencies     []struct {
		Code   string `json:"code"`
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
	} `json:"currencies"`
	Languages []struct {
		Iso6391    string `json:"iso639_1"`
		Iso6392    string `json:"iso639_2"`
		Name       string `json:"name"`
		NativeName string `json:"nativeName"`
	} `json:"languages"`
	Translations struct {
		De string `json:"de"`
		Es string `json:"es"`
		Fr string `json:"fr"`
		Ja string `json:"ja"`
		It string `json:"it"`
		Br string `json:"br"`
		Pt string `json:"pt"`
		Nl string `json:"nl"`
		Hr string `json:"hr"`
		Fa string `json:"fa"`
	} `json:"translations"`
	Flag          string `json:"flag"`
	RegionalBlocs []struct {
		Acronym       string   `json:"acronym"`
		Name          string   `json:"name"`
		OtherAcronyms []string `json:"otherAcronyms"`
		OtherNames    []string `json:"otherNames"`
	} `json:"regionalBlocs"`
	Cioc string `json:"cioc"`
}

type apiError struct {
	Status  int16  `json:"status"`
	Message string `json:"message"`
}

type RestCountries struct {
	apiRoot string
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type AllOptions struct {
	Fields []string
}

type NameOptions struct {
	Name     string
	FullText bool
	Fields   []string
}

func New() *RestCountries {
	return &RestCountries{
		apiRoot: "https://restcountries.eu/rest/v2",
	}
}

func (r *RestCountries) SetApiRoot(url string) {
	r.apiRoot = url
}

func (r *RestCountries) All(options AllOptions) ([]Country, error) {

	fields := processFields(options.Fields)

	var myClient = &http.Client{Timeout: 10 * time.Second}
	content, err := getUrlContent(r.apiRoot+"/all?fields="+url.QueryEscape(fields), myClient)

	if err != nil {
		return nil, err
	}

	var countries []Country
	decodeErr := json.Unmarshal([]byte(content), &countries)
	if decodeErr != nil {

		var basicResponse apiError
		basicResponseErr := json.Unmarshal([]byte(content), &basicResponse)
		if basicResponseErr != nil {
			return nil, decodeErr
		}

		if basicResponse.Status == 404 {
			return countries, nil
		} else {
			return nil, errors.New(basicResponse.Message)
		}

	}

	return countries, nil
}

func (r *RestCountries) Name(options NameOptions) ([]Country, error) {

	fields := processFields(options.Fields)

	base, _ := url.Parse(r.apiRoot)

	base.Path += "/name/" + options.Name // this encodes the user input properly with %20 for space and others

	params := url.Values{}
	params.Add("fields", fields)
	if options.FullText {
		params.Add("fullText", "true")
	}
	base.RawQuery = params.Encode()

	var myClient = &http.Client{Timeout: 10 * time.Second}
	content, err := getUrlContent(base.String(), myClient)

	if err != nil {
		return nil, err
	}

	var countries []Country
	decodeErr := json.Unmarshal([]byte(content), &countries)
	if decodeErr != nil {

		var basicResponse apiError
		basicResponseErr := json.Unmarshal([]byte(content), &basicResponse)
		if basicResponseErr != nil {
			return nil, decodeErr
		}

		if basicResponse.Status == 404 {
			return countries, nil
		} else {
			return nil, errors.New(basicResponse.Message)
		}

	}

	return countries, nil
}
