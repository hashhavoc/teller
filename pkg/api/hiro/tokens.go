package hiro

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *APIClient) GetAllTokens() ([]TokenResult, error) {
	var allResults []TokenResult
	offset := 0
	limit := 60 // or any other limit you want to set

	for {
		url := fmt.Sprintf("%s/metadata/v1/ft?offset=%d&limit=%d", c.BaseURL, offset, limit)
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

		var response Response
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}

		allResults = append(allResults, response.Results...)

		if len(allResults) >= response.Total {
			break
		}

		offset += limit
	}

	return allResults, nil
}

func (c *APIClient) GetNFTHoldings(principal string) ([]NFTHoldingResponseResults, error) {
	url := fmt.Sprintf("%s/extended/v1/tokens/nft/holdings?principal=%s", c.BaseURL, principal)
	method := "GET"

	req, err := http.NewRequest(method, url, nil)
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

	var response NFTHoldingResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response.Results, nil
}
