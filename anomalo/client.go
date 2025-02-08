// Package anomalo is a Go Client for interacting with the Anomalo API.
//
// The package provides structs to work with request and response objects, makes
// it easy to construct HTTP requests, and provides a few convenience methods
// for common operations.
package anomalo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"golang.org/x/exp/maps"
)

// HttpClientProvider An interface that provides an HTTP client - used for
// testing and for package users that need control over network requests.
type HttpClientProvider func() *http.Client

var (
	// ValidNotificationChannels A set of supported notification channel types
	ValidNotificationChannels = map[string]struct{}{
		"email":     {},
		"email_all": {},
		"msteams":   {},
		"opsgenie":  {},
		"pagerduty": {},
		"slack":     {},
		"webhook":   {},
	}
)

// Client The Anomalo client - used to authenticate & make calls to the Anomalo
// API.
type Client struct {
	Token          string `json:"Token,omitempty"`
	Host           string `json:"Host,omitempty"`
	ClientProvider HttpClientProvider
	client         *http.Client
}

func closeBody(body io.ReadCloser) {
	err := body.Close()
	if err != nil {
		log.Printf("Error closing body. %s", err.Error())
	}
}

// Memo-ized client
func (c *Client) getClient() *http.Client {
	if c.client == nil {
		if c.ClientProvider != nil {
			c.client = c.ClientProvider()
		} else {
			c.client = &http.Client{}
		}
	}
	return c.client
}

func (c *Client) buildUrl(endpoint string) string {
	return fmt.Sprintf("%s/api/public/v1/%s", c.Host, endpoint)
}

func (c *Client) apiCall(endpoint string, method string) (*http.Response, error) {
	return c.apiCallWithBody(endpoint, method, "{}")
}

// apiCallWithBody Builds an HTTP request to Anomalo with the given JSON
// parameters. Encodes them in the request body for PUT and POST requests, and
// encodes them as URL parameters for all other HTTP method types.
func (c *Client) apiCallWithBody(endpoint string, method string, jsonParams string) (*http.Response, error) {
	var req *http.Request
	var err error
	if method == http.MethodPost || method == http.MethodPut {
		req, err = http.NewRequest(method, c.buildUrl(endpoint), bytes.NewBuffer([]byte(jsonParams)))
		if err != nil {
			return nil, err
		}
	} else {
		req, err = http.NewRequest(method, c.buildUrl(endpoint), nil)
		if err != nil {
			return nil, err
		}
		var parsed map[string]interface{}
		err = json.Unmarshal([]byte(jsonParams), &parsed)
		if err != nil {
			return nil, err
		}
		params := req.URL.Query()
		for key, value := range parsed {
			params.Add(key, fmt.Sprintf("%v", value))
		}
		req.URL.RawQuery = params.Encode()
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.getClient().Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("response code %d. unable to read response body. got %w", resp.StatusCode, err)
		}
		return nil, fmt.Errorf(string(bodyBytes))
	}

	return resp, nil
}

func (c *Client) Ping() (*PingResponse, error) {
	var data *PingResponse
	resp, err := c.apiCall("ping", http.MethodGet)
	if err != nil {
		return nil, err
	}
	body := resp.Body
	defer closeBody(body)
	if err := json.NewDecoder(body).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

// GetTableInformation looks up a table by `tableName`.
// `tableName` must start with its warehouse name.
func (c *Client) GetTableInformation(tableName string) (*GetTableResponse, error) {
	var data *GetTableResponse
	req := fmt.Sprintf("{\"table_name\": \"%s\"}", tableName)
	resp, err := c.apiCallWithBody("get_table_information", http.MethodGet, req)
	if err != nil {
		return nil, err
	}
	body := resp.Body
	defer closeBody(body)
	if err := json.NewDecoder(body).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

// GetTableInformationWithWarehouseID takes two arguments
// - `tableName` is the name of the schema without the warehouse string prefix
// - `warehouseID` is the integer ID of the corresponding warehouse
//
// This method was added because `GetTableInformation` will fail if a table's
// warehouse name is not unique within an Anomalo workspace. Therefore, if
// there are multiple warehouses with the same name, then you should differentiate
// via the warehouseID parameter instead.
func (c *Client) GetTableInformationWithWarehouseID(tableName string, warehouseID int) (*GetTableResponse, error) {
	var data *GetTableResponse
	req := fmt.Sprintf("{\"table_name\": \"%s\", \"warehouse_id\": \"%d\"}", tableName, warehouseID)
	resp, err := c.apiCallWithBody("get_table_information", http.MethodGet, req)
	if err != nil {
		return nil, err
	}
	body := resp.Body
	defer closeBody(body)
	if err := json.NewDecoder(body).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) ConfigureTable(req ConfigureTableRequest) (*ConfigureTableResponse, error) {
	var data *ConfigureTableResponse
	reqJson, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := c.apiCallWithBody("configure_table", http.MethodPost, string(reqJson))
	if err != nil {
		return nil, err
	}
	body := resp.Body
	defer closeBody(body)
	if err := json.NewDecoder(body).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) GetChecks(tableID int) (*GetChecksResponse, error) {
	var data *GetChecksResponse
	req := fmt.Sprintf("{\"table_id\": \"%d\"}", tableID)
	resp, err := c.apiCallWithBody("get_checks_for_table", http.MethodGet, req)
	if err != nil {
		return nil, err
	}
	body := resp.Body
	defer closeBody(body)
	if err := json.NewDecoder(body).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

// GetCheckByStaticID Wrapper around GetChecks that additionally filters checks
// by static ID. Returns nil if a matching check is not found.
//
// Note that this method does a linear search through all checks, since the
// Anomalo API does not allow queries by check static ID. This shouldn't be
// untenable since it's unlikely to have so many checks on a table that this
// iteration becomes slow.
func (c *Client) GetCheckByStaticID(tableID int, staticID int) (*Check, error) {
	data, err := c.GetChecks(tableID)
	if err != nil {
		return nil, err
	}

	var relevantCheck *Check
	for _, check := range data.Checks {
		if check.CheckStaticID == staticID {
			if relevantCheck != nil {
				return nil, fmt.Errorf("saw more than one check with the same static ID %d for table %d. "+
					"check IDs %d & %d", staticID, tableID, relevantCheck.CheckID, check.CheckID)
			}
			temp := check // Necessary because the address of `check` stays the same
			relevantCheck = &temp
		}
	}

	return relevantCheck, nil
}

// GetCheckByRef Wrapper around GetChecks that additionally filters checks
// by Ref. Returns nil if a matching check is not found.
//
// Note that this method does a linear search through all checks, since the
// Anomalo API does not allow queries by check static ID. This shouldn't be
// untenable since it's unlikely to have so many checks on a table that this
// iteration becomes slow.
func (c *Client) GetCheckByRef(tableID int, ref string) (*Check, error) {
	data, err := c.GetChecks(tableID)
	if err != nil {
		return nil, err
	}

	var relevantCheck *Check
	for _, check := range data.Checks {
		if check.Ref == ref {
			if relevantCheck != nil {
				return nil, fmt.Errorf("saw more than one check with the same ref %s for table %d. "+
					"check IDs %d & %d", ref, tableID, relevantCheck.CheckID, check.CheckID)
			}
			temp := check // Necessary because the address of `check` stays the same
			relevantCheck = &temp
		}
	}

	return relevantCheck, nil
}

func (c *Client) CreateCheck(req CreateCheckRequest) (*CreateCheckResponse, error) {
	var data *CreateCheckResponse
	reqJson, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := c.apiCallWithBody("create_check", http.MethodPost, string(reqJson))
	if err != nil {
		return nil, err
	}
	body := resp.Body
	defer closeBody(body)
	if err := json.NewDecoder(body).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) DeleteCheck(req DeleteCheckRequest) (*DeleteCheckResponse, error) {
	var data *DeleteCheckResponse
	reqJson, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := c.apiCallWithBody("delete_check", http.MethodPost, string(reqJson))
	if err != nil {
		return nil, err
	}
	body := resp.Body
	defer closeBody(body)
	if err := json.NewDecoder(body).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) RunChecks(req RunChecksRequest) (*RunChecksResponse, error) {
        var data *RunChecksResponse
        reqJson, err := json.Marshal(req)
        if err != nil {
                return nil, err
        }
        resp, err := c.apiCallWithBody("run_checks", http.MethodPost, string(reqJson))
        if err != nil {
                return nil, err
        }
        body := resp.Body
        defer closeBody(body)
        if err := json.NewDecoder(body).Decode(&data); err != nil {
                return nil, err
        }
        return data, nil
}

func (c *Client) GetNotificationChannels() (*GetNotificationChannelsResponse, error) {
	var data *GetNotificationChannelsResponse
	resp, err := c.apiCall("list_notification_channels", http.MethodGet)
	if err != nil {
		return nil, err
	}
	body := resp.Body
	defer closeBody(body)
	if err := json.NewDecoder(body).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

// GetNotificationChannelWithDescriptionContaining Wrapper around
// GetNotificationChannels that additionally filters the channels by name. Note
// that Anomalo does not provide a Description that makes it easy to
// programmatically match by name, so we use strings.Contains.
//
// Returns nil if a matching channel is not found.
//
// Also note that this method does a linear search through all notification
// channels due to limitations with the Anomalo API. This is unlikely to be an
// issue since we don't expect large numbers of notification channels.
func (c *Client) GetNotificationChannelWithDescriptionContaining(
	name string,
	channelType string,
) (*NotificationChannel, error) {
	if _, ok := ValidNotificationChannels[channelType]; !ok {
		return nil, fmt.Errorf("channelType must be one of %v", maps.Keys(ValidNotificationChannels))
	}

	channels, err := c.GetNotificationChannels()
	if err != nil {
		return nil, err
	}

	var matchingChannel *NotificationChannel
	for i, channel := range channels.NotificationChannels {
		if channel.ChannelType != channelType || !strings.Contains(channel.Description, name) {
			continue
		}
		if matchingChannel != nil {
			return nil, fmt.Errorf("Found at least two channels with descriptions containing string %s. "+
				"Channel 1 description: %s\nChannel 2 description: %s",
				name, matchingChannel.Description, channel.Description)
		}
		matchingChannel = &channels.NotificationChannels[i]
	}
	return matchingChannel, nil
}

func (c *Client) GetOrganizations() ([]*Organization, error) {
	var data []*Organization
	resp, err := c.apiCall("organizations", http.MethodGet)
	if err != nil {
		return nil, err
	}
	body := resp.Body
	defer closeBody(body)
	if err := json.NewDecoder(body).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

// GetOrganizationByName Wrapper around GetOrganizations that looks for an
// organization by name.
//
// Returns nil if a channel with a case-sensitive exact match is not found.
//
// Also note that this method does a linear search through all notification
// channels due to limitations with the Anomalo API. This is unlikely to be an
// issue since we don't expect large numbers of notification channels.
func (c *Client) GetOrganizationByName(name string) (*Organization, error) {
	orgs, err := c.GetOrganizations()
	if err != nil {
		return nil, err
	}

	names := make([]string, len(orgs))
	for i, org := range orgs {
		names[i] = org.Name
		if org.Name == name {
			return org, nil
		}
	}

	return nil, fmt.Errorf("did not find an organization with name %s. Make sure it exists, and this API "+
		"token has access. Found names: %s", name, strings.Join(names, ", "))
}

// ChangeOrganization API keys have permissions scoped to a given Organization. An API key can only act within the scope
// of one organization at a time. Call ChangeOrganization to change the Organization the API Key is acting within.
func (c *Client) ChangeOrganization(orgId int64) (*ChangeOrganizationResponse, error) {
	var data *ChangeOrganizationResponse
	req := fmt.Sprintf("{\"id\": \"%d\"}", orgId)
	resp, err := c.apiCallWithBody("organization", http.MethodPut, req)
	if err != nil {
		return nil, err
	}
	body := resp.Body
	defer closeBody(body)
	if err := json.NewDecoder(body).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

// For debugging
func responseToString(resp *http.Response) string {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(bodyBytes)
}
