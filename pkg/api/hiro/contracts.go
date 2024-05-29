package hiro

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hashhavoc/teller/pkg/utils/uint128"
)

func (c *APIClient) GetContractDetails(contractID string) (ContractDetailsResponse, error) {
	url := fmt.Sprintf("%s/extended/v1/contract/%s", c.BaseURL, contractID)
	method := "GET"

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return ContractDetailsResponse{}, err
	}

	req.Header.Add("Accept", "application/json")

	res, err := c.Client.Do(req)
	if err != nil {
		return ContractDetailsResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return ContractDetailsResponse{}, fmt.Errorf("failed to get contract details: %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return ContractDetailsResponse{}, err
	}

	var response ContractDetailsResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return ContractDetailsResponse{}, err
	}

	return response, nil
}

func (c *APIClient) GetContractSource(id string) (string, error) {
	split, err := ContractValidateSplit(id)
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("%s/v2/contracts/source/%s/%s", c.BaseURL, split[0], split[1])
	method := "GET"

	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return "", err
	}
	req.Header.Add("Accept", "application/json")

	res, err := c.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", fmt.Errorf("failed to get contract source: %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var response ContractSourceResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	return response.Source, nil
}

func ContractValidateSplit(c string) ([]string, error) {
	split := strings.Split(c, ".")
	if len(split) != 2 {
		return nil, errors.New("invalid contract ID format")
	}
	if strings.Contains(split[1], "::") {
		return nil, errors.New("invalid contract ID format please remove :: and aything after")
	}
	return split, nil
}

func (c *APIClient) GetContractReadOnly(id string, function string, responseType string, arguments []string) (string, error) {
	split, err := ContractValidateSplit(id)
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("%s/v2/contracts/call-read/%s/%s/%s", c.BaseURL, split[0], split[1], function)
	method := "POST"

	rawPayload := ReadOnlyPayload{
		// random address
		Sender:    "SP3D49HARD6Y36MKPT3PKP2YHG0ZNQMK0YP70RZHS",
		Arguments: arguments,
	}

	jsonPayload, err := json.Marshal(rawPayload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := c.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", fmt.Errorf("failed to get contract source: %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var response ReadOnlyResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	switch t := responseType; t {
	case "string":
		if len(response.Result) < 14 {
			return "", fmt.Errorf("invalid response: %s", response.Result)
		}
		decoded, err := hex.DecodeString(response.Result[14:])
		if err != nil {
			return "", err
		}
		switch tt := function; tt {
		case "get-token-uri":
			return string(decoded)[1:], nil
		default:
			return string(decoded), nil
		}
	case "uint128":
		if len(response.Result) < 6 {
			return "", fmt.Errorf("invalid response: %s", response.Result)
		}
		u, err := uint128.FromString(response.Result[6:])
		if err != nil {
			return "", err
		}
		s := u.String()
		return s, nil

	default:
		if len(response.Result) < 14 {
			return "", fmt.Errorf("invalid response: %s", response.Result)
		}
		decoded, err := hex.DecodeString(response.Result[14:])
		if err != nil {
			return "", err
		}
		return string(decoded), nil
	}
}
