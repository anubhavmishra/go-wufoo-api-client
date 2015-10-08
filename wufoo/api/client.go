package api

import (
	"encoding/json"
	"fmt"
	"github.com/itembase/go-wufoo-api-client/wufoo"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Client represents API client for Wufoo
type Client struct {
	Config     wufoo.Config
	formsApi   *FormsApi
	fieldsApi  *FieldsApi
	entriesApi *EntriesApi
}

func (c *Client) FormsApi() *FormsApi {
	if c.formsApi == nil {
		c.formsApi = new(FormsApi)
		c.formsApi.Client = c
	}

	return c.formsApi
}

func (c *Client) FieldsApi() *FieldsApi {
	if c.fieldsApi == nil {
		c.fieldsApi = new(FieldsApi)
		c.fieldsApi.Client = c
	}

	return c.fieldsApi
}

func (c *Client) EntriesApi() *EntriesApi {
	if c.entriesApi == nil {
		c.entriesApi = new(EntriesApi)
		c.entriesApi.Client = c
	}

	return c.entriesApi
}

// get function makes a GET request to wufoo API and returns
func (c Client) Get(api string, params map[string]string, response interface{}) (err error) {
	req, err := http.NewRequest("GET", c.PrepareUrl(api, params), nil)
	if err != nil {
		return
	}

	req.SetBasicAuth(c.Config.ApiKey, "footastic")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return
	}

	return nil
}

// get function makes a GET request to wufoo API and returns
func (c Client) Post(api string, postData url.Values, response interface{}) (err error) {
	req, err := http.NewRequest("POST", c.PrepareUrl(api, nil), strings.NewReader(postData.Encode()))
	if err != nil {
		return
	}

	req.SetBasicAuth(c.Config.ApiKey, "footastic")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if response != nil {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		err = json.Unmarshal(body, &response)
		if err != nil {
			return err
		}
	}

	return nil
}

// prepareUrl function returns URL for specific API with GET parameters if specified
func (c Client) PrepareUrl(api string, params map[string]string) (urlValue string) {
	urlValue = fmt.Sprintf(wufoo.WUFOO_API_URL_PATTERN, c.Config.Subdomain, api)

	if params != nil && len(params) > 0 {
		getParams := url.Values{}
		for key, value := range params {
			getParams.Add(key, value)
		}

		urlValue += "?" + getParams.Encode()
	}

	return
}
