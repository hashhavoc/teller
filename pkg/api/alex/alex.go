package alex

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const DefaultApiBase = "https://api.alexgo.io"
const DefaultGraphQLEndpoint = "https://gql.alexlab.co/v1/graphql"

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

func (c *APIClient) FetchLatestPrices() (TokenPriceResponse, error) {
	const query = `{"query":"query FetchLatestPrices { laplace_current_token_price { avg_price_usd token } }"}`

	body, err := c.ExecuteGraphQLQuery(query)
	if err != nil {
		return TokenPriceResponse{}, err
	}

	var response TokenPriceResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return TokenPriceResponse{}, err
	}

	return response, nil
}

func (c *APIClient) ExecuteGraphQLQuery(query string) ([]byte, error) {
	payload := strings.NewReader(query)

	req, err := http.NewRequest("POST", DefaultGraphQLEndpoint, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("failed to make graphql request: %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *APIClient) GetPairs() ([]CurrencyPair, error) {
	url := fmt.Sprintf("%s/v2/coin-gecko/tickers", c.BaseURL)

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
		return nil, fmt.Errorf("failed to get pairs: %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response []CurrencyPair
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
