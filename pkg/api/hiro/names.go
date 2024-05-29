package hiro

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *APIClient) GetAllNames() ([]Names, error) {
	var allResults []Names
	offset := 0
	limit := 100000

	for {
		url := fmt.Sprintf("%s/v2/names?offset=%d&limit=%d", c.BaseURL, offset, limit)
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

		var response NamesListResponse
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

func (c *APIClient) GetName(name string) (NameDetails, error) {
	url := fmt.Sprintf("%s/v1/names/%s", c.BaseURL, name)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return NameDetails{}, err
	}
	req.Header.Add("Accept", "application/json")

	res, err := c.Client.Do(req)
	if err != nil {
		return NameDetails{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return NameDetails{}, fmt.Errorf("failed to fetch names: %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return NameDetails{}, err
	}

	var response NameDetails
	err = json.Unmarshal(body, &response)
	if err != nil {
		return NameDetails{}, err
	}

	return response, nil
}

func (c *APIClient) GetNamesByAddress(address string) (NameReverseLookupResponse, error) {
	url := fmt.Sprintf("%s/v1/addresses/stacks/%s", c.BaseURL, address)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return NameReverseLookupResponse{}, err
	}
	req.Header.Add("Accept", "application/json")

	res, err := c.Client.Do(req)
	if err != nil {
		return NameReverseLookupResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return NameReverseLookupResponse{}, fmt.Errorf("failed to fetch names: %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return NameReverseLookupResponse{}, err
	}

	var response NameReverseLookupResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return NameReverseLookupResponse{}, err
	}

	return response, nil
}

func (c *APIClient) GetNameZoneFile(name string) (NameZoneFileResponse, error) {
	url := fmt.Sprintf("%s/v1/names/%s/zonefile", c.BaseURL, name)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return NameZoneFileResponse{}, err
	}
	req.Header.Add("Accept", "application/json")

	res, err := c.Client.Do(req)
	if err != nil {
		return NameZoneFileResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return NameZoneFileResponse{}, fmt.Errorf("failed to fetch names: %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return NameZoneFileResponse{}, err
	}

	var response NameZoneFileResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return NameZoneFileResponse{}, err
	}

	return response, nil
}
