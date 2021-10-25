package restcountries

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestCurrencySimple(t *testing.T) {
	testClient := New("TEST_API_KEY")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `[{"name":"Antarctica", "capital": ""}]`)
	}))
	defer server.Close()
	testClient.SetApiRoot(server.URL)

	result, _ := testClient.Currency(CurrencyOptions{
		Fields:   []string{"Name", "Capital"},
		Currency: "GBP",
	})

	got := result[0].Name
	want := "Antarctica"

	if got != want {
		t.Fatalf("got %s; want %s", got, want)
	}
}

func TestCurrencyErrorUrl(t *testing.T) {
	testClient := New("TEST_API_KEY")

	testClient.SetApiRoot("not a url")

	_, gotErr := testClient.Currency(CurrencyOptions{
		Currency: "GBP",
	})

	wantErr := `Get "not%20a%20url/currency/GBP?access_key=TEST_API_KEY&fields=": unsupported protocol scheme ""`

	if gotErr == nil || gotErr.Error() != wantErr {
		t.Fatalf("got %s; want %s", gotErr, wantErr)
	}
}

func TestCurrency(t *testing.T) {
	testClient := New("TEST_API_KEY")

	tests := []struct {
		input             CurrencyOptions
		response          string
		want              []Country
		wantErr           string
		checkTypeAndEmpty string
	}{
		{
			// empty search term
			input: CurrencyOptions{
				Fields:   []string{"Name", "Capital"},
				Currency: "",
			},
			response: ``,
			wantErr:  `Search term is empty`,
		},
		{
			// not found
			input: CurrencyOptions{
				Fields:   []string{"Name", "Capital"},
				Currency: "ABC",
			},
			response:          `{"status": 404, "message": "Not Found"}`,
			want:              []Country{},
			checkTypeAndEmpty: "[]restcountries.Country",
		},
		{
			// invalid currency
			input: CurrencyOptions{
				Fields:   []string{"Name", "Capital"},
				Currency: "ABCDEF",
			},
			response:          `{"status": 400, "message": "Bad Request"}`,
			want:              []Country{},
			checkTypeAndEmpty: "[]restcountries.Country",
		},
		{
			// invalid json
			input: CurrencyOptions{
				Fields:   []string{"Name", "Capital"},
				Currency: "GBP",
			},
			response: `[{"name""Antarc`,
			wantErr:  `invalid character '"' after object key`,
		},
		{
			// custom error
			input: CurrencyOptions{
				Fields:   []string{"Name", "Capital"},
				Currency: "GBP",
			},
			response: `{"status": 500, "message": "Custom Message"}`,
			wantErr:  `Custom Message`,
		},
		{
			// single country
			input: CurrencyOptions{
				Fields:   []string{"Name", "Capital"},
				Currency: "IDR",
			},
			response: `[{"name":"Indonesia", "capital": "Jakarta"}]`,
			want: []Country{
				{Name: "Indonesia", Capital: "Jakarta"},
			},
		},
		{
			// multiple countries
			input: CurrencyOptions{
				Fields:   []string{"Name", "Capital"},
				Currency: "SGD",
			},
			response: `[{"name":"Brunei Darussalam", "capital": "Bandar Seri Begawan"}, {"name":"Singapore", "capital": "Singapore"}]`,
			want: []Country{
				{Name: "Brunei Darussalam", Capital: "Bandar Seri Begawan"},
				{Name: "Singapore", Capital: "Singapore"},
			},
		},
	}

	for _, test := range tests {

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, test.response)
		}))
		defer server.Close()
		testClient.SetApiRoot(server.URL)

		got, err := testClient.Currency(test.input)

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
