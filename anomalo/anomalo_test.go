package anomalo

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	fakeAnomalo = Client{
		Host:           "",
		Token:          "",
		ClientProvider: func() *http.Client { return &http.Client{} },
	}
)

func TestGetParameters(t *testing.T) {
	server := setupServer(
		t,
		"get_table_information?table_name=fake_table_name",
		`{"description": "Fake Description"}"`,
		http.StatusOK,
	)
	defer server.Close()

	fakeAnomalo.Host = server.URL
	resp, err := fakeAnomalo.GetTableInformation("fake_table_name")
	assert.Nil(t, err)
	assert.Equal(t, "Fake Description", resp.Description)
}

func TestPostBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, fmt.Sprintf("/api/public/v1/%s", "configure_table"), r.RequestURI)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		bodyBytes, err := io.ReadAll(r.Body)
		assert.Nil(t, err)
		assert.Equal(t, `{"table_id":123,"check_cadence_type":"Daily","definition":"defn"}`, string(bodyBytes))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"name": "table_name"}`))
	}))

	defer server.Close()

	checkCadence := "Daily"
	req := ConfigureTableRequest{
		TableID:          123,
		Definition:       "defn",
		CheckCadenceType: &checkCadence,
	}

	fakeAnomalo.Host = server.URL
	resp, err := fakeAnomalo.ConfigureTable(req)
	assert.Nil(t, err)
	assert.Equal(t, "table_name", resp.Name)
}

func TestAnomaloBadRequest(t *testing.T) {
	server := setupServer(t, "ping", `["API Error", "API Error 2"]`, http.StatusBadRequest)
	defer server.Close()

	fakeAnomalo.Host = server.URL
	_, err := fakeAnomalo.Ping()
	assert.NotNil(t, err)
	assert.Equal(t, `["API Error", "API Error 2"]`, err.Error())
}

func TestHttpError(t *testing.T) {
	server := setupServer(t, "ping", `{"ping": "pong"}"`, http.StatusOK) // Valid server
	server.Close()                                                       // No defer

	fakeAnomalo.Host = server.URL
	_, err := fakeAnomalo.Ping()
	assert.NotNil(t, err)
}

func TestLoadClientNoCreds(t *testing.T) {
	client, err := CreateClient()
	assert.Nil(t, client)
	assert.Errorf(t, err, "could not find anomalo credentials. exiting")
}

func TestLoadClientWithEnvVars(t *testing.T) {
	t.Setenv("ANOMALO_API_SECRET_TOKEN", "token")
	t.Setenv("ANOMALO_INSTANCE_HOST", "host")

	client, err := CreateClient()

	assert.Nil(t, err)
	assert.Equal(t, "token", client.Token)
	assert.Equal(t, "host", client.Host)

	server := setupServer(t, "ping", `{"ping": "pong"}"`, http.StatusOK)
	defer server.Close()

	fakeAnomalo.Host = server.URL
	resp, err := fakeAnomalo.Ping()
	assert.Nil(t, err)
	assert.Equal(t, resp.Ping, "pong")
}

func setupServer(t *testing.T, expectedEndpoint string, responseJson string, status int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, fmt.Sprintf("/api/public/v1/%s", expectedEndpoint), r.RequestURI)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(status)
		w.Write([]byte(responseJson))
	}))
}
