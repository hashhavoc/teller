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

// const DefaultApiBase = "https://api.hiro.so"
const DefaultApiBase = "https://stacks.hashhavoc.com"

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

func (c *APIClient) GetTransactions(principal string) ([]Transaction, error) {
	var allTxs []Transaction
	var total int
	offset := 0
	limit := 50 // or any other limit you want to set
	batchSize := 5
	resultsChan := make(chan []Transaction)
	errChan := make(chan error)
	doneChan := make(chan bool)

	go func() {
		for {
			currentOffset := offset
			for i := 0; i < batchSize; i++ {
				go func(o int) {
					url := fmt.Sprintf("%s/extended/v2/addresses/%s/transactions?offset=%d&limit=%d", c.BaseURL, principal, o, limit)
					req, err := http.NewRequest("GET", url, nil)
					if err != nil {
						errChan <- err
						return
					}
					req.Header.Add("Accept", "application/json")

					res, err := c.Client.Do(req)
					if err != nil {
						errChan <- err
						return
					}
					defer res.Body.Close()

					if res.StatusCode != 200 {
						errChan <- fmt.Errorf("failed to get transactions: %s", res.Status)
						return
					}

					body, err := io.ReadAll(res.Body)
					if err != nil {
						errChan <- err
						return
					}

					var txResp TransactionsResponse
					err = json.Unmarshal(body, &txResp)
					if err != nil {
						errChan <- err
						return
					}
					total = txResp.Total
					resultsChan <- txResp.Results
				}(currentOffset)
				currentOffset += limit
			}

			for i := 0; i < batchSize; i++ {
				select {
				case err := <-errChan:
					fmt.Println(err)
					return
				case results := <-resultsChan:
					allTxs = append(allTxs, results...)
				}
			}

			offset += limit * batchSize
			if len(allTxs) >= total { // Assuming a fixed number to break the loop, adjust based on your needs
				doneChan <- true
				return
			}
		}
	}()

	<-doneChan
	return allTxs, nil
}

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

func (c *APIClient) GetTokenHolders(contractID string, block int) (ContractHoldersResponse, error) {
	var url string
	if block == 0 {
		url = fmt.Sprintf("%s/extended/v1/address/%s/holders", c.BaseURL, contractID)
	} else {
		url = fmt.Sprintf("%s/extended/v1/address/%s/holders?until_block=%d", c.BaseURL, contractID, block)
	}
	method := "GET"

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return ContractHoldersResponse{}, err
	}

	req.Header.Add("Accept", "application/json")

	res, err := c.Client.Do(req)
	if err != nil {
		return ContractHoldersResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return ContractHoldersResponse{}, fmt.Errorf("failed to get contract holders: %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return ContractHoldersResponse{}, err
	}

	var response ContractHoldersResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return ContractHoldersResponse{}, err
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

func (c *APIClient) GetAccountBalance(principal string, block int) (BalanceResponse, error) {
	var url string
	if block == 0 {
		url = fmt.Sprintf("%s/extended/v1/address/%s/balances", c.BaseURL, principal)
	} else {
		url = fmt.Sprintf("%s/extended/v1/address/%s/balances?until_block=%d", c.BaseURL, principal, block)
	}
	method := "GET"

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return BalanceResponse{}, err
	}

	req.Header.Add("Accept", "application/json")

	res, err := c.Client.Do(req)
	if err != nil {
		return BalanceResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return BalanceResponse{}, fmt.Errorf("failed to get contract source: %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return BalanceResponse{}, err
	}

	var response BalanceResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return BalanceResponse{}, err
	}

	return response, nil
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

func (c *APIClient) GetAllNames() ([]Names, error) {
	var allResults []Names
	offset := 0
	limit := 100000 // or any other limit you want to set

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

// func (c *APIClient) GetAllNames() ([]string, error) {
// 	var allResults []string
// 	page := 1

// 	for {
// 		url := fmt.Sprintf("%s/v1/names?page=%d", c.BaseURL, page)
// 		req, err := http.NewRequest("GET", url, nil)
// 		if err != nil {
// 			return nil, err
// 		}
// 		req.Header.Add("Accept", "application/json")

// 		res, err := c.Client.Do(req)
// 		if err != nil {
// 			return nil, err
// 		}
// 		defer res.Body.Close()

// 		if res.StatusCode != 200 {
// 			return nil, fmt.Errorf("failed to fetch names: %s", res.Status)
// 		}

// 		body, err := io.ReadAll(res.Body)
// 		if err != nil {
// 			return nil, err
// 		}

// 		var response []string
// 		err = json.Unmarshal(body, &response)
// 		if err != nil {
// 			return nil, err
// 		}

// 		if len(response) == 0 {
// 			break
// 		}

// 		allResults = append(allResults, response...)
// 		page++
// 	}

// 	return allResults, nil
// }

func (c *APIClient) GetNames(page int) ([]string, error) {
	url := fmt.Sprintf("%s/v1/names?page=%d", c.BaseURL, page)
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
		return nil, fmt.Errorf("failed to fetch names: %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response []string
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
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
