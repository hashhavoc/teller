package stxtools

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const DefaultApiBase = "https://api.stxtools.io"

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

func (c *APIClient) GetAllTokens() ([]Token, error) {
	url := fmt.Sprintf("%s/tokens", c.BaseURL)
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
		return nil, fmt.Errorf("failed to get tokens: %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response []Token
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *APIClient) GetAllHolders(contractId string) (HoldersData, error) {
	var allResults HoldersData
	page := 0
	limit := 50 // or any other limit you want to set

	for {
		url := fmt.Sprintf("%s/tokens/%s/top-holders?page=%d&limit=%d", c.BaseURL, contractId, page, limit)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return HoldersData{}, err
		}
		req.Header.Add("Accept", "application/json")

		res, err := c.Client.Do(req)
		if err != nil {
			return HoldersData{}, err
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			return HoldersData{}, fmt.Errorf("failed to get contract source: %s", res.Status)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return HoldersData{}, err
		}

		var response HoldersResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			return HoldersData{}, err
		}

		allResults.TopHolders = append(allResults.TopHolders, response.Data.TopHolders...) // Assuming TopHolders is the slice of TokenResult

		// Assuming the Response struct has a Page field that includes TotalElements
		if len(allResults.TopHolders) >= response.Page.TotalElements {
			break
		}
		page += 1
	}

	return allResults, nil
}

func (c *APIClient) GetAllSwaps(contractId string) ([]SwapsData, error) {
	var allSwaps []SwapsData
	page := 0
	limit := 50 // Adjust the limit as needed

	for {
		url := fmt.Sprintf("%s/tokens/%s/swaps?page=%d&limit=%d", c.BaseURL, contractId, page, limit)
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
			return nil, fmt.Errorf("failed to get swaps: %s", res.Status)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		var response SwapsResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}

		allSwaps = append(allSwaps, response.Data...) // Assuming Data is the slice of SwapsData

		// Check if we've fetched all available data
		if len(allSwaps) >= response.Page.TotalElements {
			break
		}
		page += 1
	}

	return allSwaps, nil
}

func (c *APIClient) GetAllTransfers(contractId string) ([]Transaction, error) {
	var allTransfers []Transaction
	page := 0
	limit := 50 // Adjust the limit as needed

	for {
		url := fmt.Sprintf("%s/tokens/%s/transfers?page=%d&limit=%d", c.BaseURL, contractId, page, limit)
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
			return nil, fmt.Errorf("failed to get transfers: %s", res.Status)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		var response TransfersResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}

		allTransfers = append(allTransfers, response.Data...) // Assuming Data is the slice of Transaction

		// Check if we've fetched all available data
		if len(allTransfers) >= response.Page.TotalElements {
			break
		}
		page += 1
	}

	return allTransfers, nil
}
