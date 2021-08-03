package restcountries

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestCallingCodeSimple(t *testing.T) {
	testClient := New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `[{"name":"Estonia", "capital": "Tallinn"}]`)
	}))
	defer server.Close()
	testClient.SetApiRoot(server.URL)

	result, _ := testClient.CallingCode(CallingCodeOptions{
		Fields:      []string{"Name", "Capital"},
		CallingCode: "372",
	})

	got := result[0].Name
	want := "Estonia"

	if got != want {
		t.Fatalf("got %s; want %s", got, want)
	}
}

func TestCallingCodeErrorUrl(t *testing.T) {
	testClient := New()

	testClient.SetApiRoot("not a url")

	_, gotErr := testClient.CallingCode(CallingCodeOptions{
		CallingCode: "372",
	})

	wantErr := `Get "not%20a%20url/callingcode/372?fields=": unsupported protocol scheme ""`

	if gotErr == nil || gotErr.Error() != wantErr {
		t.Fatalf("got %s; want %s", gotErr, wantErr)
	}
}

func TestCallingCode(t *testing.T) {
	testClient := New()

	tests := []struct {
		input             CallingCodeOptions
		response          string
		want              []Country
		wantErr           string
		checkTypeAndEmpty string
	}{
		{
			// empty search term
			input: CallingCodeOptions{
				Fields:      []string{"Name", "Capital"},
				CallingCode: "",
			},
			response: ``,
			wantErr:  `Search term is empty`,
		},
		{
			// not found
			input: CallingCodeOptions{
				Fields:      []string{"Name", "Capital"},
				CallingCode: "ABC",
			},
			response:          `{"status": 404, "message": "Not Found"}`,
			want:              []Country{},
			checkTypeAndEmpty: "[]restcountries.Country",
		},
		{
			// invalid json
			input: CallingCodeOptions{
				Fields:      []string{"Name", "Capital"},
				CallingCode: "372",
			},
			response: `[{"name""Eston`,
			wantErr:  `invalid character '"' after object key`,
		},
		{
			// custom error
			input: CallingCodeOptions{
				Fields:      []string{"Name", "Capital"},
				CallingCode: "372",
			},
			response: `{"status": 500, "message": "Custom Message"}`,
			wantErr:  `Custom Message`,
		},
		{
			// single country
			input: CallingCodeOptions{
				Fields:      []string{"Name", "Capital"},
				CallingCode: "372",
			},
			response: `[{"name":"Estonia", "capital": "Tallinn"}]`,
			want: []Country{
				{Name: "Estonia", Capital: "Tallinn"},
			},
		},
		{
			// multiple countries
			input: CallingCodeOptions{
				Fields:      []string{"Name", "Capital"},
				CallingCode: "44",
			},
			response: `[{"name":"Guernsey","capital":"St. Peter Port"},{"name":"Isle of Man","capital":"Douglas"},{"name":"Jersey","capital":"Saint Helier"}]`,
			want: []Country{
				{Name: "Guernsey", Capital: "St. Peter Port"},
				{Name: "Isle of Man", Capital: "Douglas"},
				{Name: "Jersey", Capital: "Saint Helier"},
			},
		},
	}

	for _, test := range tests {

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, test.response)
		}))
		defer server.Close()
		testClient.SetApiRoot(server.URL)

		got, err := testClient.CallingCode(test.input)

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
