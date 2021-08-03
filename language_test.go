package restcountries

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestLanguageSimple(t *testing.T) {
	testClient := New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `[{"name":"American Samoa", "capital": "Pago Pago"}]`)
	}))
	defer server.Close()
	testClient.SetApiRoot(server.URL)

	result, _ := testClient.Language(LanguageOptions{
		Fields:   []string{"Name", "Capital"},
		Language: "EN",
	})

	got := result[0].Name
	want := "American Samoa"

	if got != want {
		t.Fatalf("got %s; want %s", got, want)
	}
}

func TestLanguageErrorUrl(t *testing.T) {
	testClient := New()

	testClient.SetApiRoot("not a url")

	_, gotErr := testClient.Language(LanguageOptions{
		Language: "EN",
	})

	wantErr := `Get "not%20a%20url/lang/EN?fields=": unsupported protocol scheme ""`

	if gotErr == nil || gotErr.Error() != wantErr {
		t.Fatalf("got %s; want %s", gotErr, wantErr)
	}
}

func TestLanguage(t *testing.T) {
	testClient := New()

	tests := []struct {
		input             LanguageOptions
		response          string
		want              []Country
		wantErr           string
		checkTypeAndEmpty string
	}{
		{
			// empty search term
			input: LanguageOptions{
				Fields:   []string{"Name", "Capital"},
				Language: "",
			},
			response: ``,
			wantErr:  `Search term is empty`,
		},
		{
			// not found
			input: LanguageOptions{
				Fields:   []string{"Name", "Capital"},
				Language: "ABC",
			},
			response:          `{"status": 404, "message": "Not Found"}`,
			want:              []Country{},
			checkTypeAndEmpty: "[]restcountries.Country",
		},
		{
			// invalid json
			input: LanguageOptions{
				Fields:   []string{"Name", "Capital"},
				Language: "EN",
			},
			response: `[{"name""Americ`,
			wantErr:  `invalid character '"' after object key`,
		},
		{
			// custom error
			input: LanguageOptions{
				Fields:   []string{"Name", "Capital"},
				Language: "EN",
			},
			response: `{"status": 500, "message": "Custom Message"}`,
			wantErr:  `Custom Message`,
		},
		{
			// single country
			input: LanguageOptions{
				Fields:   []string{"Name", "Capital"},
				Language: "TG",
			},
			response: `[{"name":"Tajikistan", "capital": "Dushanbe"}]`,
			want: []Country{
				{Name: "Tajikistan", Capital: "Dushanbe"},
			},
		},
		{
			// multiple countries
			input: LanguageOptions{
				Fields:   []string{"Name", "Capital"},
				Language: "FF",
			},
			response: `[{"name":"Burkina Faso", "capital": "Ouagadougou"}, {"name":"Guinea", "capital": "Conakry"}]`,
			want: []Country{
				{Name: "Burkina Faso", Capital: "Ouagadougou"},
				{Name: "Guinea", Capital: "Conakry"},
			},
		},
	}

	for _, test := range tests {

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, test.response)
		}))
		defer server.Close()
		testClient.SetApiRoot(server.URL)

		got, err := testClient.Language(test.input)

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
