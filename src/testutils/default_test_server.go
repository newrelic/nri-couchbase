package testutils

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// GetTestServer returns a test server mocked to return data from a file specified in the endpoint map
func GetTestServer(t *testing.T, dataMap map[string]string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		username, password, ok := req.BasicAuth()
		assert.True(t, ok)
		assert.Equal(t, username, "testUser")
		assert.Equal(t, password, "testPass")

		endpoint := req.RequestURI
		filepath, ok := dataMap[endpoint]
		if !ok {
			t.Errorf("bad request, was not expecting request for endpoint %s", endpoint)
		}
		data, _ := ioutil.ReadFile(filepath)
		_, err := res.Write(data)
		if err != nil {
			t.Errorf("could not write response body for endpoint %s: %v", endpoint, err)
		}
	}))
}
