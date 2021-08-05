Go REST Countries
=================

[![GoDoc](https://godoc.org/github.com/chriscross0/go-restcountries?status.svg)](http://godoc.org/github.com/chriscross0/go-restcountries)
[![Build Status](https://travis-ci.com/chriscross0/go-restcountries.svg?branch=master)](https://travis-ci.org/chriscross0/go-restcountries)
[![Coverage Status](https://coveralls.io/repos/github/chriscross0/go-restcountries/badge.svg?branch=master)](https://coveralls.io/github/chriscross0/go-restcountries?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/chriscross0/go-restcountries)](https://goreportcard.com/report/github.com/chriscross0/go-restcountries)

go-restcountries is a wrapper for the [REST Countries API](https://restcountries.eu/), written in Go. The latest (v2) version of the API is used.

## Supported API methods (all methods of the v2 API are supported)

- All - get all countries.
- Name - search countries by name, including the option of an exact or partial match.
- Capital - search countries by capital city. Uses a partial match.
- Currency - search countries by ISO 4217 currency code. Uses an exact match.
- Language - search countries by ISO 639-1 language code. Uses an exact match.
- Region - search countries by region: Africa, Americas, Asia, Europe, Oceania. Uses an exact match.
- RegionalBloc - search countries by regional bloc: EU, EFTA, CARICOM, PA etc. Uses an exact match.
- CallingCode - search countries by calling code. Uses an exact match.
- Code/List of Codes (method name is Codes) - search countries by ISO 3166-1 2-letter or 3-letter country codes. Uses an exact match.

## Usage

### Get all countries

```go
package main

import (
	"fmt"
	"github.com/chriscross0/go-restcountries"
)

func main(){
	client := restcountries.New()

	// All with no fields filter (get all countries with all fields)
	countries, err := client.All(restcountries.AllOptions{})

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Total countries: ", len(countries)) // 250
		fmt.Println("First country name: ", countries[0].Name) // Afghanistan
		fmt.Println("First country capital: ", countries[0].Capital) // Kabul
	}
}

```

### Search countries by name - partial match

```go
countries, err := client.Name(restcountries.NameOptions{
	Name: "United States",
})

fmt.Println("Total countries: ", len(countries)) // 2
fmt.Println("First country name: ", countries[0].Name) // United States Minor Outlying Islands
fmt.Println("Second country name: ", countries[1].Name) // United States of America
```

### Search countries by name - exact match

```go
countries, err := client.Name(restcountries.NameOptions{
	Name: "United States of America",
	FullText: true, // true turns exact match on
})

fmt.Println("Total countries: ", len(countries)) // 1
fmt.Println("First country name: ", countries[0].Name) // United States of America
```

### Search countries by capital city - partial match with single country found

```go
countries, err := client.Capital(restcountries.CapitalOptions{
	Name: "London",
})

fmt.Println("Total countries: ", len(countries)) // 1
fmt.Println("First country name: ", countries[0].Name) // United Kingdom of Great Britain and Northern Ireland
```

### Search countries by capital city - partial match with multiple countries found

```go
countries, err := client.Capital(restcountries.CapitalOptions{
	Name: "Lon",
})

fmt.Println("Total countries: ", len(countries)) // 3
fmt.Println("First country name: ", countries[0].Name) // Malawi
fmt.Println("Second country name: ", countries[1].Name) // Svalbard and Jan Mayen
fmt.Println("Third country name: ", countries[2].Name) // United Kingdom of Great Britain and Northern Ireland
```

### Search countries by currency code - exact match with single country found

```go
countries, err := client.Currency(restcountries.CurrencyOptions{
	Currency: "IDR",
})

fmt.Println("Total countries: ", len(countries)) // 1
fmt.Println("First country name: ", countries[0].Name) // Indonesia
```

### Search countries by currency code - exact match with multiple countries found

```go
countries, err := client.Capital(restcountries.CurrencyOptions{
	Currency: "SGD",
})

fmt.Println("Total countries: ", len(countries)) // 2
fmt.Println("First country name: ", countries[0].Name) // Brunei Darussalam
fmt.Println("Second country name: ", countries[1].Name) // Singapore
```

### Search countries by language code - exact match with single country found

```go
countries, err := client.Language(restcountries.LanguageOptions{
	Language: "TG",
})

fmt.Println("Total countries: ", len(countries)) // 1
fmt.Println("First country name: ", countries[0].Name) // Tajikistan
```

### Search countries by language code - exact match with multiple countries found

```go
countries, err := client.Language(restcountries.LanguageOptions{
	Language: "FF",
})

fmt.Println("Total countries: ", len(countries)) // 2
fmt.Println("First country name: ", countries[0].Name) // Burkina Faso
fmt.Println("Second country name: ", countries[1].Name) // Guinea
```

### Search countries by region - exact match with multiple countries found

```go
countries, err := client.Region(restcountries.RegionOptions{
	Region: "Oceania",
})

fmt.Println("Total countries: ", len(countries)) // 27
fmt.Println("First country name: ", countries[0].Name) // American Samoa
fmt.Println("Second country name: ", countries[1].Name) // Australia
```

### Search countries by regional bloc - exact match with multiple countries found

```go
countries, err := client.RegionalBloc(restcountries.RegionalBlocOptions{
	RegionalBloc: "PA",
})

fmt.Println("Total countries: ", len(countries)) // 4
fmt.Println("First country name: ", countries[0].Name) // Chile
fmt.Println("Second country name: ", countries[1].Name) // Colombia
```

### Search countries by calling code - exact match with single country found

```go
countries, err := client.CallingCode(restcountries.CallingCodeOptions{
	CallingCode: "372",
})

fmt.Println("Total countries: ", len(countries)) // 1
fmt.Println("First country name: ", countries[0].Name) // Estonia
```

### Search countries by calling code - exact match with multiple countries found

```go
countries, err := client.CallingCode(restcountries.CallingCodeOptions{
	CallingCode: "44",
})

fmt.Println("Total countries: ", len(countries)) // 4
fmt.Println("First country name: ", countries[0].Name) // Guernsey
fmt.Println("Second country name: ", countries[1].Name) // Isle of Man
```

### Search countries by country code - exact match with single country found

```go
countries, err := client.Codes(restcountries.CodesOptions{
	Codes: []string{"CO"}, // single code
})

fmt.Println("Total countries: ", len(countries)) // 1
fmt.Println("First country name: ", countries[0].Name) // Colombia
```

### Search countries by country code - exact match with multiple countries found

```go
countries, err := client.Codes(restcountries.CodesOptions{
	Codes: []string{"CO", "GB"}, // multiple codes
})

fmt.Println("Total countries: ", len(countries)) // 2
fmt.Println("First country name: ", countries[0].Name) // Colombia
fmt.Println("Second country name: ", countries[1].Name) // United Kingdom of Great Britain and Northern Ireland
```

### Fields Filtering

By default, all fields are returned from the API and populated to the Country type. Below is how to specify a whitelist of fields you would like and all others will not be returned. The `Fields` property is supported on the `All()`, `Name()`, `Capital()`, `Currency()`, `Language()`, `Region()`, `RegionalBloc()`, `CallingCode()` and `Codes()` methods, which return a slice of countries.

```go
// Get all countries with fields filter, to include only the country Name and Capital
countries, err := client.All(restcountries.AllOptions{
	Fields: []string{"Name", "Capital"},
})

fmt.Println(countries[0].Name) // Afghanistan
fmt.Println(countries[0].Capital) // Kabul
fmt.Println(countries[0].Region) // empty because this field was not requested
```

## Configuration

### `SetTimeout()`

The default timeout for the HTTP client is `0` meaning no timeout. Use `SetTimeout()` to override the default timeout, using a time.Duration.

```go
	client := restcountries.New()
	client.SetTimeout(10 * time.Duration)
	...
```

## Supported Fields

All fields in the v2 restcountries APi are supported. Below is the Country type:

```go
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

```
