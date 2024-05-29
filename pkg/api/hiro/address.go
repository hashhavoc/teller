package hiro

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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
