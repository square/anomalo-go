package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	beyondclient "github.com/squareup/beyondclient-go"
)

// TODO
//   Error propagation from all methods

type AnomaloClient struct {
	Token  string `json:"Token,omitempty"`
	Host   string `json:"Host,omitempty"`
	Client *http.Client
}

func (c AnomaloClient) getClient() *http.Client {
	if c.Client == nil {
		c.Client = beyondclient.NewClient()
	}
	return c.Client
}

func (c AnomaloClient) buildUrl(endpoint string) string {
	return fmt.Sprintf("https://%s/api/public/v1/%s", c.Host, endpoint)
}

func (c AnomaloClient) apiCall(endpoint string, method string) io.ReadCloser {
	return c.apiCallWithBody(endpoint, method, "{}")
}

func (c AnomaloClient) apiCallWithBody(endpoint string, method string, jsonParams string) io.ReadCloser {
	// TODO environment variables for Host & Token. Check for presence of Host & Token with real error message.

	var req *http.Request
	var err error
	if method == http.MethodPost || method == http.MethodPut {
		req, err = http.NewRequest(method, c.buildUrl(endpoint), bytes.NewBuffer([]byte(jsonParams)))
		if err != nil {
			fmt.Println(err)
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, c.buildUrl(endpoint), bytes.NewBuffer([]byte(jsonParams)))
		if err != nil {
			fmt.Println(err)
		}
		var parsed map[string]interface{}
		err = json.Unmarshal([]byte(jsonParams), &parsed)
		if err != nil {
			fmt.Println(err)
		}
		params := req.URL.Query()
		for key, value := range parsed {
			params.Add(key, fmt.Sprintf("%v", value))
		}
		req.URL.RawQuery = params.Encode()
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("X-Anomalo-Token", c.Token)

	r, err := c.getClient().Do(req)
	if err != nil {
		fmt.Println(err)
	}

	return r.Body
}

func (c AnomaloClient) Ping() PingResponse {
	// TODO is there a way to centralize this unmarshalling?
	var data PingResponse
	body := c.apiCall("ping", http.MethodGet)
	defer body.Close()
	//x, _ := io.ReadAll(body)
	//fmt.Println(string(x))
	if err := json.NewDecoder(body).Decode(&data); err != nil {
		fmt.Println(err)
	}
	return data
}

func (c AnomaloClient) GetTableInformation(tableName string) GetTableResponse {
	var data GetTableResponse
	req := fmt.Sprintf("{\"table_name\": \"%s\"}", tableName)
	body := c.apiCallWithBody("get_table_information", http.MethodGet, req)
	defer body.Close()
	//x, _ := io.ReadAll(body)
	//fmt.Println(string(x))
	if err := json.NewDecoder(body).Decode(&data); err != nil {
		fmt.Println(err)
	}
	return data
}

func (c AnomaloClient) ConfigureTable(params string) ConfigureTableResponse {
	var data ConfigureTableResponse
	body := c.apiCallWithBody("configure_table", http.MethodPost, params)
	defer body.Close()
	//x, _ := io.ReadAll(body)
	//fmt.Println(string(x))
	if err := json.NewDecoder(body).Decode(&data); err != nil {
		fmt.Println(err)
	}
	return data
}

func main() {
	var client AnomaloClient
	jsonFile, err := os.Open("anomalo_secrets.json")
	if err != nil {
		fmt.Println(err)
	}
	contents, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(contents, &client)
	fmt.Println(client.Ping())
	fmt.Print("\n")
	table := client.GetTableInformation("cml_cloud_bot.inventory_health_development.variation_report")
	tableId := table.ID
	fmt.Print("\n")
	fmt.Printf("table ID: %d", tableId)
	fmt.Print("\n")
	config := `{
	  "table_id":"2011641",
	  "definition":"",
	  "check_cadence_run_at_duration":"PT6H",
	  "check_cadence_type":"daily",
	  "time_column_type":null,
	  "notify_after":null,
	  "notification_channel_id":"65",
	  "time_columns":null,
	  "fresh_after":null,
	  "interval_skip_expr":null
	}`
	fmt.Println(client.ConfigureTable(config))
}
