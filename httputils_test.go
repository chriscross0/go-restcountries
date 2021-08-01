package restcountries

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type ClientMock struct {
	DoFunc func(*http.Request) (*http.Response, error)
}

func (c *ClientMock) Do(req *http.Request) (*http.Response, error) {
	return c.DoFunc(req)
}

func TestGetUrlContentErrorEOF(t *testing.T) {

	var mockedClient = &ClientMock{}
	mockedClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{}, errors.New("unexpected EOF")
	}

	gotContent, gotErr := getUrlContent("", mockedClient)

	wantContent := ""
	wantErr := "unexpected EOF"

	if gotContent != wantContent {
		t.Errorf("got content %s; wanted %s", gotContent, wantContent)
	}

	if gotErr == nil || gotErr.Error() != wantErr {
		t.Errorf("got err %v; wanted %s", gotErr, wantErr)
	}
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test read error")
}
func TestGetUrlContentErrorRead(t *testing.T) {

	// a server which returns no response but content-length:1 to cause an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "1")
	}))
	defer server.Close()

	var myClient = &http.Client{Timeout: 10 * time.Second}
	gotContent, gotErr := getUrlContent(server.URL, myClient)

	wantContent := ""
	wantErr := "unexpected EOF"

	if gotContent != wantContent {
		t.Errorf("got content %s; wanted %s", gotContent, wantContent)
	}

	if gotErr == nil || gotErr.Error() != wantErr {
		t.Errorf("got err %v; wanted %s", gotErr, wantErr)
	}
}
