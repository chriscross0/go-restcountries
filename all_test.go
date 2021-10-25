package restcountries

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestAllSimple(t *testing.T) {
	testClient := New("TEST_API_KEY")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `[{"name":"TestName", "capital": "testCap"}, {"name":"TestName2", "capital": "testCap2"}]`)
	}))
	defer server.Close()
	testClient.SetApiRoot(server.URL)
	testClient.SetTimeout(10 * time.Second)

	result, _ := testClient.All(AllOptions{
		Fields: []string{"Name", "Capital"},
	})

	got := len(result)
	want := 2

	if got != want {
		t.Fatalf("got len %d; wanted %d", got, want)
	}
}

func TestAllErrorUrl(t *testing.T) {
	testClient := New("TEST_API_KEY")

	testClient.SetApiRoot("not a url")

	_, gotErr := testClient.All(AllOptions{})

	wantErr := `Get "not%20a%20url/all?access_key=TEST_API_KEY&fields=": unsupported protocol scheme ""`

	if gotErr == nil || gotErr.Error() != wantErr {
		t.Fatalf("got %s; want %s", gotErr, wantErr)
	}
}

func TestAll(t *testing.T) {

	testClient := New("TEST_API_KEY")

	tests := []struct {
		input             AllOptions
		response          string
		want              []Country
		wantErr           string
		checkTypeAndEmpty string
	}{
		{
			// not found
			input: AllOptions{
				Fields: []string{"Name", "Capital"},
			},
			response:          `{"status": 404, "message": "Not Found"}`,
			want:              []Country{},
			checkTypeAndEmpty: "[]restcountries.Country",
		},
		{
			// invalid json
			input: AllOptions{
				Fields: []string{"Name", "Capital"},
			},
			response: `[{"name""Fran`,
			wantErr:  `invalid character '"' after object key`,
		},
		{
			// custom error
			input: AllOptions{
				Fields: []string{"Name", "Capital"},
			},
			response: `{"status": 500, "message": "Custom Message"}`,
			wantErr:  `Custom Message`,
		},
		{
			// single country
			input: AllOptions{
				Fields: []string{"Name", "Capital"},
			},
			response: `[{"name":"France", "capital": "Paris"}]`,
			want: []Country{
				{Name: "France", Capital: "Paris"},
			},
		},
		{
			// multiple countries
			input: AllOptions{
				Fields: []string{"Name", "Capital"},
			},
			response: `[{"name":"United States of America", "capital": "Washington DC"}, {"name":"United Arab Emirates", "capital": "Abu Dhabi"}]`,
			want: []Country{
				{Name: "United States of America", Capital: "Washington DC"},
				{Name: "United Arab Emirates", Capital: "Abu Dhabi"},
			},
		},
	}

	for _, test := range tests {

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, test.response)
		}))
		defer server.Close()
		testClient.SetApiRoot(server.URL)

		got, err := testClient.All(test.input)

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
