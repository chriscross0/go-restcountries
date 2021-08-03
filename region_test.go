package restcountries

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestRegionSimple(t *testing.T) {
	testClient := New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `[{"name":"American Samoa", "capital": "Pago Pago"}]`)
	}))
	defer server.Close()
	testClient.SetApiRoot(server.URL)

	result, _ := testClient.Region(RegionOptions{
		Fields: []string{"Name", "Capital"},
		Region: "Oceania",
	})

	got := result[0].Name
	want := "American Samoa"

	if got != want {
		t.Fatalf("got %s; want %s", got, want)
	}
}

func TestRegionErrorUrl(t *testing.T) {
	testClient := New()

	testClient.SetApiRoot("not a url")

	_, gotErr := testClient.Region(RegionOptions{
		Region: "Oceania",
	})

	wantErr := `Get "not%20a%20url/region/Oceania?fields=": unsupported protocol scheme ""`

	if gotErr == nil || gotErr.Error() != wantErr {
		t.Fatalf("got %s; want %s", gotErr, wantErr)
	}
}

func TestRegion(t *testing.T) {
	testClient := New()

	tests := []struct {
		input             RegionOptions
		response          string
		want              []Country
		wantErr           string
		checkTypeAndEmpty string
	}{
		{
			// empty search term
			input: RegionOptions{
				Fields: []string{"Name", "Capital"},
				Region: "",
			},
			response: ``,
			wantErr:  `Search term is empty`,
		},
		{
			// not found
			input: RegionOptions{
				Fields: []string{"Name", "Capital"},
				Region: "ABC",
			},
			response:          `{"status": 404, "message": "Not Found"}`,
			want:              []Country{},
			checkTypeAndEmpty: "[]restcountries.Country",
		},
		{
			// invalid json
			input: RegionOptions{
				Fields: []string{"Name", "Capital"},
				Region: "Oceania",
			},
			response: `[{"name""Americ`,
			wantErr:  `invalid character '"' after object key`,
		},
		{
			// custom error
			input: RegionOptions{
				Fields: []string{"Name", "Capital"},
				Region: "Oceania",
			},
			response: `{"status": 500, "message": "Custom Message"}`,
			wantErr:  `Custom Message`,
		},
		{
			// multiple countries
			input: RegionOptions{
				Fields: []string{"Name", "Capital"},
				Region: "Oceania",
			},
			response: `[{"name":"American Samoa", "capital": "Pago Pago"}, {"name":"Australia","capital":"Canberra"}, {"name":"Christmas Island","capital":"Flying Fish Cove"}]`,
			want: []Country{
				{Name: "American Samoa", Capital: "Pago Pago"},
				{Name: "Australia", Capital: "Canberra"},
				{Name: "Christmas Island", Capital: "Flying Fish Cove"},
			},
		},
	}

	for _, test := range tests {

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, test.response)
		}))
		defer server.Close()
		testClient.SetApiRoot(server.URL)

		got, err := testClient.Region(test.input)

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
