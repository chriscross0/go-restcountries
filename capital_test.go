package restcountries

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestCapitalSimple(t *testing.T) {
	testClient := New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `[{"name":"France", "capital": "Paris"}]`)
	}))
	defer server.Close()
	testClient.SetApiRoot(server.URL)

	result, _ := testClient.Capital(CapitalOptions{
		Fields:  []string{"Name", "Capital"},
		Capital: "Paris",
	})

	got := result[0].Name
	want := "France"

	if got != want {
		t.Fatalf("got %s; want %s", got, want)
	}
}

func TestCapitalErrorUrl(t *testing.T) {
	testClient := New()

	testClient.SetApiRoot("not a url")

	_, gotErr := testClient.Capital(CapitalOptions{
		Capital: "Paris",
	})

	wantErr := `Get "not%20a%20url/capital/Paris?fields=": unsupported protocol scheme ""`

	if gotErr == nil || gotErr.Error() != wantErr {
		t.Fatalf("got %s; want %s", gotErr, wantErr)
	}
}

func TestCapital(t *testing.T) {
	testClient := New()

	tests := []struct {
		input             CapitalOptions
		response          string
		want              []Country
		wantErr           string
		checkTypeAndEmpty string
	}{
		{
			// not found
			input: CapitalOptions{
				Fields:  []string{"Name", "Capital"},
				Capital: "ABCDEF",
			},
			response:          `{"status": 404, "message": "Not Found"}`,
			want:              []Country{},
			checkTypeAndEmpty: "[]restcountries.Country",
		},
		{
			// invalid json
			input: CapitalOptions{
				Fields:  []string{"Name", "Capital"},
				Capital: "Paris",
			},
			response: `[{"name""Fran`,
			wantErr:  `invalid character '"' after object key`,
		},
		{
			// custom error
			input: CapitalOptions{
				Fields:  []string{"Name", "Capital"},
				Capital: "Paris",
			},
			response: `{"status": 500, "message": "Custom Message"}`,
			wantErr:  `Custom Message`,
		},
		{
			// single country
			input: CapitalOptions{
				Fields:  []string{"Name", "Capital"},
				Capital: "Paris",
			},
			response: `[{"name":"France", "capital": "Paris"}]`,
			want: []Country{
				{Name: "France", Capital: "Paris"},
			},
		},
		{
			// multiple countries
			input: CapitalOptions{
				Fields:  []string{"Name", "Capital"},
				Capital: "Lon",
			},
			response: `[{"name":"Malawi", "capital": "Lilongwe"}, {"name":"Svalbard and Jan Mayen", "capital": "Longyearbyen"}, {"name":"United Kingdom of Great Britain and Northern Ireland", "capital": "London"}]`,
			want: []Country{
				{Name: "Malawi", Capital: "Lilongwe"},
				{Name: "Svalbard and Jan Mayen", Capital: "Longyearbyen"},
				{Name: "United Kingdom of Great Britain and Northern Ireland", Capital: "London"},
			},
		},
	}

	for _, test := range tests {

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, test.response)
		}))
		defer server.Close()
		testClient.SetApiRoot(server.URL)

		got, err := testClient.Capital(test.input)

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
