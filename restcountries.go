package restcountries

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"
)

// Country represents a Country from the API
// A slice of Country is returned by the methods which return countries, e.g. All and Name
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

// RestCountries represents an app/client using the API
type RestCountries struct {
	apiRoot string
}

// httpClient is used for mocking the http client
type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// AllOptions represents options for the All() method
type AllOptions struct {
	Fields []string
}

// NameOptions represents options for the Name() method
type NameOptions struct {
	Name     string
	FullText bool
	Fields   []string
}

// CapitalOptions represents options for the Capital() method
type CapitalOptions struct {
	Capital string
	Fields  []string
}

// CurrencyOptions represents options for the Currency() method
type CurrencyOptions struct {
	Currency string
	Fields   []string
}

// LanguageOptions represents options for the Language() method
type LanguageOptions struct {
	Language string
	Fields   []string
}

// New creates and returns a new instance of the client
func New() *RestCountries {
	return &RestCountries{
		apiRoot: "https://restcountries.eu/rest/v2",
	}
}

// SetApiRoot overrides the API root url - used for unit testing
func (r *RestCountries) SetApiRoot(url string) {
	r.apiRoot = url
}

// All method returns all countries
// The optional AllOptions.Fields allows filtering fields by specifying the fields you want, instead of all fields
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
		}
		return nil, errors.New(basicResponse.Message)

	}

	return countries, nil
}

// Name method searches countries by name
// The optional NameOptions.FullText boolean when true, will search for an exact match. Otherwise, partial matches are returned
// The optional NameOptions.Fields allows filtering fields by specifying the fields you want, instead of all fields
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
		}
		return nil, errors.New(basicResponse.Message)

	}

	return countries, nil
}

// Capital method searches countries by capital city using a partial match
// The optional CapitalOptions.Fields allows filtering fields by specifying the fields you want, instead of all fields
func (r *RestCountries) Capital(options CapitalOptions) ([]Country, error) {

	fields := processFields(options.Fields)

	base, _ := url.Parse(r.apiRoot)

	base.Path += "/capital/" + options.Capital // this encodes the user input properly with %20 for space and others

	params := url.Values{}
	params.Add("fields", fields)
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
		}
		return nil, errors.New(basicResponse.Message)

	}

	return countries, nil
}

// Currency method searches countries by currency code using an exact match
// The optional CurrencyOptions.Fields allows filtering fields by specifying the fields you want, instead of all fields
func (r *RestCountries) Currency(options CurrencyOptions) ([]Country, error) {

	fields := processFields(options.Fields)

	base, _ := url.Parse(r.apiRoot)

	base.Path += "/currency/" + options.Currency // this encodes the user input properly with %20 for space and others

	params := url.Values{}
	params.Add("fields", fields)
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

		if basicResponse.Status == 404 || basicResponse.Status == 400 { // 400 is returned for invalid search values
			return countries, nil
		}
		return nil, errors.New(basicResponse.Message)

	}

	return countries, nil
}

// Language method searches countries by language code using an exact match
// The optional LanguageOptions.Fields allows filtering fields by specifying the fields you want, instead of all fields
func (r *RestCountries) Language(options LanguageOptions) ([]Country, error) {

	fields := processFields(options.Fields)

	base, _ := url.Parse(r.apiRoot)

	base.Path += "/lang/" + options.Language // this encodes the user input properly with %20 for space and others

	params := url.Values{}
	params.Add("fields", fields)
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
		}
		return nil, errors.New(basicResponse.Message)

	}

	return countries, nil
}
