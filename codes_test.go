package restcountries

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestCodesSimple(t *testing.T) {
	testClient := New("TEST_API_KEY")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `[{"name":"Colombia", "capital": "Bogotá"}]`)
	}))
	defer server.Close()
	testClient.SetApiRoot(server.URL)

	result, _ := testClient.Codes(CodesOptions{
		Fields: []string{"Name", "Capital"},
		Codes:  []string{"CO"},
	})

	got := result[0].Name
	want := "Colombia"

	if got != want {
		t.Fatalf("got %s; want %s", got, want)
	}
}

func TestCodesErrorUrl(t *testing.T) {
	testClient := New("TEST_API_KEY")

	testClient.SetApiRoot("not a url")

	_, gotErr := testClient.Codes(CodesOptions{
		Codes: []string{"CO"},
	})

	wantErr := `Get "not%20a%20url/alpha/?access_key=TEST_API_KEY&codes=CO%3B&fields=": unsupported protocol scheme ""`

	if gotErr == nil || gotErr.Error() != wantErr {
		t.Fatalf("got %s; want %s", gotErr, wantErr)
	}
}

func TestCodes(t *testing.T) {
	testClient := New("TEST_API_KEY")

	tests := []struct {
		input             CodesOptions
		response          string
		want              []Country
		wantErr           string
		checkTypeAndEmpty string
	}{
		{
			// empty search term
			input: CodesOptions{
				Fields: []string{"Name", "Capital"},
			},
			response: ``,
			wantErr:  `Search term is empty`,
		},
		{
			// not found
			input: CodesOptions{
				Fields: []string{"Name", "Capital"},
				Codes:  []string{"ABC"},
			},
			response:          `{"status": 400, "message": "Bad Request"}`,
			want:              []Country{},
			checkTypeAndEmpty: "[]restcountries.Country",
		},
		{
			// invalid json
			input: CodesOptions{
				Fields: []string{"Name", "Capital"},
				Codes:  []string{"CO"},
			},
			response: `[{"name""Colom`,
			wantErr:  `invalid character '"' after object key`,
		},
		{
			// custom error
			input: CodesOptions{
				Fields: []string{"Name", "Capital"},
				Codes:  []string{"CO"},
			},
			response: `{"status": 401, "message": "Custom Message"}`,
			wantErr:  `Custom Message`,
		},
		{
			// single country
			input: CodesOptions{
				Fields: []string{"Name", "Capital"},
				Codes:  []string{"CO"},
			},
			response: `[{"name":"Colombia", "capital": "Bogotá"}]`,
			want: []Country{
				{Name: "Colombia", Capital: "Bogotá"},
			},
		},
		{
			// multiple countries
			input: CodesOptions{
				Fields: []string{"Name", "Capital"},
				Codes:  []string{"CO", "US"},
			},
			response: `[{"name":"Colombia", "capital": "Bogotá"}, {"name":"United States of America", "capital": "Washington, D.C"}]`,
			want: []Country{
				{Name: "Colombia", Capital: "Bogotá"},
				{Name: "United States of America", Capital: "Washington, D.C"},
			},
		},
	}

	for _, test := range tests {

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, test.response)
		}))
		defer server.Close()
		testClient.SetApiRoot(server.URL)

		got, err := testClient.Codes(test.input)

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
