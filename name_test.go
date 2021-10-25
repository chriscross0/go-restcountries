package restcountries

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNameSimple(t *testing.T) {
	testClient := New("TEST_API_KEY")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `[{"name":"France", "capital": "Paris"}]`)
	}))
	defer server.Close()
	testClient.SetApiRoot(server.URL)

	result, _ := testClient.Name(NameOptions{
		Fields: []string{"Name", "Capital"},
		Name:   "France",
	})

	got := result[0].Name
	want := "France"

	if got != want {
		t.Fatalf("got %s; want %s", got, want)
	}
}

func TestNameErrorUrl(t *testing.T) {
	testClient := New("TEST_API_KEY")

	testClient.SetApiRoot("not a url")

	_, gotErr := testClient.Name(NameOptions{
		Name: "France",
	})

	wantErr := `Get "not%20a%20url/name/France?access_key=TEST_API_KEY&fields=": unsupported protocol scheme ""`

	if gotErr == nil || gotErr.Error() != wantErr {
		t.Fatalf("got %s; want %s", gotErr, wantErr)
	}
}

func TestName(t *testing.T) {
	testClient := New("TEST_API_KEY")

	tests := []struct {
		input             NameOptions
		response          string
		want              []Country
		wantErr           string
		checkTypeAndEmpty string
	}{
		{
			// empty search term
			input: NameOptions{
				Fields: []string{"Name", "Capital"},
				Name:   "",
			},
			response: ``,
			wantErr:  `Search term is empty`,
		},
		{
			// not found
			input: NameOptions{
				Fields: []string{"Name", "Capital"},
				Name:   "ABCDEF",
			},
			response:          `{"status": 404, "message": "Not Found"}`,
			want:              []Country{},
			checkTypeAndEmpty: "[]restcountries.Country",
		},
		{
			// invalid json
			input: NameOptions{
				Fields: []string{"Name", "Capital"},
				Name:   "France",
			},
			response: `[{"name""Fran`,
			wantErr:  `invalid character '"' after object key`,
		},
		{
			// custom error
			input: NameOptions{
				Fields: []string{"Name", "Capital"},
				Name:   "France",
			},
			response: `{"status": 500, "message": "Custom Message"}`,
			wantErr:  `Custom Message`,
		},
		{
			// single country
			input: NameOptions{
				Fields: []string{"Name", "Capital"},
				Name:   "France",
			},
			response: `[{"name":"France", "capital": "Paris"}]`,
			want: []Country{
				{Name: "France", Capital: "Paris"},
			},
		},
		{
			// multiple countries
			input: NameOptions{
				Fields: []string{"Name", "Capital"},
				Name:   "United",
			},
			response: `[{"name":"United States of America", "capital": "Washington DC"}, {"name":"United Arab Emirates", "capital": "Abu Dhabi"}]`,
			want: []Country{
				{Name: "United States of America", Capital: "Washington DC"},
				{Name: "United Arab Emirates", Capital: "Abu Dhabi"},
			},
		},
		{
			// FullText on (exact name match)
			input: NameOptions{
				Fields:   []string{"Name", "Capital"},
				Name:     "France",
				FullText: true,
			},
			response: `[{"name":"France", "capital": "Paris"}]`,
			want: []Country{
				{Name: "France", Capital: "Paris"},
			},
		},
	}

	for _, test := range tests {

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, test.response)
		}))
		defer server.Close()
		testClient.SetApiRoot(server.URL)

		got, err := testClient.Name(test.input)

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
