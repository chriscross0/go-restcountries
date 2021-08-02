Go Rest Countries
=================

[![Build Status](https://travis-ci.com/chriscross0/go-restcountries.svg?branch=master)](https://travis-ci.org/chriscross0/go-restcountries)
[![Coverage Status](https://coveralls.io/repos/github/chriscross0/go-restcountries/badge.svg?branch=master)](https://coveralls.io/github/chriscross0/go-restcountries?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/chriscross0/go-restcountries)](https://goreportcard.com/report/github.com/chriscross0/go-restcountries)

go-restcountries is a wrapper for the restcountries API, written in Go.

## Supported API methods

- All
- Name

## Usage

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
		fmt.Println("First country capitcal: ", countries[0].Capital) // Kabul
	}
}

```

## Fields Filtering

By default, all fields are returned from the API and populated to the Country type. Below is how to specify a whitelist of fields you would like and all others will not be returned. The `Fields` property is supported on the `All()` and `Name()` methods, which return a slice of countries.

```go
// Get all countries with fields filter, to include only the country Name and Capital
countries, err := client.All(restcountries.AllOptions{
	Fields: []string{"Name", "Capital"},
})

fmt.Println(countries[0].Name) // Afghanistan
fmt.Println(countries[0].Capital) // Kabul
fmt.Println(countries[0].Region) // empty because this field was not requested
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
