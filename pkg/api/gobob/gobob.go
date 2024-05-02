package gobob

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/go-querystring/query"
)

const DefaultApiBase = "https://explorer.gobob.xyz"

type APIClient struct {
	BaseURL string
	Client  *http.Client
}

func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}

func (c *APIClient) GetAllTokens() ([]TokenItems, error) {
	var allResults []TokenItems
	var params NextPageParams

	for {
		var url string
		if params == (NextPageParams{}) {
			url = fmt.Sprintf("%s/api/v2/tokens", c.BaseURL)
		} else {
			v, _ := query.Values(params)
			url = fmt.Sprintf("%s/api/v2/tokens?%s", c.BaseURL, v.Encode())
		}

		fmt.Println(url)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Add("Accept", "application/json")

		res, err := c.Client.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			return nil, fmt.Errorf("failed to get contract source: %s", res.Status)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		var response TokenResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}
		allResults = append(allResults, response.Items...)
		params = response.NextPageParams

		if response.NextPageParams == (NextPageParams{}) {
			break
		}
	}

	return allResults, nil
}
