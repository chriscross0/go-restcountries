package restcountries

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestRegionalBlocSimple(t *testing.T) {
	testClient := New("TEST_API_KEY")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `[{"name":"Chile", "capital": "Santiago"}]`)
	}))
	defer server.Close()
	testClient.SetApiRoot(server.URL)

	result, _ := testClient.RegionalBloc(RegionalBlocOptions{
		Fields:       []string{"Name", "Capital"},
		RegionalBloc: "PA",
	})

	got := result[0].Name
	want := "Chile"

	if got != want {
		t.Fatalf("got %s; want %s", got, want)
	}
}

func TestRegionalBlocErrorUrl(t *testing.T) {
	testClient := New("TEST_API_KEY")

	testClient.SetApiRoot("not a url")

	_, gotErr := testClient.RegionalBloc(RegionalBlocOptions{
		RegionalBloc: "PA",
	})

	wantErr := `Get "not%20a%20url/regionalbloc/PA?access_key=TEST_API_KEY&fields=": unsupported protocol scheme ""`

	if gotErr == nil || gotErr.Error() != wantErr {
		t.Fatalf("got %s; want %s", gotErr, wantErr)
	}
}

func TestRegionalBloc(t *testing.T) {
	testClient := New("TEST_API_KEY")

	tests := []struct {
		input             RegionalBlocOptions
		response          string
		want              []Country
		wantErr           string
		checkTypeAndEmpty string
	}{
		{
			// empty search term
			input: RegionalBlocOptions{
				Fields:       []string{"Name", "Capital"},
				RegionalBloc: "",
			},
			response: ``,
			wantErr:  `Search term is empty`,
		},
		{
			// not found
			input: RegionalBlocOptions{
				Fields:       []string{"Name", "Capital"},
				RegionalBloc: "ABC",
			},
			response:          `{"status": 404, "message": "Not Found"}`,
			want:              []Country{},
			checkTypeAndEmpty: "[]restcountries.Country",
		},
		{
			// invalid json
			input: RegionalBlocOptions{
				Fields:       []string{"Name", "Capital"},
				RegionalBloc: "PA",
			},
			response: `[{"name""Chi`,
			wantErr:  `invalid character '"' after object key`,
		},
		{
			// custom error
			input: RegionalBlocOptions{
				Fields:       []string{"Name", "Capital"},
				RegionalBloc: "PA",
			},
			response: `{"status": 500, "message": "Custom Message"}`,
			wantErr:  `Custom Message`,
		},
		{
			// multiple countries
			input: RegionalBlocOptions{
				Fields:       []string{"Name", "Capital"},
				RegionalBloc: "PA",
			},
			response: `[{"name":"Chile","capital":"Santiago"},{"name":"Colombia","capital":"Bogotá"},{"name":"Mexico","capital":"Mexico City"}]`,
			want: []Country{
				{Name: "Chile", Capital: "Santiago"},
				{Name: "Colombia", Capital: "Bogotá"},
				{Name: "Mexico", Capital: "Mexico City"},
			},
		},
	}

	for _, test := range tests {

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, test.response)
		}))
		defer server.Close()
		testClient.SetApiRoot(server.URL)

		got, err := testClient.RegionalBloc(test.input)

		// checking for an error
		if test.wantErr != "" {
			if err.Error() != test.wantErr {
				t.Fatalf("want err: %s, got: %s", test.wantErr, err.Error())
			}
			continue
		}

		if test.checkTypeAndEmpty != "" { // used for checking the slice type is correct and it is empty
			typeStr := reflect.TypeOf(got).String()
			if typeStr != test.checkTypeAndEmpty {
				t.Fatalf("want: empty %v, got: %v", test.checkTypeAndEmpty, typeStr)
			} else if len(got) != 0 {
				t.Fatalf("want: empty %v with length 0, got: %v with length %d", test.checkTypeAndEmpty, typeStr, len(got))
			}
			continue
		}

		if !reflect.DeepEqual(test.want, got) {
			t.Fatalf("want: %v, got: %v", test.want, got)
		}

	}

}
