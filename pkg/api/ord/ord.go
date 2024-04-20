package ord

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const DefaultApiBase = "http://localhost:8080"

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

func (c *APIClient) GetAllRunes() ([]Entry, error) {
	var allEntries []Entry
	offset := 1

	for {
		url := fmt.Sprintf("%s/runes/%d", c.BaseURL, offset)
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
			return nil, fmt.Errorf("failed to get runes: %s", res.Status)
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

		for _, entry := range response.Entries {
			var e Entry
			e.ID = entry[0].(string)
			details, ok := entry[1].(map[string]interface{})
			if !ok {
				return nil, errors.New("failed to convert entry details")

			}
			e.Details.Block = int64(details["block"].(float64))
			e.Details.Burned = int64(details["burned"].(float64))
			e.Details.Divisibility = int64(details["divisibility"].(float64))
			e.Details.Etching = details["etching"].(string)
			e.Details.Mints = int64(details["mints"].(float64))
			e.Details.Number = int64(details["number"].(float64))
			e.Details.Premine = int64(details["premine"].(float64))
			e.Details.SpacedRune = details["spaced_rune"].(string)
			if terms, ok := details["terms"].(map[string]interface{}); ok {
				if amount, ok := terms["amount"].(float64); ok {
					e.Details.Terms.Amount = int64(amount)
				}
			}
			allEntries = append(allEntries, e)

		}

		if !response.More {
			break
		}

		offset += 1
	}
	return allEntries, nil
}
