package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"generals-quarters", "/rooms/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"majors-suite", "/rooms/majors-suite", "GET", []postData{}, http.StatusOK},
	{"search-availability", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"make-res", "/make-reservation", "GET", []postData{}, http.StatusOK},
	{"post-search-avail", "/search-availability", "POST", []postData{
		{key: "start", value: "2024-05-08"}, {key: "end", value: "2024-05-10"},
	}, http.StatusOK},
	{"post-search-avail-json", "/search-availability-json", "POST", []postData{
		{key: "start", value: "2024-05-08"}, {key: "end", value: "2024-05-10"},
	}, http.StatusOK},
	{"make reservation post", "/make-reservation", "POST", []postData{
		{key: "firstName", value: "Yahia"}, {key: "lastName", value: "Salah"}, {key: "email", value: "a@abc.com"}, {key: "phone", value: "0123456789"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()

	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		t.Logf("Running test %s", e.name)
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("For %s, expected status code is %d, but found %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		} else if e.method == "POST" {
			values := url.Values{}
			for _, x := range e.params {
				values.Add(x.key, x.value)
			}

			resp, err := ts.Client().PostForm(ts.URL+e.url, values)
			if err != nil {
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("For %s, expected status code is %d, but found %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}
